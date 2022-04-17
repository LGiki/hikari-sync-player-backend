package main

import (
	"hikari_listen/pkg/logging"
	"hikari_listen/pkg/setting"
	"hikari_listen/pkg/ws"
	"hikari_listen/routers"
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
