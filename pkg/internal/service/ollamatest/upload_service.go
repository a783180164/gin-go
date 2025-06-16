package ollamatest

import (
	"bufio"
	"context"
	"fmt"
	"gin-go/pkg/internal/embed"
	repository "gin-go/pkg/internal/repository/ollamatest"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"io"
	"mime/multipart"
	"strings"
	"sync"
)

const (
	// 单个文件最大 10MB
	MaxSingleFileSize = 10 << 20
	// 整个请求（所有文件）最大 50MB
	MaxRequestSize = 50 << 20
	// 并发上传的最大 goroutine 数
	MaxConcurrency = 5
)

type UploadModel struct {
	Collection string
}

func (s *service) Upload(model *UploadModel, files []*multipart.FileHeader) (int32, error) {
	id := int32(0)
	qd := repository.NewQueryBuilder()
	// 并发控制信号量
	sem := make(chan struct{}, MaxConcurrency)
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)

	for _, fh := range files {
		// 3. 文件大小检查
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

			// 读取文件内容
			reader := bufio.NewReader(src)
			var sb strings.Builder
			for {
				line, err := reader.ReadString('\n')
				sb.WriteString(line)
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
			res, err := embed.CallOllamaEmbed(sb.String())
			if err != nil {
				mu.Lock()
				firstErr = err
				mu.Unlock()
				return
			}

			// docID := qdrant.NewID(uuid.NewString())
			u := uuid.NewSHA1(uuid.NameSpaceURL, []byte(sb.String()))
			docID := qdrant.NewID(u.String())
			// 取出第一个（也是唯一一个）向量
			if len(res.Embeddings) == 0 {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("empty embeddings")
				}
				mu.Unlock()
				return // 或者报错：“empty embeddings”
			}
			vector := res.Embeddings[0]
			err = qd.AddDocument(s.qd, context.Background(), model.Collection, docID, qdrant.NewVectors(vector...) /* vector 留空由 QA 层填充 */, qdrant.NewValueMap(map[string]any{"text": sb.String()}))
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("写入 Qdrant 失败(%s): %w", docID, err)
				}
				mu.Unlock()
				return
			}

		}(fh)
	}

	wg.Wait()

	if firstErr != nil {

		return 0, firstErr
	}

	return id, nil
}
