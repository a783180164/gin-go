// pkg/internal/service/user/user.go
package user

import (
	"fmt"
	repository "gin-go/pkg/internal/repository/user"
	"github.com/google/uuid"
)

type CreateUserData struct {
	User     string // 用户名
	Password string // 密码
}

func (s *service) CreateUser(user *CreateUserData) (int32, error) {
	repo := repository.NewUserRepository()
	var u repository.User
	u.User = user.User
	u.Password = user.Password
	u.UUID = uuid.NewString()
	fmt.Println(u, user)
	id, err := repo.CreatetUser(&u)
	if err != nil {
		return 0, err
	}
	return id, nil
}
