// ollama_embed_refactor.go
package ollamatest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-go/configs"
	"gin-go/pkg/code"
	Core "gin-go/pkg/internal/core"
	"gin-go/pkg/logger"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type EmbedRequest struct {
	Text string `json:"text" binding:"required"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

// handler 结构体略

// EmbedWithOllama 处理 HTTP 请求，并调用封装好的 POST 方法
func (h *handler) EmbedWithOllama(c *gin.Context) {
	req := new(EmbedRequest)
	core := Core.NewContext(c)
	if err := core.ShouldBindJSON(req); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Info("参数绑定错误")
		c.JSON(http.StatusOK, &code.Failure{
			Code:    code.ParamsError,
			Data:    nil,
			Message: code.Text(code.ParamsError),
		})
		return
	}

	// 调用封装方法
	respData, err := h.callOllamaEmbed(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &code.Failure{
			Code:    code.JsonBodyError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &code.Failure{
		Code:    code.Success,
		Data:    respData,
		Message: code.Text(code.Success),
	})
}

// callOllamaEmbed 封装向 Ollama 发起 POST 请求并解析返回
func (h *handler) callOllamaEmbed(text string) (*interface{}, error) {
	cfg := configs.Get().OLLAMA

	// 构造请求体
	bodyObj := map[string]interface{}{"model": cfg.Model, "input": text}
	bodyBytes, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	url := fmt.Sprintf("%s:%d/api/embed", cfg.Host, cfg.Port)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Ollama API 失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("非 200 响应: %d, body: %s", resp.StatusCode, string(data))
	}

	// 解析 JSON
	var embedResp interface{}
	if err := json.Unmarshal(data, &embedResp); err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	return &embedResp, nil
}
