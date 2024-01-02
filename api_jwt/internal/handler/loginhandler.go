package handler

import (
	"fmt"
	"go_code/go_zero/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go_code/go_zero/api_study/user/api_jwt/internal/logic"
	"go_code/go_zero/api_study/user/api_jwt/internal/svc"
	"go_code/go_zero/api_study/user/api_jwt/internal/types"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			fmt.Println("登录信息获取失败")
			response.Response(r, w, resp, err)
		}

		// 设置 Token 到请求头中
		r.Header.Set("Authorization", "Bearer "+resp)
		response.Response(r, w, resp, err)
	}
}
