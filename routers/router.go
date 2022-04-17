package routers

import (
	"github.com/gin-gonic/gin"
	"hikari_sync_player/middleware/cors"
	"hikari_sync_player/pkg/ws"
	v1 "hikari_sync_player/routers/api/v1"
)

func SetupRouter() *gin.Engine {
	hub := ws.NewHub()
	go hub.Run()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Cors())
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/room", v1.CreateRoom)
		apiV1.GET("/room/:id", func(context *gin.Context) {
			v1.RoomSync(context, hub)
		})
	}
	return router
}
