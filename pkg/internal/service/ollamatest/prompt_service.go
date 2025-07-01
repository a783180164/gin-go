package ollamatest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"

	"gin-go/pkg/internal/embed"
	"gin-go/pkg/internal/mysql"
	collectionRepository "gin-go/pkg/internal/repository/collection"
	repository "gin-go/pkg/internal/repository/ollamatest"
)

type Prompt struct {
	Text string
	UUID string
}

// normalize 对字符串去标点、去多余空格、统一小写
func normalize(s string) string {
	re := regexp.MustCompile(`[，。？！；：,.?!;:\n\r]`)
	s = re.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.ToLower(s)
}

// hashMD5 返回 s 的 MD5 十六进制摘要
func hashMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// cosineSim 计算两个向量的余弦相似度
func cosineSim(a, b []float32) float64 {
	var dot, na, nb float64
	for i := range a {
		x, y := float64(a[i]), float64(b[i])
		dot += x * y
		na += x * x
		nb += y * y
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

// buildPrompt 将去重后的文本列表拼成最终的 prompt
func buildPrompt(texts []string) string {
	var lines []string
	for i, t := range texts {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, t))
	}
	return fmt.Sprintf(
		"请参考以下内容，回答接下来的问题：\n\n%s\n\n（请基于以上内容，不要凭空编造。）",
		strings.Join(lines, "\n"),
	)
}

func (s *service) Prompt(prompt *Prompt) (string, error) {

	qb := collectionRepository.NewQueryBuilder()
	qb.WhereUUid(mysql.EqualPredicate, prompt.UUID)
	Info, err := qb.QueryOne(s.db)

	if err != nil {
		return "", err
	}
	// 1. 初次 embed + qdrant 查询
	embRes, err := embed.CallOllamaEmbed(prompt.Text)
	if err != nil {
		return "", err
	}
	qd := repository.NewQueryBuilder()

	qd.WhereCollection(Info.Name + "_" + Info.UUID)
	qd.WhereQuery(embRes.Embeddings[0])
	raw, err := qd.QueryAll(
		s.qd,
	)
	if err != nil {
		return "", err
	}

	// 2. MD5 去重，收集文本和原始向量
	type item struct {
		text string
		vec  []float32
	}
	seen := make(map[string]bool)
	var candidates []item
	for _, pt := range raw {
		payloadVal, ok := pt.Payload["content"]
		if !ok {
			continue
		}
		text := payloadVal.GetStringValue()
		norm := normalize(text)
		key := hashMD5(norm)
		if seen[key] {
			continue
		}
		seen[key] = true
		// 获取 qdrant 返回的向量（float64 slice），然后转换为 float32
		if pt.Vectors == nil {

			continue
		}
		vecOut := pt.Vectors.GetVector()
		float64Vec := vecOut.Data
		vec32 := make([]float32, len(float64Vec))
		for i, v := range float64Vec {
			vec32[i] = float32(v)
		}
		candidates = append(candidates, item{text: text, vec: vec32})
	}

	// 3. 二次向量去重（阈值 0.95）
	const threshold = 0.95
	var finalTexts []string
	var finalVecs [][]float32
	for _, it := range candidates {
		dup := false
		for _, v := range finalVecs {
			if cosineSim(it.vec, v) > threshold {
				dup = true
				break
			}
		}
		if !dup {
			finalTexts = append(finalTexts, it.text)
			finalVecs = append(finalVecs, it.vec)
		}
	}

	// 4. 拼成最终 Prompt
	finalPrompt := buildPrompt(finalTexts)
	return finalPrompt, nil
}
