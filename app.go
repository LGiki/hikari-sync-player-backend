package main

import (
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
}

func main() {
	hub := ws.NewHub()
	go hub.Run()
	router := routers.SetupRouter()
	server := &http.Server{
		Addr:           "0.0.0.0:12312",
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
