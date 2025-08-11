package xhs

// api 错误码 通用的错误码
// 业务错误码在业务内部定义

var (
	CodeError   = -1 // 通用服务器异常
	CodeSuccess = 0  // 通用接口请求错误
)

const (
	RoleSuperAdmin = "superAdmin"
)

const (
	CodeTokenError        = iota + 1000 // 登录失效或其他
	CodeRefreshTokenError               // 刷新token验证
	CodeCurrentLimiting                 // 限流
	CodeParamError                      // 参数错误
	CodeUserNotExist
	CodeMethodNotAllowed
	CodeNotFound
	CodeNoPermissions
	CodeInvoke = 8000
)
