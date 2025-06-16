package ollamatest

import (
	"gin-go/pkg/internal/service/ollamatest"
	"github.com/gin-gonic/gin"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler interface {
	i()

	Upload(c *gin.Context)

	EmbedWithOllama(c *gin.Context)

	Prompt(c *gin.Context)
}

type handler struct {
	logger            *logrus.Logger
	db                *gorm.DB
	qd                *qdrant.Client
	ollamatestService ollamatest.Service
}

func New(logger *logrus.Logger, db *gorm.DB, qd *qdrant.Client) Handler {
	return &handler{
		logger:            logger,
		db:                db,
		qd:                qd,
		ollamatestService: ollamatest.New(db, qd),
	}
}

func (h *handler) i() {}
