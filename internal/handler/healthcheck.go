package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gosaas/internal/httpx"
	"gosaas/internal/svc"
	"gosaas/internal/types"
)

const version = "1.0.0"

type HealthCheckLogic struct {
	logger *slog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *HealthCheckLogic) HealthCheck() (resp *types.HealthResponse, err error) {
	return &types.HealthResponse{
		Status:    "healthy",
		Version:   version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func HealthCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &HealthCheckLogic{
			logger: slog.Default(),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.HealthCheck()
		if err != nil {
			httpx.ErrorResponse(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
