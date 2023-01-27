package main

import (
	"flag"
	"go-pano/config"
	router_v1 "go-pano/router"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()

	// 顯示監聽地址
	addr := flag.String("addr", config.Server.Addr, "Address to listen and serve")
	flag.Parse()

	// 關閉顏色 增加效能
	if config.Server.Mode == gin.ReleaseMode {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	// 使用預設 logger recovery
	app := gin.Default()

	// 靜態檔案
	app.Static("/static", config.Server.StaticDir)
	app.Static("/docs", "./docs")

	// 配置檔案上傳大小
	app.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20

	// 路由添加
	router_v1.NewRouter(app)

	// 監聽
	app.Run(*addr)
}
