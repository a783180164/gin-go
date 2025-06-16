package ollamatest

import (
	"fmt"
	"gin-go/pkg/code"
	// Core "gin-go/pkg/internal/core"
	"gin-go/pkg/internal/service/ollamatest"
	"gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedTypes = map[string][]string{
	// ".pdf":  {"application/pdf"},
	// ".doc":  {"application/msword"},
	// ".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	".txt": {"text/plain", "text/plain; charset=utf-8", "application/octet-stream"},
	// ".png":  {"image/png"},
	// ".jpg":  {"image/jpeg"},
	// ".jpeg": {"image/jpeg"},
}
var (
	firstErr error
)

// 检查文件类型是否允许
func isAllowed(fileHeader *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	allowedMimes, ok := allowedTypes[ext]
	if !ok {
		return false
	}

	// 打开文件内容，检测实际 mime
	f, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		return false
	}

	detectedMime := http.DetectContentType(buf)
	fmt.Println("detectedMime", detectedMime)
	for _, m := range allowedMimes {
		if m == detectedMime {
			return true
		}
	}

	return false
}

const (
	MaxRequestSize = 50 << 20
)

type UploadRequest struct {
	Collection string `form:"collection" binding:"required"`
}

func (h *handler) Upload(c *gin.Context) {
	req := new(UploadRequest)
	// core := Core.NewContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Info("参数绑定错误")
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: code.Text(code.ParamsError),
		})
		return
	}

	// 1. 限制整个请求体大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxRequestSize)

	// 2. 解析 multipart 表单（内存/临时文件限 32MB）
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    4000,
			Message: fmt.Sprintf("解析表单失败: %v", err),
		})
		return
	}

	files := c.Request.MultipartForm.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, &code.Failure{
			Code:    5000,
			Message: "没有上传任何文件",
		})
		return
	}
	for _, fh := range files {
		if !isAllowed(fh) {
			firstErr = fmt.Errorf("文件类型不被支持: %s", fh.Filename)
			break
		}
	}
	if firstErr != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    60000,
			Data:    "",
			Message: firstErr.Error(),
		})
		return
	}

	model := new(ollamatest.UploadModel)
	model.Collection = req.Collection
	id, err := h.ollamatestService.Upload(model, files)

	if err != nil {
		c.JSON(http.StatusOK, &code.Failure{
			Code:    50000,
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    id,
		Message: code.Text(code.Success),
	})
}
