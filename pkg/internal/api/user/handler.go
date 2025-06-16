package user

import (
	"gin-go/pkg/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	logger      *logrus.Logger
	userService user.Service
}

func New(logger *logrus.Logger, db *gorm.DB) Handler {
	return &handler{
		logger:      logger,
		userService: user.New(db),
	}
}

func (h *handler) i() {}
