package servers

import (
	"context"
	"log"
	"myapp/config"
	"myapp/database"
	logger "myapp/log"
	"myapp/servers/music"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServerManager struct {
	// Add fields for managing servers if needed
	cfg          *config.Config
	db           *gorm.DB
	ctx          *context.Context
	r            *gin.Engine
	rg           *gin.RouterGroup
	musicService *music.MusicService
}

func NewServerManager(ctx *context.Context, config *config.Config) *ServerManager {
	// 初始化数据库
	db, err := database.InitDB(&config.DbConfig)
	if err != nil {
		logger.ZFatal(ctx, "初始化数据库失败", err)
	}

	r := gin.Default()

	// // 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	rg := r.Group("/api")

	// 初始化音乐服务等
	musicService := music.NewMusicService(*ctx, &config.MusicConfig, db, rg)
	if musicService == nil {
		logger.ZFatal(ctx, "初始化音乐服务失败", nil)
	}

	return &ServerManager{
		cfg:          config,
		db:           db,
		musicService: musicService,
		ctx:          ctx,
		r:            r,
		rg:           rg,
	}
}

func (srvMgr *ServerManager) StartAllServers() {
	// 启动所有服务器的逻辑
	if srvMgr.musicService != nil {
		// 启动音乐服务
		srvMgr.musicService.Start()
	}

	// 启动HTTP服务器
	addr := srvMgr.cfg.SrvConfig.Addr + ":" + srvMgr.cfg.SrvConfig.Port

	if err := srvMgr.r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
