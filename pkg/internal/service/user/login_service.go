package user

import (
	"gin-go/pkg/internal/mysql"
	repo "gin-go/pkg/internal/repository/user"
)

type UserData struct {
	User     string
	Password string
}

func (s *service) Login(u *UserData) (*repo.User, error) {

	qb := repo.NewQueryBuilder()

	if u.User != "" {
		qb.WhereUser(mysql.EqualPredicate, u.User)
	}

	// if u.Password != "" {
	// 	qb.WherePassword(mysql.EqualPredicate, u.Password)
	// }

	info, err := qb.QueryOne(s.db)
	if err != nil {
		return nil, err
	}
	return info, nil
}
