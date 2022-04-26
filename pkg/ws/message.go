package ws

type WebsocketJsonMessage struct {
	Event  string      `json:"event"`
	UserId string      `json:"userId"`
	Data   interface{} `json:"data"`
}

type WebsocketBroadcastMessage struct {
	Data   []byte `json:"data"`
	RoomId string `json:"roomId"`
}
