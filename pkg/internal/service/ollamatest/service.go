package ollamatest

import (
	"mime/multipart"

	"github.com/qdrant/go-client/qdrant"
	"gorm.io/gorm"
)

type Service interface {
	i()
	Upload(model *UploadModel, file []*multipart.FileHeader) (int32, error)
	Prompt(prompt *Prompt) (string, error)
	Create(model *CreateCollection) error
	List(data *ListCollections) (list []*qdrant.ScoredPoint, err error)
	Count(data *ListCollections) (uint64, error)
}

type service struct {
	db *gorm.DB
	qd *qdrant.Client
}

func New(db *gorm.DB, qd *qdrant.Client) Service {
	return &service{
		db: db,
		qd: qd,
	}
}

func (s *service) i() {}
