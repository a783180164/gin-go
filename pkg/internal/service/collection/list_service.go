package collection

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/collection"
)

type ListCollection struct {
	Name        string
	Description string
	Prompt      string
	Page        int
	PageSize    int
}

func (s *service) List(data *ListCollection) (list []*collection.Collection, err error) {

	page := data.Page
	if page == 0 {
		page = 1
	}

	pageSize := data.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	qb := collection.NewQueryBuilder()

	if data.Name != "" {
		qb.WhereName(mysql.LikePredicate, "%"+data.Name+"%")
	}

	if data.Description != "" {
		qb.WhereDesc(mysql.LikePredicate, data.Description)
	}

	if data.Prompt != "" {
		qb.WherePrompt(mysql.LikePredicate, data.Prompt)
	}

	lists, err := qb.Limit(pageSize).Offset(offset).OrderById(false).QueryAll(s.db)

	if err != nil {
		return nil, err
	}
	return lists, nil
}
