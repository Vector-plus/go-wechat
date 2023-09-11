package re

const (
	SUCCSE = 200

	ERROR           = 500
	ERROR_PARAM     = 501
	ERROR_DATA_EIXT = 502

	ERROR_TOKEN_EIXT    = 1001
	ERROR_TOKEN_WRONG   = 1002
	ERROR_TOKEN_EXPIRED = 1003
	ERROR_TOKEN_FAIL    = 1004

	ERROR_USER_EIXT        = 1020
	ERROR_USER_PASSWORD    = 1021
	ERROR_USER_ROLE        = 1022
	ERROR_USER_PASSWORD_RE = 1023
	ERROR_USER_PHONE       = 1024

	ERROR_APPLICATION_EIXT = 1080

	ERROR_WS_CONN = 2000

	ERROR_FILE_READ   = 2500
	ERROR_FILE_CREATE = 2501
)

var codeMsg = map[int]string{
	SUCCSE:                 "OK",
	ERROR:                  "FAIL",
	ERROR_TOKEN_EIXT:       "token不存在",
	ERROR_TOKEN_WRONG:      "token格式错误",
	ERROR_TOKEN_EXPIRED:    "token已过期",
	ERROR_TOKEN_FAIL:       "token无效,重新登录",
	ERROR_USER_EIXT:        "该用户不存在,重新输入",
	ERROR_USER_PASSWORD:    "密码错误,重新输入",
	ERROR_USER_ROLE:        "用户没有权限",
	ERROR_USER_PASSWORD_RE: "用户输入密码不一致",
	ERROR_USER_PHONE:       "用户手机号已存在",

	ERROR_APPLICATION_EIXT: "该申请不存在",
	ERROR_PARAM:            "参数错误，请重新输入",
	ERROR_DATA_EIXT:        "该数据已经存在",

	ERROR_WS_CONN: "WebSocket连接错误",

	ERROR_FILE_READ:   "文件读取错误",
	ERROR_FILE_CREATE: "文件创建错误",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
