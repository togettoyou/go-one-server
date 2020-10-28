package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-one-server/model"
	"go-one-server/router"
	"go-one-server/util"
	"go-one-server/util/conf"
	"go-one-server/util/logger"
	"go-one-server/util/times"
	"go-one-server/util/validator"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func init() {
	conf.Setup()
	logger.Setup()
	validator.Setup()
	model.Setup()
}

// @title go-one-server
// @version 1.0
// @description 基于Gin进行快速构建RESTful API 服务的脚手架
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	reload := make(chan int, 1)
	conf.OnConfigChange(func() { reload <- 1 })
	startServer()
	for {
		select {
		case <-reload:
			util.Reset()
		}
	}
}

func startServer() {
	timeLocal := time.FixedZone("CST", 8*3600)
	time.Local = timeLocal
	zap.L().Info(time.Now().Format(times.TimeFormat))
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
	fmt.Printf(`swagger文档地址：http://localhost%s/swagger/index.html
   ____   ____             ____   ____   ____             ______ ______________  __ ___________ 
  / ___\ /  _ \   ______  /  _ \ /    \_/ __ \   ______  /  ___// __ \_  __ \  \/ // __ \_  __ \
 / /_/  >  <_> ) /_____/ (  <_> )   |  \  ___/  /_____/  \___ \\  ___/|  | \/\   /\  ___/|  | \/
 \___  / \____/           \____/|___|  /\___  >         /____  >\___  >__|    \_/  \___  >__|   
/_____/                              \/     \/               \/     \/                 \/       

`, httpPort)
}
