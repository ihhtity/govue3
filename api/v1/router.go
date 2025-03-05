package v1

import (
	"govue3/controllers"
	"govue3/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, userController *controllers.UserController, authController *controllers.AuthController) {
	// 用户相关路由
	users := r.Group("/users")
	users.Use(middleware.JWTAuth())
	{
		users.GET("/", userController.GetUser)
		// 其他用户路由...
	}

	// 认证相关路由
	// auth := r.Group("/auth")
	// {
	// 	auth.POST("/login", authController.Login)
	// 	auth.POST("/register", authController.Register)
	// 	// 其他认证路由...
	// }
}
