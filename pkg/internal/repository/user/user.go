// pkg/internal/repository/user.go
package user

import (
	"errors"
	"fmt"
	"gin-go/pkg/internal/model"
	"gin-go/pkg/internal/mysql"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type userQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func NewModel() *User {
	return new(User)
}

// NewUserRepository 构造函数
func NewUserRepository() *UserRepository {
	return &UserRepository{db: mysql.Instance()}
}

func NewQueryBuilder() *userQueryBuilder {
	return new(userQueryBuilder)
}

// FindByUsername 根据用户名查用户
func (r *userQueryBuilder) FindByUsername(db *gorm.DB, user string) (*model.User, error) {
	var u model.User
	if err := db.Where("user = ?", user).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreatetUser(u *User) (int32, error) {
	var existing User
	if err := r.db.Where("user = ?", u.User).First(&existing).Error; err == nil {
		// 用户已存在
		return 0, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 查询过程中出现其他错误
		return 0, err
	}

	result := r.db.Create(u)
	if result.Error != nil {
		return 0, result.Error
	}
	return u.ID, nil
}

func (qb *userQueryBuilder) WherePassword(p mysql.Predicate, value string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereUser(p mysql.Predicate, value string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "user", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
	ret := db
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *userQueryBuilder) QueryOne(db *gorm.DB) (*User, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *userQueryBuilder) QueryAll(db *gorm.DB) ([]*User, error) {
	var ret []*User
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}
