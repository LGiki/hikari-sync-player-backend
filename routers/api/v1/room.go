package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hikari_listen/pkg/app"
	"hikari_listen/pkg/e"
	"hikari_listen/pkg/logging"
	"hikari_listen/pkg/ws"
	"hikari_listen/pkg/xiaoyuzhou_parser"
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

type CreateRoomRequestBody struct {
	Url string `json:"url"`
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
		appG.Response(http.StatusOK, e.SUCCESS, episode)
		return
	}

	appG.Response(http.StatusBadRequest, e.UnsupportedUrl, nil)
}

func RoomSync(context *gin.Context, hub *ws.Hub) {
	conn, err := webSocketUpgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		logging.Error(err)
		return
	}
	roomId := context.Param("id")
	client := ws.NewClient(hub, conn, roomId)
	hub.Register <- client
	go client.WriteLoop()
	go client.ReadLoop()
}
