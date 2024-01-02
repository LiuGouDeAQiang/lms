package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go_code/go_zero/api_study/user/api_jwt/internal/svc"
	"go_code/go_zero/api_study/user/api_jwt/internal/types"
	"go_code/go_zero/common/jwts"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (*types.UserInfoResponse, error) {
	auth := l.svcCtx.Config.Auth
	_, err := jwts.ParseToken(l.svcCtx.UserAgentMiddleware.JWTAuth, auth.AccessSecret) //验证token
	if err != nil {
		fmt.Printf("jwt解密失败%v132", l.svcCtx.UserAgentMiddleware.JWTAuth)
		return nil, err
	}
	a := &types.UserInfoResponse{
		UserId:   1,
		Username: "张飞",
		Addr:     "9",
		Id:       0,
	}

	return a, nil
}
