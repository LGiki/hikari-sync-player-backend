package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hikari_sync_player/pkg/app"
	"hikari_sync_player/pkg/cowtransfer_parser"
	"hikari_sync_player/pkg/e"
	"hikari_sync_player/pkg/gredis"
	"hikari_sync_player/pkg/logging"
	"hikari_sync_player/pkg/ws"
	"hikari_sync_player/pkg/xiaoyuzhou_parser"
	"hikari_sync_player/util"
	"net/http"
	"strings"
)

var webSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	RoomTypePodcast = "podcast"
	RoomTypeVideo   = "video"
)

const RedisExpireTime = 24 * 60 * 60

type CreateRoomRequestBody struct {
	Url string `json:"url"`
}

type CreateRoomResponseBody struct {
	Type   string `json:"type"`
	RoomId string `json:"roomId"`
}

type RoomDetail struct {
	Type      string      `json:"type"`
	MediaData interface{} `json:"mediaData"`
}

func CreateRoom(context *gin.Context) {
	appG := app.Gin{Context: context}
	var createRoomRequestBody CreateRoomRequestBody
	err := context.BindJSON(&createRoomRequestBody)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	// parse xiaoyuzhou episode link
	if strings.Contains(createRoomRequestBody.Url, "xiaoyuzhoufm.com") {
		episode, err := xiaoyuzhou_parser.ParseEpisode(createRoomRequestBody.Url)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.FailToParseXiaoYuZhouEpisode, nil)
			return
		}
		roomId := util.GenerateUUID()
		roomDetail := &RoomDetail{
			Type:      RoomTypePodcast,
			MediaData: episode,
		}
		err = gredis.Set(roomId, roomDetail, RedisExpireTime)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.FailToSaveToRedis, nil)
			return
		}
		createRoomResponseBody := CreateRoomResponseBody{
			Type:   RoomTypePodcast,
			RoomId: roomId,
		}
		appG.Response(http.StatusOK, e.SUCCESS, createRoomResponseBody)
		return
	}

	// parse cow transfer share link
	if strings.Contains(createRoomRequestBody.Url, "cowtransfer.com") {
		video, err := cowtransfer_parser.ParseShareLink(createRoomRequestBody.Url)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.FailToParseCowTransferShareLink, nil)
			return
		}
		roomId := util.GenerateUUID()
		roomDetail := &RoomDetail{
			Type:      RoomTypeVideo,
			MediaData: video,
		}
		err = gredis.Set(roomId, roomDetail, RedisExpireTime)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.FailToSaveToRedis, nil)
			return
		}
		createRoomResponseBody := CreateRoomResponseBody{
			Type:   RoomTypeVideo,
			RoomId: roomId,
		}
		appG.Response(http.StatusOK, e.SUCCESS, createRoomResponseBody)
		return
	}

	appG.Response(http.StatusBadRequest, e.UnsupportedUrl, nil)
}

func GetRoomDetail(context *gin.Context) {
	appG := app.Gin{Context: context}
	roomId := context.Param("roomId")
	if !gredis.Exists(roomId) {
		appG.Response(http.StatusBadRequest, e.KeyNotExists, nil)
		return
	}
	roomDetailBytes, err := gredis.Get(roomId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.FailToReadFromRedis, nil)
		return
	}
	var roomDetail RoomDetail
	err = json.Unmarshal(roomDetailBytes, &roomDetail)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.FailToUnmarshal, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, roomDetail)
}

func RoomSync(context *gin.Context, hub *ws.Hub) {
	conn, err := webSocketUpgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logging.Error(err)
		return
	}
	roomId := context.Param("roomId")
	client := ws.NewClient(hub, conn, roomId)
	hub.Register <- client
	go client.WriteLoop()
	go client.ReadLoop()
}
