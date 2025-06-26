package collection

import (
	"gin-go/pkg/errors"
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/collection"
	"github.com/google/uuid"
)

type CreateCollection struct {
	UUID        string
	Name        string
	Description string
	Prompt      string
	IsDelete    int32
}

func (s *service) Create(data *CreateCollection) (string, error) {
	model := collection.NewModel()

	qb := collection.NewQueryBuilder()

	qb.WhereName(mysql.EqualPredicate, data.Name)

	qb.WhereIsDelete(mysql.EqualPredicate, 0)

	info, err := qb.QueryOne(s.db)
	if err != nil {
		return "", err
	}

	if info != nil {
		return "", errors.New("知识库已经存在")
	}
	model.Name = data.Name
	model.Prompt = data.Prompt
	model.IsDeleted = false
	model.UUID = uuid.NewString()
	model.Description = data.Description
	uuid, err := model.Create(s.db)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
