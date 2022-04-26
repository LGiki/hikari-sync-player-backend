package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hikari_sync_player/pkg/gredis"
	"hikari_sync_player/pkg/logging"
	"hikari_sync_player/pkg/setting"
	"hikari_sync_player/pkg/ws"
	"hikari_sync_player/routers"
	"log"
	"net/http"
	"time"
)

func init() {
	err := setting.Setup()
	if err != nil {
		log.Fatalf("init setting err: %v", err)
	}
	err = logging.Setup()
	if err != nil {
		log.Fatalf("init logging err: %v", err)
	}
	err = gredis.Setup()
	if err != nil {
		logging.Fatal(fmt.Sprintf("init redis err: %v", err))
	}
}

func main() {
	hub := ws.NewHub()
	go hub.Run()
	gin.SetMode(setting.GlobalSettings.App.RunningMode)
	router := routers.SetupRouter()
	addr := fmt.Sprintf("%s:%d", setting.GlobalSettings.App.Host, setting.GlobalSettings.App.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		logging.Error(err)
	}
}
