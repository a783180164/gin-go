package collection

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/collection"
)

type UpdateCollection struct {
	UUID        string
	Description string
	Prompt      string
}

func (s *service) Update(data *UpdateCollection) error {

	qb := collection.NewQueryBuilder()
	model := map[string]interface{}{
		"description": data.Description,
		"prompt":      data.Prompt,
	}

	qb.WhereUUid(mysql.EqualPredicate, data.UUID)

	err := qb.Updates(s.db, model)

	if err != nil {
		return err
	}
	return nil
}
