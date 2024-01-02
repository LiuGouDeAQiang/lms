package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response http返回
func Response(r *http.Request, w http.ResponseWriter, res any, err error) {
	if err != nil {
		//成功返回
		r := Body{
			Code: 10086,
			Msg:  "错误",
			Data: nil,
		}
		httpx.WriteJson(w, http.StatusOK, r)
		return
	}
	body := &Body{
		Code: 0,
		Msg:  "成功",
		Data: res,
	}
	httpx.WriteJson(w, http.StatusOK, body)

}
