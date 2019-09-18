package error

var MsgFlags = map[int]string{
	SUCCESS:              "ok",
	ERROR:                "服务器错误,请稍后重试",
	INVALID_PARAMS:       "请求参数错误",
	ERROR_NOT_EXIST_USER: "用户不存在",
}

// GetMsg 通过code码获取msg
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
