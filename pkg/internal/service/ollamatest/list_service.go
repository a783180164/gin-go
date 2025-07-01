package ollamatest

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/ollamatest"

	collectionRepository "gin-go/pkg/internal/repository/collection"

	"github.com/qdrant/go-client/qdrant"
)

type ListCollections struct {
	Name        string
	Description string
	Prompt      string
	Page        int
	PageSize    int
	UUID        string
}

func (s *service) List(data *ListCollections) ([]*qdrant.ScoredPoint, error) {
	page := data.Page
	if page == 0 {
		page = 1
	}

	pageSize := data.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	qb := collectionRepository.NewQueryBuilder()
	qb.WhereUUid(mysql.EqualPredicate, data.UUID)
	info, err := qb.QueryOne(s.db)
	if err != nil {
		return nil, err
	}

	qd := ollamatest.NewQueryBuilder()
	qd.Limit(uint64(pageSize))
	qd.Offset(uint64(offset))
	qd.WhereCollection(info.Name + "_" + info.UUID)
	list, err := qd.QueryAll(s.qd)

	if err != nil {
		return nil, err
	}
	return list, nil
}
