package collection

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/collection"
)

func (s *service) Count(data *ListCollection) (total int64, err error) {

	qb := collection.NewQueryBuilder()

	qb = qb.WhereIsDelete(mysql.EqualPredicate, 0)
	if data.Name != "" {
		qb.WhereName(mysql.LikePredicate, "%"+data.Name+"%")
	}

	if data.Description != "" {
		qb.WhereDesc(mysql.LikePredicate, data.Description)
	}

	if data.Prompt != "" {
		qb.WhereDesc(mysql.LikePredicate, data.Prompt)
	}
	qb.Limit(-1)
	total, err = qb.Count(s.db)
	if err != nil {
		return 0, err
	}

	return
}
