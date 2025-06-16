package ollamatest

import (
	"context"
	"fmt"
	"gin-go/pkg/internal/embed"
	repository "gin-go/pkg/internal/repository/ollamatest"
	"github.com/qdrant/go-client/qdrant"
)

type Prompt struct {
	Text string
}

func (s *service) Prompt(prompt *Prompt) ([]*qdrant.ScoredPoint, error) {
	fmt.Println("pro", prompt.Text)
	res, err := embed.CallOllamaEmbed(prompt.Text)
	if err != nil {
		return nil, err
	}
	qd := repository.NewQueryBuilder()
	vector := res.Embeddings[0]

	data, err := qd.QueryDocument(s.qd, context.Background(), "txt_collection", qdrant.NewQuery(vector...))
	if err != nil {
		return nil, err
	}

	return data, nil
}
