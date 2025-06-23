// pkg/internal/qtstant/repository.go
package qtstant

import (
	"context"
	"fmt"
	qdrantcl "gin-go/pkg/internal/qdrant"
	"github.com/qdrant/go-client/qdrant"
)

type ollamatestQueryBuilder struct {
}

// QdrantRepository 用于以仓库模式操作 Qdrant，
// 封装了向量的写入和检索等功能。
type QdrantRepository struct {
	client *qdrant.Client
}

// NewQdrantRepository 构造 QdrantRepository 实例，
// 使用全局已初始化的客户端。
func NewQdrantRepository() *QdrantRepository {
	return &QdrantRepository{client: qdrantcl.Instance()}
}

func NewQueryBuilder() *ollamatestQueryBuilder {
	return new(ollamatestQueryBuilder)
}

// AddDocument 向指定 collection 插入或更新一个点。
// id 表示点的唯一标识，vector 为向量嵌入值，
// payload 可以存放原始文本等元数据。
func (r *ollamatestQueryBuilder) AddDocument(qd *qdrant.Client, ctx context.Context, collection string, id *qdrant.PointId, vector *qdrant.Vectors, payload map[string]*qdrant.Value) error {
	wait := true
	// 构造 upsert 请求体
	upsert := &qdrant.UpsertPoints{
		CollectionName: collection,
		Points: []*qdrant.PointStruct{{
			Id:      id,
			Vectors: vector,
			Payload: payload,
		}},
		Wait: &wait,
	}
	fmt.Println("upsert", upsert)
	// 执行写入操作
	info, err := qd.Upsert(ctx, upsert)
	fmt.Println("info", info)
	if err != nil {
		return fmt.Errorf("向 Qdrant 写入数据失败：%w", err)
	}
	return nil
}

// AddDocument 向指定 collection 插入或更新一个点。
// id 表示点的唯一标识，vector 为向量嵌入值，
// payload 可以存放原始文本等元数据。
func (r *ollamatestQueryBuilder) QueryDocument(qd *qdrant.Client, ctx context.Context, collection string, vector *qdrant.Query) ([]*qdrant.ScoredPoint, error) {
	// 构造 upsert 请求体
	query := &qdrant.QueryPoints{
		CollectionName: collection,
		Query:          vector,
		// 返回添加 payload
		WithPayload: qdrant.NewWithPayload(true),
		// 返回添加 vectors
		WithVectors: qdrant.NewWithVectors(true),
	}
	fmt.Println("query", query)
	// 执行写入操作
	info, err := qd.Query(ctx, query)
	fmt.Println("info", info)
	if err != nil {
		return nil, fmt.Errorf("向 Qdrant 读取失败%w", err)
	}
	return info, nil
}
