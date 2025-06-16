package weather

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler interface {
	i()
	Now(c *gin.Context)
}

type handler struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger, db *gorm.DB) Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) i() {}
