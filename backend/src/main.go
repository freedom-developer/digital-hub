package main

import (
	"context"
	"myapp/config"
	logger "myapp/log"
	"myapp/servers"
)

// var srvManager *servers.ServerManager

func main() {
	cfg := config.LoadConfig()

	if logger.InitLogger(&cfg.LogConfig) == nil {
		panic("初始化日志失败")
	}

	ctx := context.Background()
	ctx = *logger.SetTagInContext(&ctx, "main")

	srvManager := servers.NewServerManager(&ctx, cfg)
	if srvManager == nil {
		logger.ZFatal(&ctx, "初始化服务器管理器失败", nil)
	}
	srvManager.StartAllServers()

}
