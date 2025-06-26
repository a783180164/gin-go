package collection

import (
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/repository/collection"
)

func (s *service) Delete(uuids []string) error {

	data := map[string]interface{}{
		"is_deleted": 1,
	}
	qb := collection.NewQueryBuilder()

	qb.WhereUUids(mysql.InPredicate, uuids)
	err := qb.Updates(s.db, data)

	if err != nil {
		return err
	}
	return nil
}
