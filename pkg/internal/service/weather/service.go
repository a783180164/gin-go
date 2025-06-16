package weather

import (
	"gorm.io/gorm"
)

type Service interface {
	i()
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
