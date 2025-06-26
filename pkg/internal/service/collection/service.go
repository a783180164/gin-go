package collection

import (
	"gin-go/pkg/internal/repository/collection"
	"gorm.io/gorm"
)

type Service interface {
	i()
	Create(data *CreateCollection) (string, error)
	Delete(uuid []string) error
	Update(data *UpdateCollection) error

	List(data *ListCollection) (list []*collection.Collection, err error)

	Count(data *ListCollection) (total int64, err error)
}

type service struct {
	db *gorm.DB
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

func (s *service) i() {}
