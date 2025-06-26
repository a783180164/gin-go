package ollamatest

import (
	"github.com/qdrant/go-client/qdrant"
	"gorm.io/gorm"
	"mime/multipart"
)

type Service interface {
	i()
	Upload(model *UploadModel, file []*multipart.FileHeader) (int32, error)
	Prompt(prompt *Prompt) (string, error)
	Create(model *CreateCollection) error
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
