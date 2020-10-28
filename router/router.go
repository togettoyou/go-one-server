package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "go-one-server/docs"
	"go-one-server/handler/v1/examples"
	"go-one-server/router/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.Logger())
	//debug模式开启性能分析
	pprof.Register(r)
	//swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//api路由分组v1版本
	apiV1 := r.Group("/api/v1")
	initExamplesRouter(apiV1)
	return r
}

func initExamplesRouter(api *gin.RouterGroup) {
	examplesRouterGroup := api.Group("/examples")
	{
		examplesRouterGroup.GET("/get", examples.GetExamples)
	}
}
