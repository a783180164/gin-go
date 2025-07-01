package ollamatest

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/ollamatest"

	collectionRepository "gin-go/pkg/internal/repository/collection"
)

func (s *service) Count(data *ListCollections) (uint64, error) {

	qb := collectionRepository.NewQueryBuilder()
	qb.WhereUUid(mysql.EqualPredicate, data.UUID)
	info, err := qb.QueryOne(s.db)
	if err != nil {
		return 0, err
	}

	qd := ollamatest.NewQueryBuilder()
	qd.WhereCollection(info.Name + "_" + info.UUID)
	count, err := qd.Count(s.qd)

	if err != nil {
		return 0, err
	}
	return count, nil
}
