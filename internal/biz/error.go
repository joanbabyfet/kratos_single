package biz

import "github.com/go-kratos/kratos/v2/errors"

//code 在框架里被当成 HTTP status / gRPC status / protobuf status 所以不支持-1, -2, -3
var (
	// ================= 系统级错误 =================
	ErrSystem 		= errors.New(500, "-1211", "系统错误")
	// ================= 业务错误 =================
	ErrParam    	= errors.New(400, "-1", "参数错误")
	ErrInvalidID    = errors.New(400, "-1", "ID不能为空")
)