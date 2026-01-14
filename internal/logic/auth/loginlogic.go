package auth

import (
	"context"
	"fmt"
	"time"

	"gosaas/internal/svc"
	"gosaas/internal/types"

	levee "github.com/almatuck/levee-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// Use local auth when Levee is disabled
	if l.svcCtx.UseLocal() {
		return l.loginLocal(req)
	}

	// Use Levee when enabled
	if l.svcCtx.Levee == nil {
		return nil, fmt.Errorf("auth service not configured")
	}

	// Authenticate via Levee SDK
	authResp, err := l.svcCtx.Levee.Auth.Login(l.ctx, &levee.SDKLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		l.Errorf("Login failed for %s: %v", req.Email, err)
		return nil, err
	}

	// Parse expiry time
	expiresAt, _ := time.Parse(time.RFC3339, authResp.ExpiresAt)

	l.Infof("User logged in: %s", req.Email)

	return &types.LoginResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    expiresAt.UnixMilli(),
	}, nil
}

// loginLocal handles login with local SQLite auth
func (l *LoginLogic) loginLocal(req *types.LoginRequest) (*types.LoginResponse, error) {
	if l.svcCtx.Auth == nil {
		return nil, fmt.Errorf("local auth service not configured")
	}

	authResp, err := l.svcCtx.Auth.Login(l.ctx, req.Email, req.Password)
	if err != nil {
		l.Errorf("Login failed for %s: %v", req.Email, err)
		return nil, err
	}

	l.Infof("User logged in (local): %s", req.Email)

	return &types.LoginResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    authResp.ExpiresAt.UnixMilli(),
	}, nil
}
