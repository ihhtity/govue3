package main

import (
	"context"
	v1 "govue3/api/v1"
	"govue3/config"
	"govue3/controllers"
	"govue3/db"
	_ "govue3/docs"
	"govue3/middleware"
	"govue3/services"
	"govue3/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title Govue3 API
// @version 1.0
// @description Govue3 项目API文档
// @BasePath /api/v1
func main() {
	// 初始化日志记录器
	logger, err := utils.InitLogger()
	if err != nil {
		log.Println(os.Stderr, "初始化日志记录器失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 初始化配置
	config.InitConfig("dev")
	config.WatchConfig()

	// 初始化数据库
	db.InitMySQL()
	db.InitRedis()
	db.InitMongoDB()

	// 创建Gin实例
	r := gin.Default()

	// 初始化限流中间件
	rateLimiter := middleware.NewRateLimiter(rate.Every(1*time.Second), 5) // 每秒1个请求，最多允许5个突发请求

	// 注册全局中间件（所有路由都会使用）
	r.Use(rateLimiter.Middleware())                    // 使用限流中间件
	r.Use(middleware.Logger())                         // 使用日志中间件
	r.Use(middleware.CORS())                           // 使用跨域中间件
	r.Use(middleware.RequestTimeout(30 * time.Second)) // 设置请求超时时间

	// 初始化服务和控制器
	userService := services.NewUserService()
	authService := services.NewAuthService()
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// 创建 BaseController 实例
	base := &controllers.BaseController{}

	// 路由分组
	api := r.Group("/api")
	{
		// Swagger 文档
		if config.Cfg.Swagger.Enabled { // 通过配置文件控制 Swagger 启用选项
			// 注册 Swagger UI 路由
			api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			log.Println("Swagger 文档已启用，可通过 /api/swagger/index.html 访问")
		} else {
			// 如果 Swagger 未启用，返回 404 或自定义提示信息
			api.GET("/swagger/*any", func(c *gin.Context) {
				base.Error(c, "Swagger 文档未启用")
			})
			log.Println("Swagger 文档未启用")
		}

		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			status := db.HealthCheck()
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"details": status,
			})
		})

		// Prometheus 监控
		api.GET("/metrics", gin.WrapH(promhttp.Handler()))

		// 不需要认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.POST("/reftoken", authController.RefreshToken)
		}

		// 需要认证的路由
		authorized := api.Group("/v1")
		{
			v1.RegisterRoutes(authorized, userController, authController)
		}
	}

	// 启动服务
	srv := &http.Server{
		Addr:    config.Cfg.Server.Port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	utils.GracefulShutdown(srv)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	log.Println("服务已关闭")
}
