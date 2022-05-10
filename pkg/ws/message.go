package ws

import (
	"encoding/json"
	"time"
)

// WebsocketJsonMessage is the JSON format used in websocket
type WebsocketJsonMessage struct {
	Event  string          `json:"event"`
	UserId string          `json:"userId"`
	Data   json.RawMessage `json:"data"`
}

type PlayState struct {
	// PlayTime in seconds
	PlayTime  int64   `json:"playTime"`
	IsPlaying bool    `json:"isPlaying"`
	PlaySpeed float32 `json:"playSpeed"`
	// Timestamp is the milliseconds elapsed since the UNIX epoch
	Timestamp int64 `json:"timestamp"`
}

// GetCurrentPlayTime returns the play time at current timestamp
func (ps *PlayState) GetCurrentPlayTime() int64 {
	currentTimestamp := time.Now().UnixMilli()
	playSpeed := int64(ps.PlaySpeed * 100)
	return (playSpeed*(currentTimestamp-ps.Timestamp))/(1000*100) + ps.PlayTime
}

// GetCurrentPlayState returns a new PlayState where PlayState.PlayTime is at the current timestamp.
func (ps *PlayState) GetCurrentPlayState() *PlayState {
	if ps == nil {
		return nil
	}
	return &PlayState{
		PlayTime:  ps.GetCurrentPlayTime(),
		IsPlaying: ps.IsPlaying,
		PlaySpeed: ps.PlaySpeed,
		Timestamp: ps.Timestamp,
	}
}

// WebsocketBroadcastMessage is used to send messages to the Hub.Broadcast
type WebsocketBroadcastMessage struct {
	Data   []byte `json:"data"`
	RoomId string `json:"roomId"`
}
