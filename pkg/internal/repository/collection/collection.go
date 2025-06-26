package collection

import (
	"fmt"
	"gin-go/pkg/errors"
	"gin-go/pkg/internal/mysql"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	db *gorm.DB
}

type CollectionQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func NewModel() *Collection {
	return new(Collection)
}

// NewUserRepository 构造函数
func NewUserRepository() *CollectionRepository {
	return &CollectionRepository{db: mysql.Instance()}
}

func NewQueryBuilder() *CollectionQueryBuilder {
	return new(CollectionQueryBuilder)
}

func (qb *CollectionQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
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

func (qb *CollectionQueryBuilder) QueryOne(db *gorm.DB) (*Collection, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *CollectionQueryBuilder) QueryAll(db *gorm.DB) ([]*Collection, error) {
	var ret []*Collection
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *CollectionQueryBuilder) Limit(limit int) *CollectionQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *CollectionQueryBuilder) Offset(offset int) *CollectionQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *CollectionQueryBuilder) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	db = db.Model(&Collection{})

	for _, where := range qb.where {
		db.Where(where.prefix, where.value)
	}

	if err = db.Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *CollectionQueryBuilder) Count(db *gorm.DB) (int64, error) {
	var c int64
	res := qb.buildQuery(db).Model(&Collection{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (t *Collection) Create(db *gorm.DB) (uuid string, err error) {
	if err := db.Create(t).Error; err != nil {
		return "", errors.Wrap(err, "create err")
	}
	return t.UUID, nil
}

func (qb *CollectionQueryBuilder) WherePassword(p mysql.Predicate, value string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", p),
		value,
	})
	return qb
}
func (qb *CollectionQueryBuilder) OrderById(asc bool) *CollectionQueryBuilder {

	order := "DESC"
	if asc {
		order = "ASC"
	}
	qb.order = append(qb.order, "id  "+order)
	return qb
}

func (qb *CollectionQueryBuilder) WhereName(p mysql.Predicate, value string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "name", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) WhereDesc(p mysql.Predicate, value string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "description", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) WherePrompt(p mysql.Predicate, value string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "prompt", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) WhereUUid(p mysql.Predicate, value string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "uuid", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) WhereUUids(p mysql.Predicate, value []string) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v (?)", "uuid", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) WhereIsDelete(p mysql.Predicate, value int32) *CollectionQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", p),
		value,
	})
	return qb
}

func (qb *CollectionQueryBuilder) Delete(db *gorm.DB) (err error) {
	for _, where := range qb.where {
		db = db.Where(where.prefix, where.value)
	}

	if err = db.Delete(&Collection{}).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}
