package collection

import (
	"gin-go/pkg/internal/service/collection"
	"gin-go/pkg/internal/service/ollamatest"
	"github.com/gin-gonic/gin"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler interface {
	i()

	Create(c *gin.Context)

	Delete(c *gin.Context)

	Update(c *gin.Context)

	List(c *gin.Context)
}

type handler struct {
	logger            *logrus.Logger
	db                *gorm.DB
	qd                *qdrant.Client
	collectionService collection.Service
	ollamatestService ollamatest.Service
}

func New(logger *logrus.Logger, db *gorm.DB, qd *qdrant.Client) Handler {
	return &handler{
		logger:            logger,
		db:                db,
		qd:                qd,
		collectionService: collection.New(db),
		ollamatestService: ollamatest.New(db, qd),
	}
}

func (h *handler) i() {}
