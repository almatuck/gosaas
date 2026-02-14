package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// RegisterHandler - Register new user
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := &RegisterLogic{
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.LoginResponse, err error) {
	// Use local auth when Levee is disabled
	if l.svcCtx.UseLocal() {
		return l.registerLocal(req)
	}

	// Use Levee when enabled
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("auth service not configured")
	}

	// Default to free plan if not specified
	plan := req.Plan
	if plan == "" {
		plan = "free"
	}

	// Build success/cancel URLs from config
	baseURL := l.svcCtx.Config.App.BaseURL
	successURL := baseURL + l.svcCtx.Config.Levee.CheckoutSuccessURL
	cancelURL := baseURL + l.svcCtx.Config.Levee.CheckoutCancelURL

	// Register via Levee SDK
	authResp, err := l.svcCtx.Levee.Auth.Register(l.ctx, &levee.SDKRegisterRequest{
		Email:         req.Email,
		Password:      req.Password,
		Name:          req.Name,
		PriceNickname: plan,
		SuccessUrl:    successURL,
		CancelUrl:     cancelURL,
	})
	if err != nil {
		l.Errorf("Registration failed for %s: %v", req.Email, err)
		return nil, err
	}

	// Parse expiry time
	expiresAt, _ := time.Parse(time.RFC3339, authResp.ExpiresAt)

	l.Infof("User registered: %s (plan: %s)", req.Email, plan)

	return &types.LoginResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    expiresAt.UnixMilli(),
		CheckoutUrl:  authResp.CheckoutUrl,
	}, nil
}

// registerLocal handles registration with local SQLite auth
func (l *RegisterLogic) registerLocal(req *types.RegisterRequest) (*types.LoginResponse, error) {
	if l.svcCtx.Auth == nil {
		return nil, fmt.Errorf("local auth service not configured")
	}

	// Register user locally
	authResp, err := l.svcCtx.Auth.Register(l.ctx, req.Email, req.Password, req.Name)
	if err != nil {
		l.Errorf("Registration failed for %s: %v", req.Email, err)
		return nil, err
	}

	l.Infof("User registered (local): %s", req.Email)

	// If user wants a paid plan, create checkout session
	var checkoutURL string
	plan := req.Plan
	if plan != "" && plan != "free" && l.svcCtx.Billing != nil {
		user, err := l.svcCtx.Auth.GetUserByEmail(l.ctx, req.Email)
		if err == nil {
			checkoutURL, _ = l.svcCtx.Billing.CreateCheckoutSession(l.ctx, user.ID, req.Email, plan)
		}
	}

	return &types.LoginResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    authResp.ExpiresAt.UnixMilli(),
		CheckoutUrl:  checkoutURL,
	}, nil
}
