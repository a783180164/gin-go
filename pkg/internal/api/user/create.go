// pkg/internal/api/login/login.go
package user

import (
	// "gin-go/pkg/jwt"
	"gin-go/pkg/code"
	"gin-go/pkg/crypto"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/user"
	"gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// LoginRequest 登录请求体
type CreateRequest struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatResponse struct {
	ID int32 `json:"id"`
}

func (h *handler) RegisterUser(c *gin.Context) {

	core := Core.NewContext(c)
	req := new(CreateRequest)
	// 解析 JSON 请求体
	if err := core.ShouldBindJSON(req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"err":  err.Error(),
			"user": req.User,
		}).Info("Attempting to register user")
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.CreateUserParam,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	createData := new(user.CreateUserData)
	createData.User = req.User
	password, err := crypto.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.HashPasswordError,
			Data:    nil,
			Message: err.Error(),
		})
	}
	createData.Password = password

	// 创建用户
	id, err := h.userService.CreateUser(createData)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.CreatHasUserError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code: code.Success,
		Data: &CreatResponse{
			ID: id,
		},
		Message: code.Text(code.Success),
	})
}
