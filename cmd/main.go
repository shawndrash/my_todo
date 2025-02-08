package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"my_todo/internal/database"
	"my_todo/internal/handler"
	"my_todo/internal/router"
)

func main() {
	// 初始化数据库连接
	db, err := database.NewDatabase("todo.db") // 使用 "todo.db" 作为数据库文件名
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 使用 defer 关闭数据库连接

	// 创建 Gin 引擎
	r := gin.Default()

	// 创建处理器，并将数据库连接传递给处理器
	todoHandler := handler.NewTodoHandler(db)

	// 初始化路由
	router.SetupRoutes(r, todoHandler)

	// 载入静态文件 (如果需要)
	r.Static("/static", "./static")

	// 启动 Web 服务器
	log.Fatal(r.Run(":8080"))
}

