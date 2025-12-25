package main

import (
	"database/sql"
	"log"
	"net/http"

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
	// 连接数据库
	var err error
	db, err = sql.Open("mysql", "wsb:admin1234@tcp(127.0.0.1:3306)/myapp")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库Ping失败:", err)
	}
	log.Println("数据库连接成功!")

	// 创建Gin路由
	r := gin.Default()

	// 配置CORS（允许跨域）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// 定义路由
	r.GET("/api/user", getUser)

	// 启动服务器
	log.Println("后端服务启动在 :8888")
	r.Run(":8888")
}

func getUser(c *gin.Context) {
	var user User

	// 查询数据库
	err := db.QueryRow("SELECT id, name, age FROM users LIMIT 1").Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}
