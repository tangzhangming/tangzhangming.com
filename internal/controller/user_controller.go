package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Name  string `json:"name" label:"name" validate:"required" error:"姓名是必须的"`
	Email string `json:"email" validate:"email" label:"邮箱" `
	// Age   int    `json:"age" validate:"gte=0,lte=150"`
}

func User_view(c *gin.Context) {

}

func User_create(c *gin.Context) {

	// 从 Gin 上下文中获取验证器
	vd, exists := c.Get("vd")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validator not found"})
		return
	}

	// 使用验证器进行参数验证
	user := User{
		Name:  c.Query("name"),
		Email: c.Query("email"),
	}

	err := vd.(*validator.Validate).Struct(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
		return
	}

	c.String(200, "user_create")
}
