package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"learningAssistant-backend/config"
	"learningAssistant-backend/database"
	_ "learningAssistant-backend/docs" // 导入生成的 docs
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/routes"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	database.InitDatabase()

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	log.Printf("Server mode: %s", config.AppConfig.Server.Mode)

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
