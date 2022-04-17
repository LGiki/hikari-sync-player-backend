package e

var MsgFlags = map[int]string{
	SUCCESS:       "success",
	ERROR:         "fail",
	InvalidParams: "invalid params",

	FailToParseXiaoYuZhouEpisode: "fail to parse xiaoyuzhou episode url",
	UnsupportedUrl:               "unsupported url",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
