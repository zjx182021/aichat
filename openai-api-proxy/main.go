package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"openai-api-proxy/pkg/config"
	"openai-api-proxy/pkg/log"
	"openai-api-proxy/routers"
)

var (
	configFile = flag.String("config", "dev.config.yaml", "")
)

func main() {
	flag.Parse()
	config.InitConfig(*configFile)
	conf := config.GetConfig()

	gin.SetMode(conf.Http.Mode)

	//初始化日志组件
	log.SetOutput(log.GetRotateWriter(conf.Log.LogPath))
	log.SetLevel(conf.Log.Level)
	log.SetPrintCaller(true)

	r := gin.Default()
	routers.InitRouters(r)
	r.Run(fmt.Sprintf("%s:%d", conf.Http.Host, conf.Http.Port))
}
