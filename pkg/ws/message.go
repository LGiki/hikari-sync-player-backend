package ws

const (
	MessageEventLatestTime = "latestTime"
)

type WebsocketJsonMessage struct {
	Event string      `json:"type"`
	Data  interface{} `json:"data"`
}

type WebsocketBroadcastMessage struct {
	Data   []byte `json:"data"`
	RoomId string `json:"roomId"`
}
