package ollamatest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"

	"gin-go/pkg/internal/embed"
	"gin-go/pkg/internal/mysql"
	collectionRepository "gin-go/pkg/internal/repository/collection"
	repository "gin-go/pkg/internal/repository/ollamatest"
)

const (
	// 单个文件最大 10MB
	MaxSingleFileSize = 10 << 20
	// 整个请求（所有文件）最大 50MB
	MaxRequestSize = 50 << 20
	// 并发上传的最大 goroutine 数
	MaxConcurrency = 5
	// 每段分块大小 1MB
	ChunkSize = 1 << 20
)

type UploadModel struct {
	Collection string
	UUID       string
	Tags       []string // 可选标签
}

func (s *service) Upload(model *UploadModel, files []*multipart.FileHeader) (int32, error) {
	qd := repository.NewQueryBuilder()
	qb := collectionRepository.NewQueryBuilder()
	qb.WhereUUid(mysql.EqualPredicate, model.UUID)
	info, err := qb.QueryOne(s.db)
	if err != nil {
		return 0, err
	}
	// 并发控制
	sem := make(chan struct{}, MaxConcurrency)
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)

	for _, fh := range files {
		if fh.Size > MaxSingleFileSize {
			firstErr = fmt.Errorf("文件 %s 超过单文件大小限制(%d bytes)", fh.Filename, MaxSingleFileSize)
			break
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()
			defer func() { <-sem }()

			src, err := fileHeader.Open()
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			defer src.Close()

			reader := bufio.NewReader(src)
			chunkIndex := 0
			for {
				buf := make([]byte, ChunkSize)
				n, err := reader.Read(buf)
				if n > 0 {
					textChunk := string(buf[:n])

					// 生成 embedding
					res, err := embed.CallOllamaEmbed(textChunk)
					if err != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = err
						}
						mu.Unlock()
						return
					}
					if len(res.Embeddings) == 0 {
						mu.Lock()
						if firstErr == nil {
							firstErr = fmt.Errorf("empty embeddings")
						}
						mu.Unlock()
						return
					}
					vector := res.Embeddings[0]

					// 生成唯一ID: 文件名+chunkIndex
					u := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fileHeader.Filename+fmt.Sprint(chunkIndex)))
					docID := qdrant.NewID(u.String())

					// 构造 payload
					payloadStruct := &repository.CollectionPoint{
						CreatedAt:  time.Now().Format(time.RFC3339), // 或者你想要的任何时间格式
						Filename:   fileHeader.Filename,
						ChunkIndex: int64(chunkIndex),
						ChunkSize:  int64(n), // 本次读取的字节数
						Content:    textChunk,
					}

					// 1. 序列化 struct
					bts, err := json.Marshal(payloadStruct)
					if err != nil {
						firstErr = fmt.Errorf("marshal payload: %w", err)
					}

					// 2. 反序列化成 map
					var payloadMap map[string]any
					if err := json.Unmarshal(bts, &payloadMap); err != nil {
						firstErr = fmt.Errorf("unmarshal to map: %w", err)
					}

					// upsert 到 Qdrant
					err = qd.AddDocument(
						s.qd,
						context.Background(),
						info.Name+"_"+info.UUID,
						docID,
						qdrant.NewVectors(vector...),
						qdrant.NewValueMap(payloadMap),
					)
					if err != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = fmt.Errorf("写入 Qdrant 失败(%s): %w", docID, err)
						}
						mu.Unlock()
						return
					}

					chunkIndex++
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}
			}

		}(fh)
	}

	wg.Wait()
	if firstErr != nil {
		return 0, firstErr
	}

	return 0, nil
}
