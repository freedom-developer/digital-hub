package servers

import (
	"context"
	"log"
	"myapp/config"
	"myapp/database"
	logger "myapp/log"
	"myapp/servers/music"
	"myapp/servers/user"

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
	musicService *music.MusicService
	us           *user.UserService
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

	us := user.NewUserService(nil, db, r)
	if us == nil {
		logger.ZFatal(ctx, "初始化用户服务失败", nil)
	}

	// 初始化音乐服务等
	musicService := music.NewMusicService(*ctx, &config.MusicConfig, db, r)
	if musicService == nil {
		logger.ZFatal(ctx, "初始化音乐服务失败", nil)
	}

	return &ServerManager{
		cfg:          config,
		db:           db,
		musicService: musicService,
		ctx:          ctx,
		r:            r,
		us:           us,
	}
}

func (srvMgr *ServerManager) StartAllServers() {
	// 启动所有服务器的逻辑
	if srvMgr.musicService != nil {
		// 启动音乐服务
		srvMgr.musicService.Start()

		// 启动用户服务
		srvMgr.us.Start()
	}

	// 添加外键约束
	// 为 user_id 添加外键约束
	srvMgr.db.Exec(`
		ALTER TABLE user_music 
		ADD CONSTRAINT fk_user_music_user 
		FOREIGN KEY (user_id) REFERENCES users(id) 
		ON DELETE CASCADE
	`)
	// 为 music_id 添加外键约束
	srvMgr.db.Exec(`
		ALTER TABLE user_music 
		ADD CONSTRAINT fk_user_music_music 
		FOREIGN KEY (music_id) REFERENCES musics(id) 
		ON DELETE CASCADE
	`)

	// 启动HTTP服务器
	addr := srvMgr.cfg.SrvConfig.Addr + ":" + srvMgr.cfg.SrvConfig.Port

	if err := srvMgr.r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
