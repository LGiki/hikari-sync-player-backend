package e

var MsgFlags = map[int]string{
	SUCCESS:       "success",
	ERROR:         "fail",
	InvalidParams: "invalid params",

	FailToParseXiaoYuZhouEpisode:    "fail to parse xiaoyuzhou episode url",
	FailToParseCowTransferShareLink: "fail to parse cow transfer share link",
	UnsupportedUrl:                  "unsupported url",
	FailToSaveToRedis:               "fail to save to redis",
	FailToReadFromRedis:             "fail to read from redis",
	KeyNotExists:                    "key not exists",
	FailToUnmarshal:                 "fail to unmarshal",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
