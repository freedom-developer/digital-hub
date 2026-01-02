package user

import (
	logger "myapp/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserConfig struct {
	// Add configuration fields here
}

type UserService struct {
	// Add service methods here
	cfg *UserConfig
	db  *gorm.DB
	rg  *gin.RouterGroup
}

func NewUserService(cfg *UserConfig, db *gorm.DB, r *gin.Engine) *UserService {
	err := db.AutoMigrate(&User{})
	if err != nil {
		logger.ZError(nil, "数据库自动迁移失败", err)
		return nil
	}

	rg := r.Group("/user")

	return &UserService{
		cfg: cfg,
		db:  db,
		rg:  rg,
	}
}

func (us *UserService) Start() {
	// Implementation for retrieving user by ID
	us.RegisterRouters()
}
