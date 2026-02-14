package handler

import (
	"context"
	"net/http"
	"time"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const version = "1.0.0"

type HealthCheckLogic struct {
	logx.Logger
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
			Logger: logx.WithContext(r.Context()),
			ctx:    r.Context(),
			svcCtx: svcCtx,
		}
		resp, err := l.HealthCheck()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
