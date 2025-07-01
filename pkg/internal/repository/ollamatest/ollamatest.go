// pkg/internal/qtstant/repository.go
package ollamatest

import (
	"context"
	"fmt"
	qdrantcl "gin-go/pkg/internal/qdrant"

	"github.com/qdrant/go-client/qdrant"
)

type ollamatestQueryBuilder struct {
	filter      qdrant.Filter
	limit       uint64
	offset      uint64
	WithPayload bool
	query       []float32
	collection  string
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
	return &ollamatestQueryBuilder{
		filter: qdrant.Filter{
			Must:    []*qdrant.Condition{},
			Should:  []*qdrant.Condition{},
			MustNot: []*qdrant.Condition{},
		},
		limit:       0, // 例如默认 limit=100
		offset:      0, // 默认 offset=0
		WithPayload: true,
		query:       nil,
	}
}

func (r *ollamatestQueryBuilder) buildQuery() *qdrant.QueryPoints {
	// 构建通用的 QueryPoints
	qp := &qdrant.QueryPoints{
		CollectionName: r.collection,
		Filter:         &r.filter,
		WithPayload:    qdrant.NewWithPayload(r.WithPayload),
	}
	// 分页参数：只有非负值才生效
	if r.limit > 0 {
		qp.Limit = &r.limit
	}
	if r.offset > 0 {
		qp.Offset = &r.offset
	}
	// 如果有向量，则执行向量查询
	if len(r.query) > 0 {
		// 展开切片；请确保 r.Query 长度和 collection 的向量维度一致
		qp.Query = qdrant.NewQuery(r.query...)
	}
	fmt.Println(qp, r)
	return qp
}

func (r *ollamatestQueryBuilder) buildCount() *qdrant.CountPoints {
	// 构建通用的 QueryPoints
	qp := &qdrant.CountPoints{
		CollectionName: r.collection,
		Filter:         &r.filter,
	}
	return qp
}

// QueryAll 支持：
// - Query 为空时，只按 filter + 分页滚动获取
// - Query 非空时，执行向量检索 + filter + 分页
func (r *ollamatestQueryBuilder) QueryAll(
	qd *qdrant.Client,
) ([]*qdrant.ScoredPoint, error) {
	qp := r.buildQuery()
	// 发起请求
	resp, err := qd.Query(context.Background(), qp)
	if err != nil {
		return nil, fmt.Errorf("查询 Qdrant 失败：%w", err)
	}
	return resp, nil
}

func (r *ollamatestQueryBuilder) QueryOne(qd *qdrant.Client) (*qdrant.ScoredPoint, error) {
	resp, err := r.QueryAll(qd)
	if len(resp) > 0 {
		return resp[0], err
	}
	return nil, err
}

func (r *ollamatestQueryBuilder) WhereCollection(name string) {
	r.collection = name
}

func (r *ollamatestQueryBuilder) Limit(limit uint64) {
	r.limit = limit
}

func (r *ollamatestQueryBuilder) Offset(offset uint64) {
	r.offset = offset
}

func (r *ollamatestQueryBuilder) WhereQuery(query []float32) {
	r.query = query
}

func (r *ollamatestQueryBuilder) Count(qd *qdrant.Client) (uint64, error) {
	qp := r.buildCount()
	// qd.Get
	count, err := qd.Count(context.Background(), qp)
	if err != nil {
		return 0, err
	}
	return count, nil
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
	_, err := qd.Upsert(ctx, upsert)
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

func (r *ollamatestQueryBuilder) CreateCollection(qd *qdrant.Client, ctx context.Context, name string, uuid string, size uint64) error {
	err := qd.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: fmt.Sprintf("%v_%v", name, uuid),
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     size,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	fmt.Println(fmt.Printf("%v_%v", name, uuid))
	if err != nil {
		return err
	}
	return nil
}
