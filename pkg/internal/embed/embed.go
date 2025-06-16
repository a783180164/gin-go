package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-go/configs"
	"io/ioutil"
	"net/http"
)

type Embed struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

// callOllamaEmbed 封装向 Ollama 发起 POST 请求并解析返回
func CallOllamaEmbed(text string) (*Embed, error) {
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
	var embedResp Embed
	if err := json.Unmarshal(data, &embedResp); err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	return &embedResp, nil
}
