package user

import (
	"gin-go/pkg/code"
	"gin-go/pkg/crypto"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/user"
	"gin-go/pkg/jwt"
	"gin-go/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       int32  `json:"id"`
	User     string `json:"user"`
	UUID     string `json:"uuid"`
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (h *handler) Login(c *gin.Context) {
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
	userData := new(user.UserData)
	userData.User = req.User
	userData.Password = req.Password

	info, err := h.userService.Login(userData)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.HashPasswordError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.HashPasswordError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	if info == nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.HashPasswordError,
			Data:    nil,
			Message: "没有用户",
		})
		return
	}
	if crypto.CheckPasswordHash(userData.Password, info.Password) == false {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.HashPasswordError,
			Data:    nil,
			Message: "密码不正确",
		})
		return
	}

	res := &LoginResponse{
		ID:       info.ID,
		UUID:     info.UUID,
		User:     info.User,
		NickName: info.Nickname,
		UserName: info.Username,
		Avatar:   info.Avatar,
	}

	// 生成 JWT token
	jwtCfg := jwtmidd.DefaultJWTConfig()
	claims := jwt.MapClaims{
		"user_id":  info.ID,
		"username": info.Username,
	}
	token, err := jwtCfg.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.AuthorizationError,
			Data:    nil,
			Message: "生成 token 失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success,
		"data":    res,
		"token":   token,
		"message": "成功",
	})

}
