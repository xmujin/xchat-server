package controller

import (
	"xchat-server/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	// 包含 Service 引用
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// 建议命名：简洁明了
func (u *UserController) Register(c *gin.Context) {
	// 获取用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 获取token
	token := u.userService.Register(username, password)
	c.JSON(200, gin.H{
		"token": token,
	})
}

func (u *UserController) Login(c *gin.Context) {
	// 获取用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 获取token
	token := u.userService.Register(username, password)
	c.JSON(200, gin.H{
		"token": token,
	})
}

func (u *UserController) Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "我真的服了你了",
	})
}
