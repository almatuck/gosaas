package subscription

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gosaas/internal/auth"
	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

type ListBillingHistoryLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *ListBillingHistoryLogic) ListBillingHistory(req *types.ListBillingHistoryRequest) (resp *types.ListBillingHistoryResponse, err error) {
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("levee service not configured")
	}

	// Get email from JWT context
	email, err := auth.GetEmailFromContext(l.ctx)
	if err != nil {
		slog.Error("Failed to get email from context", "error", err)
		return nil, err
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	// Get invoices from Levee SDK
	invoicesResp, err := l.svcCtx.Levee.Customers.ListCustomerInvoices(l.ctx, email, pageSize)
	if err != nil {
		slog.Error("Failed to get invoices", "email", email, "error", err)
		return nil, err
	}

	// Convert to response format
	records := make([]types.BillingRecord, 0, len(invoicesResp.Invoices))
	for _, inv := range invoicesResp.Invoices {
		records = append(records, types.BillingRecord{
			Id:               inv.ID,
			Amount:           int(inv.AmountDue),
			Currency:         inv.Currency,
			Status:           inv.Status,
			Description:      inv.Description,
			InvoiceDate:      inv.CreatedAt,
			PaidAt:           inv.PaidAt,
			InvoicePdfUrl:    inv.InvoicePdfUrl,
			HostedInvoiceUrl: inv.HostedUrl,
		})
	}

	return &types.ListBillingHistoryResponse{
		Records:    records,
		TotalCount: len(records),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// ListBillingHistoryHandler handles requests to get the billing history.
func ListBillingHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListBillingHistoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorResponse(w, err)
			return
		}

		l := &ListBillingHistoryLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.ListBillingHistory(&req)
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
