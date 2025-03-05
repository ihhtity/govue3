package controllers

import (
	"govue3/services"
	"govue3/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	BaseController
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary 用户登录
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} string "登录成功"
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	// 获取登录参数
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.Error(c, "获取参数失败!")
		return
	}

	// 调用服务层验证登录
	token, user, err := a.authService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		a.Error(c, err.Error())
		return
	}

	data := gin.H{
		"token": token,
		"user":  user,
	}
	a.Success(c, "登陆成功!", data)
}

// @Summary 用户注册
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} string "注册成功"
// @Router /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	// 获取注册参数
	var registerRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		a.Error(c, "获取参数失败!")
		return
	}

	// 调用服务层进行注册
	err := a.authService.Register(
		registerRequest.Username,
		registerRequest.Password,
	)
	if err != nil {
		a.Error(c, err.Error())
		return
	}

	a.Success(c, "注册成功", nil)
}

// @Summary 刷新Token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} string "Token刷新成功"
// @Router /auth/refresh [post]
func (a *AuthController) RefreshToken(c *gin.Context) {
	// 从请求头获取refresh token
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		a.Error(c, "缺少刷新令牌")
		return
	}

	// 调用服务层刷新token
	newToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		a.Error(c, err.Error())
		return
	}

	a.Success(c, "Token刷新成功", newToken)
}
