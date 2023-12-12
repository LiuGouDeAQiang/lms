package tools

import "fmt"

var (
	OK1            = ECode{Code: 0, Message: "借书成功"}
	OK             = ECode{Code: 0, Message: "成功"}
	OK2            = ECode{Code: 0, Message: "还书成功"}
	TouristLogin   = ECode{Code: 10000, Message: "游客登录"}
	NotLogin       = ECode{Code: 10001, Message: "用户未登录"}
	ParamErr       = ECode{Code: 10002, Message: "参数错误"}
	UserErr        = ECode{Code: 10003, Message: "账号或密码错误"}
	BorrowErr      = ECode{Code: 10004, Message: "借书失败,数据库查询失败"}
	BorrowErr2     = ECode{Code: 10005, Message: "借书失败,请先登录"}
	ReturnErr      = ECode{Code: 10006, Message: "还书失败"}
	ReturnErr2     = ECode{Code: 10007, Message: "还书失败，请先登录"}
	GetImgErr      = ECode{Code: 10008, Message: "获取封面失败"}
	GetImgRedisErr = ECode{Code: 10008, Message: "获取封面失败"}
	EmailErr       = ECode{Code: 10009, Message: "获取封面失败"}
	ConflictErr    = ECode{Code: 10010, Message: "乐观锁冲突"}
)

type ECode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("code:%d,message:%s", e.Code, e.Message)
}
