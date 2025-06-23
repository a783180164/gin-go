package ollamatest

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"

	"gin-go/pkg/internal/embed"
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
	Tags       []string // 可选标签
}

func (s *service) Upload(model *UploadModel, files []*multipart.FileHeader) (int32, error) {
	qd := repository.NewQueryBuilder()
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
					payload := map[string]any{
						"created_at":  time.Now().Format(time.RFC3339),
						"file_name":   fileHeader.Filename,
						"chunk_index": chunkIndex,
						"chunk_size":  n,
						"content":     textChunk,
					}

					// upsert 到 Qdrant
					err = qd.AddDocument(
						s.qd,
						context.Background(),
						model.Collection,
						docID,
						qdrant.NewVectors(vector...),
						qdrant.NewValueMap(payload),
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
