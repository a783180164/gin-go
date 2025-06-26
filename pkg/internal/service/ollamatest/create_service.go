package ollamatest

import (
	"context"

	repository "gin-go/pkg/internal/repository/ollamatest"
)

type CreateCollection struct {
	Name string
	UUID string
	Size uint64
}

func (s *service) Create(model *CreateCollection) error {
	// 1. 初次 embed + qdrant 查询
	qd := repository.NewQueryBuilder()
	err := qd.CreateCollection(
		s.qd,
		context.Background(),
		model.Name,
		model.UUID,
		model.Size,
	)
	if err != nil {
		return err
	}
	return nil
}
