package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func main() {
	var err error
	
	dsn := "wsb:admin1234@tcp(127.0.0.1:3306)/myapp?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("数据库Ping失败:", err)
	}
	log.Println("数据库连接成功!")

	r := gin.Default()

	// ⭐⭐⭐ 关键：CORS 配置（必须在路由之前）
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://172.16.85.131:5173",  // ⭐ 你的实际 IP
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// 路由
	r.GET("/api/user", getUser)

	log.Println("后端服务启动在 0.0.0.0:8888")
	r.Run("0.0.0.0:8888")
}

func getUser(c *gin.Context) {
	var user User
	err := db.QueryRow("SELECT id, name, age FROM users LIMIT 1").Scan(
		&user.ID, &user.Name, &user.Age,
	)
	if err != nil {
		log.Println("查询错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, user)
}
