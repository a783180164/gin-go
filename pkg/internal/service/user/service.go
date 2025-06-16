package user

import (
	repo "gin-go/pkg/internal/repository/user"
	"gorm.io/gorm"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	CreateUser(user *CreateUserData) (int32, error)

	Login(u *UserData) (*repo.User, error)
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
