package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-one-piece/router"
	"go-one-piece/util"
	"go-one-piece/util/conf"
	"go-one-piece/util/logging"
	"net/http"
	"time"
)

func init() {
	conf.Setup()
	logging.Setup()
}

func main() {
	reload := make(chan int, 1)
	conf.OnConfigChange(func() { reload <- 1 })
	startServer()
	for {
		select {
		case <-reload:
			util.Reset()
			logging.Get().Infoln("OnConfigChange")
		}
	}
}

func startServer() {
	timeLocal := time.FixedZone("CST", 8*3600)
	time.Local = timeLocal
	gin.SetMode(conf.Config.Server.RunMode)
	httpPort := fmt.Sprintf(":%d", conf.Config.Server.HttpPort)
	server := &http.Server{
		Addr:           httpPort,
		Handler:        router.InitRouter(),
		ReadTimeout:    conf.Config.Server.ReadTimeout,
		WriteTimeout:   conf.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}
