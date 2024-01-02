package svc

import (
	"go_code/go_zero/api_study/user/api_jwt/internal/config"
	"go_code/go_zero/api_study/user/api_jwt/internal/middleware"
)

type ServiceContext struct {
	Config config.Config

	UserAgentMiddleware *middleware.UserAgentMiddleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UserAgentMiddleware: middleware.NewUserAgentMiddleware(),
	}
}
