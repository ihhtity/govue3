package controllers

import (
	"govue3/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	// idStr = "1"
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Error(ctx, "无效的用户ID!")
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		c.Error(ctx, "用户不存在!")
		return
	}

	c.Success(ctx, "获取成功!", user)
}
