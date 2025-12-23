package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"

	"learningAssistant-backend/models"
)

// Vector 别名，方便使用
type Vector = models.Vector

// QwenEmbeddingService Qwen向量化服务
type QwenEmbeddingService struct {
	apiKey     string
	apiURL     string
	model      string
	httpClient *http.Client
}

// QwenEmbeddingRequest Qwen embedding请求
type QwenEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// QwenEmbeddingResponse Qwen embedding响应
type QwenEmbeddingResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
	Usage struct {
		InputTokens int `json:"input_tokens"`
	} `json:"usage"`
}

// NewQwenEmbeddingService 创建Qwen embedding服务
func NewQwenEmbeddingService(apiKey string) *QwenEmbeddingService {
	if apiKey == "" {
		apiKey = os.Getenv("QWEN_API_KEY")
	}

	return &QwenEmbeddingService{
		apiKey:     apiKey,
		apiURL:     "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding",
		model:      "text-embedding-v3",
		httpClient: &http.Client{},
	}
}

// GenerateEmbedding 生成文本向量
func (s *QwenEmbeddingService) GenerateEmbedding(text string) (Vector, error) {
	if s.apiKey == "" {
		// 返回模拟向量以支持本地开发
		return s.mockEmbedding(text), nil
	}

	vectors, err := s.GenerateEmbeddings([]string{text})
	if err != nil {
		return nil, err
	}

	if len(vectors) == 0 {
		return nil, fmt.Errorf("未获得向量数据")
	}

	return vectors[0], nil
}

// GenerateEmbeddings 批量生成向量
func (s *QwenEmbeddingService) GenerateEmbeddings(texts []string) ([]Vector, error) {
	if s.apiKey == "" {
		// 返回模拟向量以支持本地开发
		vectors := make([]Vector, len(texts))
		for i, text := range texts {
			vectors[i] = s.mockEmbedding(text)
		}
		return vectors, nil
	}

	if len(texts) == 0 {
		return []Vector{}, nil
	}

	// 限制单次请求大小
	if len(texts) > 10 {
		var allVectors []Vector
		for i := 0; i < len(texts); i += 10 {
			end := i + 10
			if end > len(texts) {
				end = len(texts)
			}
			batch, err := s.callQwenAPI(texts[i:end])
			if err != nil {
				return nil, err
			}
			allVectors = append(allVectors, batch...)
		}
		return allVectors, nil
	}

	return s.callQwenAPI(texts)
}

// callQwenAPI 调用Qwen API
func (s *QwenEmbeddingService) callQwenAPI(texts []string) ([]Vector, error) {
	req := QwenEmbeddingRequest{
		Model: s.model,
		Input: texts,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var respData QwenEmbeddingResponse
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if respData.Error != nil {
		return nil, fmt.Errorf("API错误: %s (%s)", respData.Error.Message, respData.Error.Code)
	}

	vectors := make([]Vector, len(respData.Data))
	for _, item := range respData.Data {
		vectors[item.Index] = Vector(item.Embedding)
	}

	return vectors, nil
}

// CosineSimilarity 计算向量余弦相似度
func (s *QwenEmbeddingService) CosineSimilarity(vec1, vec2 Vector) float32 {
	if len(vec1) == 0 || len(vec2) == 0 {
		return 0
	}

	if len(vec1) != len(vec2) {
		return 0
	}

	var dotProduct, norm1, norm2 float64

	for i := range vec1 {
		dotProduct += float64(vec1[i] * vec2[i])
		norm1 += math.Pow(float64(vec1[i]), 2)
		norm2 += math.Pow(float64(vec2[i]), 2)
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	similarity := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
	return float32(similarity)
}

// mockEmbedding 生成模拟向量用于本地开发
func (s *QwenEmbeddingService) mockEmbedding(text string) Vector {
	// 简单的模拟：基于字符哈希生成固定维度的向量
	hash := 0
	for i, char := range text {
		hash += (int(char) * (i + 1))
	}

	dim := 1536 // Qwen embedding维度
	vec := make(Vector, dim)

	seed := float32(hash % 1000000)
	for i := 0; i < dim; i++ {
		// 使用伪随机生成
		val := math.Sin(float64(seed+float32(i))*0.001) * 0.5
		vec[i] = float32(val)
	}

	return vec
}

// LocalEmbeddingService 本地embedding服务（不需要API key）
type LocalEmbeddingService struct{}

// NewLocalEmbeddingService 创建本地embedding服务
func NewLocalEmbeddingService() *LocalEmbeddingService {
	return &LocalEmbeddingService{}
}

func (s *LocalEmbeddingService) GenerateEmbedding(text string) (Vector, error) {
	return s.simpleEmbedding(text), nil
}

func (s *LocalEmbeddingService) GenerateEmbeddings(texts []string) ([]Vector, error) {
	vectors := make([]Vector, len(texts))
	for i, text := range texts {
		vectors[i] = s.simpleEmbedding(text)
	}
	return vectors, nil
}

func (s *LocalEmbeddingService) CosineSimilarity(vec1, vec2 Vector) float32 {
	if len(vec1) == 0 || len(vec2) == 0 {
		return 0
	}

	if len(vec1) != len(vec2) {
		return 0
	}

	var dotProduct, norm1, norm2 float64
	for i := range vec1 {
		dotProduct += float64(vec1[i] * vec2[i])
		norm1 += math.Pow(float64(vec1[i]), 2)
		norm2 += math.Pow(float64(vec2[i]), 2)
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2)))
}

func (s *LocalEmbeddingService) simpleEmbedding(text string) Vector {
	// 简单词频特征向量
	dim := 256
	vec := make(Vector, dim)

	for i, char := range text {
		idx := (int(char) * (i + 1)) % dim
		vec[idx] += 0.1
	}

	// 归一化
	var norm float64
	for _, v := range vec {
		norm += math.Pow(float64(v), 2)
	}
	norm = math.Sqrt(norm)

	if norm > 0 {
		for i := range vec {
			vec[i] = float32(float64(vec[i]) / norm)
		}
	}

	return vec
}
