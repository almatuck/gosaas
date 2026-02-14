package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gosaas/app"
	"gosaas/internal/config"
	"gosaas/internal/handler"
	"gosaas/internal/mcp"
	mcpoauth "gosaas/internal/mcp/oauth"
	"gosaas/internal/middleware"
	"gosaas/internal/oauth"
	"gosaas/internal/realtime"
	"gosaas/internal/svc"
	"gosaas/internal/webhook"
	"gosaas/internal/websocket"

	levee "github.com/almatuck/levee-go"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed etc/gosaas.yaml
var embeddedConfig []byte

func main() {
	// Load config with environment variable expansion
	c, err := config.LoadConfig([]byte(os.ExpandEnv(string(embeddedConfig))))
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Determine server host based on config
	var srvHost = c.Host
	var serverPort = c.Port
	var useHTTPS = false

	if c.IsProductionMode() {
		if c.App.Domain == "" {
			fmt.Println("ERROR: App.Domain is required in production mode")
			os.Exit(1)
		}
		srvHost = c.App.Domain
		serverPort = 443
		useHTTPS = true
		fmt.Printf("Running in PRODUCTION mode - server.json will return https://%s\n", c.App.Domain)
	} else if serverPort == 443 || serverPort == 80 {
		if c.App.Domain == "" {
			fmt.Println("ERROR: App.Domain is required when using ports 80/443")
			os.Exit(1)
		}
		srvHost = c.App.Domain
		if serverPort == 443 {
			useHTTPS = true
		}
	} else {
		srvHost = "localhost"
		app.DevMode = true
		fmt.Printf("Running in DEVELOPMENT mode - server.json will return http://localhost:%d\n", serverPort)
	}

	fmt.Println("Server Host:", srvHost, "Port:", serverPort, "Use HTTPS:", useHTTPS)

	// Set server host for server.json
	app.SetServerHost(srvHost, serverPort, useHTTPS)

	// Set up SPA filesystem for static file serving
	spaFS, spaErr := app.FileSystem()
	if spaErr != nil {
		fmt.Printf("Warning: Could not load embedded SPA files: %v\n", spaErr)
		fmt.Println("App should be running separately on port 5173")
	}

	// Create chi router
	r := chi.NewRouter()

	// Standard middleware
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	svcCtx := svc.NewServiceContext(*c)
	defer svcCtx.Close()

	// Global security headers middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.IsSecurityHeadersEnabled() {
				headers := middleware.APISecurityHeaders()
				w.Header().Set("Content-Security-Policy", headers.ContentSecurityPolicy)
				w.Header().Set("X-Content-Type-Options", headers.XContentTypeOptions)
				w.Header().Set("X-Frame-Options", headers.XFrameOptions)
				w.Header().Set("X-XSS-Protection", headers.XXSSProtection)
				w.Header().Set("Referrer-Policy", headers.ReferrerPolicy)
				w.Header().Set("Permissions-Policy", headers.PermissionsPolicy)
				if c.IsForceHTTPS() {
					w.Header().Set("Strict-Transport-Security", headers.StrictTransportSecurity)
				}
				w.Header().Set("Cache-Control", headers.CacheControl)
				w.Header().Set("Pragma", headers.Pragma)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register API routes
	handler.RegisterHandlers(r, svcCtx)

	// Register Stripe webhook for standalone mode
	if svcCtx.UseLocal() && svcCtx.Config.Stripe.WebhookSecret != "" {
		r.Post("/api/webhook/stripe", webhook.StripeHandler(svcCtx))
		fmt.Println("Stripe webhook registered at /api/webhook/stripe")

		// Sync products to Stripe on startup
		if svcCtx.Billing != nil && len(svcCtx.Config.Products) > 0 {
			syncCtx := context.Background()
			syncedProducts, err := svcCtx.Billing.SyncProductsToStripe(syncCtx)
			if err != nil {
				fmt.Printf("Warning: Failed to sync products to Stripe: %v\n", err)
			} else {
				fmt.Printf("Synced %d products to Stripe\n", len(syncedProducts))
				svcCtx.Config.Products = syncedProducts
			}
		}
	}

	// Register Levee embedded handlers on default mux (reverse proxy routes to these)
	if svcCtx.Levee != nil {
		svcCtx.Levee.RegisterHandlers(http.DefaultServeMux, "",
			levee.WithUnsubscribeRedirect("/unsubscribed"),
			levee.WithConfirmRedirect("/welcome"),
			levee.WithConfirmExpiredRedirect("/confirm-expired"),
		)
	}

	// Register OAuth callback handlers
	if svcCtx.UseLocal() && c.IsOAuthEnabled() {
		oauthHandler := oauth.NewHandler(svcCtx)
		oauthHandler.RegisterRoutes(http.DefaultServeMux)
		fmt.Println("OAuth callbacks registered at /oauth/{provider}/callback")
	}

	// Register MCP (Model Context Protocol) handler
	if svcCtx.UseLocal() {
		var baseURL string
		if useHTTPS {
			baseURL = fmt.Sprintf("https://%s", srvHost)
		} else {
			baseURL = fmt.Sprintf("http://%s:%d", srvHost, serverPort)
		}

		mcpHandler := mcp.NewHandler(svcCtx, baseURL)
		http.DefaultServeMux.Handle("/mcp", mcpHandler)
		http.DefaultServeMux.Handle("/mcp/", mcpHandler)

		mcpOAuthHandler := mcpoauth.NewHandler(svcCtx, baseURL)
		mcpOAuthHandler.RegisterRoutes(http.DefaultServeMux)

		fmt.Println("MCP endpoint registered at /mcp")
		fmt.Println("MCP OAuth endpoints registered at /.well-known/oauth-* and /mcp/oauth/*")
	}

	// Create WebSocket hub for real-time events
	hub := realtime.NewHub()
	go hub.Run(context.Background())

	// Register rewrite handler for WebSocket messages
	rewriteHandler := realtime.NewRewriteHandler(svcCtx)
	rewriteHandler.Register()

	// WebSocket endpoint
	r.Get("/ws", websocket.Handler(hub))

	// SPA fallback (must be last)
	if spaErr == nil {
		r.NotFound(app.NotFoundHandler(spaFS).ServeHTTP)
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	// In development mode, serve directly
	if app.DevMode {
		fmt.Printf("Starting server on %s (dev mode)...\n", addr)
		server := &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Production mode: Start backend in background, then HTTPS/HTTP servers
	backendServer := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Starting backend server on %s...\n", addr)
		if err := backendServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Backend server error: %v\n", err)
		}
	}()

	// Set up autocert for Let's Encrypt
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(c.App.Domain, "www."+c.App.Domain),
		Email:      c.App.AdminEmail,
	}

	// Create reverse proxy to backend with connection pooling
	backendURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", c.Host, c.Port))
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		if req.Header.Get("Upgrade") != "" {
			req.Header.Set("Connection", "Upgrade")
		}
	}

	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		WriteBufferSize:     32 << 10,
		ReadBufferSize:      32 << 10,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("Proxy error: %v\n", err)
		http.Error(w, "Backend temporarily unavailable", http.StatusBadGateway)
	}

	// HTTP handler for port 80 - ACME challenges and HTTPS redirect
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") {
			certManager.HTTPHandler(nil).ServeHTTP(w, r)
			return
		}
		host, _ := strings.CutPrefix(r.Host, "www.")
		newURL := fmt.Sprintf("https://%s%s", host, r.RequestURI)
		http.Redirect(w, r, newURL, http.StatusMovedPermanently)
	})

	// HTTPS handler for port 443
	baseHTTPSHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// www → non-www redirect
		if nonWWWHost, hadPrefix := strings.CutPrefix(r.Host, "www."); hadPrefix {
			newURL := fmt.Sprintf("https://%s%s", nonWWWHost, r.RequestURI)
			http.Redirect(w, r, newURL, http.StatusMovedPermanently)
			return
		}

		// Route API requests to backend
		if strings.HasPrefix(r.URL.Path, "/api/") {
			proxy.ServeHTTP(w, r)
			return
		}

		// Route Levee webhook paths to default mux (stripe, ses)
		if r.URL.Path == "/webhooks/stripe" || r.URL.Path == "/webhooks/ses" {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}

		// Route other webhook requests to backend
		if strings.HasPrefix(r.URL.Path, "/webhooks/") {
			proxy.ServeHTTP(w, r)
			return
		}

		// Route Levee email tracking/confirmation paths to default mux
		if strings.HasPrefix(r.URL.Path, "/e/") || r.URL.Path == "/confirm-email" {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}

		// Route OAuth callbacks to default mux
		if strings.HasPrefix(r.URL.Path, "/oauth/") {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}

		// Route MCP requests to default mux
		if strings.HasPrefix(r.URL.Path, "/mcp") || strings.HasPrefix(r.URL.Path, "/.well-known/oauth-") {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}

		// Route WebSocket requests to backend
		if strings.HasPrefix(r.URL.Path, "/ws") {
			proxyWebSocket(w, r, c.Host, c.Port)
			return
		}

		// Serve static files with SPA fallback
		if spaErr == nil {
			app.SPAHandler(spaFS).ServeHTTP(w, r)
		} else {
			http.Error(w, "SPA not available", http.StatusServiceUnavailable)
		}
	})

	// Layer middlewares: Gzip → CacheControl → Handler
	httpsHandler := middleware.Gzip(middleware.CacheControl(baseHTTPSHandler))

	// Start HTTP server on port 80
	httpServer := &http.Server{
		Addr:         ":80",
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("Starting HTTP server on :80 for ACME challenges and HTTPS redirect...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Start HTTPS server on port 443
	httpsServer := &http.Server{
		Addr:    ":443",
		Handler: httpsHandler,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
			NextProtos:     []string{"h2", "http/1.1", "acme-tls/1"},
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
		},
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("Starting HTTPS server on :443 with Let's Encrypt auto-certificate...")
	fmt.Println("Auto-redirect: www → non-www")
	fmt.Println("Auto-redirect: HTTP → HTTPS")
	fmt.Println("API routes: /api/* (proxied to backend)")
	fmt.Println("Static SPA: /* (served directly from embedded FS)")
	fmt.Println("Optimizations: HTTP/2, connection pooling, compression")

	go func() {
		if err := httpsServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTPS server error: %v\n", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down servers gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpsServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTPS server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("HTTPS server stopped")
	}

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("HTTP server stopped")
	}

	if err := backendServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Backend server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Backend server stopped")
	}

	fmt.Println("All servers shut down successfully")
}

// proxyWebSocket handles WebSocket upgrade and bidirectional proxying
func proxyWebSocket(w http.ResponseWriter, r *http.Request, backendHost string, backendPort int) {
	fmt.Printf("[WS Proxy] Incoming WebSocket request: %s\n", r.URL.String())

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("[WS Proxy] ERROR: ResponseWriter does not support Hijack")
		http.Error(w, "WebSocket not supported", http.StatusInternalServerError)
		return
	}

	backendAddr := fmt.Sprintf("%s:%d", backendHost, backendPort)
	fmt.Printf("[WS Proxy] Dialing backend: %s\n", backendAddr)
	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		fmt.Printf("[WS Proxy] ERROR: Failed to dial backend: %v\n", err)
		http.Error(w, "Backend unavailable", http.StatusBadGateway)
		return
	}
	defer backendConn.Close()

	clientConn, clientBuf, err := hijacker.Hijack()
	if err != nil {
		fmt.Printf("[WS Proxy] ERROR: Hijack failed: %v\n", err)
		http.Error(w, "Hijack failed", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()
	fmt.Println("[WS Proxy] Connection hijacked successfully")

	if err := r.Write(backendConn); err != nil {
		fmt.Printf("[WS Proxy] ERROR: Failed to forward request: %v\n", err)
		return
	}
	fmt.Println("[WS Proxy] Request forwarded to backend")

	if clientBuf.Reader.Buffered() > 0 {
		buffered := make([]byte, clientBuf.Reader.Buffered())
		clientBuf.Read(buffered)
		backendConn.Write(buffered)
	}

	done := make(chan struct{}, 2)
	go func() {
		io.Copy(backendConn, clientConn)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(clientConn, backendConn)
		done <- struct{}{}
	}()
	<-done
	fmt.Println("[WS Proxy] Connection closed")
}
