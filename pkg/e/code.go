package e

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	FailToParseXiaoYuZhouEpisode    = 10000
	FailToParseCowTransferShareLink = 10001

	UnsupportedUrl = 20000

	FailToSaveToRedis   = 30000
	FailToReadFromRedis = 30001
	KeyNotExists        = 30002

	FailToUnmarshal = 40000
)
