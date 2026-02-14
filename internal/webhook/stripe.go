package webhook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"gosaas/internal/svc"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"log/slog"
)

// StripeHandler creates a raw HTTP handler for Stripe webhooks
// This must be registered directly to access raw body for signature verification
func StripeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read raw body for signature verification
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read webhook body", "error", err)
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		// Get Stripe signature header
		sig := r.Header.Get("Stripe-Signature")
		if sig == "" {
			slog.Error("Missing Stripe-Signature header")
			http.Error(w, "Missing signature", http.StatusBadRequest)
			return
		}

		// Verify webhook signature
		webhookSecret := svcCtx.Config.Stripe.WebhookSecret
		if webhookSecret == "" {
			slog.Error("Stripe webhook secret not configured")
			http.Error(w, "Webhook not configured", http.StatusInternalServerError)
			return
		}

		event, err := webhook.ConstructEventWithOptions(body, sig, webhookSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
		if err != nil {
			slog.Error("Failed to verify webhook signature", "error", err)
			http.Error(w, "Invalid signature", http.StatusBadRequest)
			return
		}

		// Handle the event
		if err := handleStripeEvent(svcCtx, &event); err != nil {
			slog.Error("Failed to handle webhook event", "type", event.Type, "error", err)
			http.Error(w, "Failed to handle event", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"received": true}`))
	}
}

// handleStripeEvent processes different Stripe event types
func handleStripeEvent(svcCtx *svc.ServiceContext, event *stripe.Event) error {
	ctx := context.Background()

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return err
		}
		return handleCheckoutCompleted(svcCtx, ctx, &session)

	case "customer.subscription.created", "customer.subscription.updated":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return err
		}
		return handleSubscriptionUpdated(svcCtx, ctx, &subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return err
		}
		return handleSubscriptionDeleted(svcCtx, ctx, &subscription)

	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return err
		}
		slog.Info("Invoice payment succeeded", "invoice_id", invoice.ID)

	case "invoice.payment_failed":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return err
		}
		slog.Warn("Invoice payment failed", "invoice_id", invoice.ID, "customer_id", invoice.Customer.ID)

	default:
		slog.Info("Unhandled Stripe event type", "type", event.Type)
	}

	return nil
}

// handleCheckoutCompleted processes successful checkout sessions
func handleCheckoutCompleted(svcCtx *svc.ServiceContext, ctx context.Context, session *stripe.CheckoutSession) error {
	customerID := ""
	if session.Customer != nil {
		customerID = session.Customer.ID
	}
	subscriptionID := ""
	if session.Subscription != nil {
		subscriptionID = session.Subscription.ID
	}

	slog.Info("Checkout completed", "session_id", session.ID, "customer_id", customerID, "subscription_id", subscriptionID)

	// Get user ID from metadata
	userID, ok := session.Metadata["user_id"]
	if !ok || userID == "" {
		slog.Info("No user_id in checkout session metadata")
		return nil
	}

	// Update subscription in database
	if svcCtx.Billing != nil {
		return svcCtx.Billing.HandleCheckoutCompleted(ctx, customerID, subscriptionID, userID)
	}

	return nil
}

// handleSubscriptionUpdated processes subscription updates
func handleSubscriptionUpdated(svcCtx *svc.ServiceContext, ctx context.Context, subscription *stripe.Subscription) error {
	slog.Info("Subscription updated", "subscription_id", subscription.ID, "status", subscription.Status)

	if svcCtx.Billing != nil {
		// Get period from subscription items (Stripe API 2025-03-31 moved these fields)
		var periodStart, periodEnd int64
		if len(subscription.Items.Data) > 0 {
			item := subscription.Items.Data[0]
			periodStart = item.CurrentPeriodStart
			periodEnd = item.CurrentPeriodEnd
		}

		return svcCtx.Billing.HandleSubscriptionUpdated(
			ctx,
			subscription.ID,
			string(subscription.Status),
			periodStart,
			periodEnd,
			subscription.CancelAtPeriodEnd,
		)
	}

	return nil
}

// handleSubscriptionDeleted processes subscription cancellations
func handleSubscriptionDeleted(svcCtx *svc.ServiceContext, ctx context.Context, subscription *stripe.Subscription) error {
	slog.Info("Subscription deleted", "subscription_id", subscription.ID)

	if svcCtx.Billing != nil {
		return svcCtx.Billing.HandleSubscriptionDeleted(ctx, subscription.ID)
	}

	return nil
}
