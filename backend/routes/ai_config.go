package routes

import (
	"os"
	"strings"
)

// getQwenAPIKey 从环境变量获取 QWEN_API_KEY
func getQwenAPIKey() string {
	return strings.TrimSpace(os.Getenv("QWEN_API_KEY"))
}

// getQwenBaseURL 从环境变量获取基础域名，未设置时使用默认 dashscope
func getQwenBaseURL() string {
	base := strings.TrimSpace(os.Getenv("QWEN_API_BASE_URL"))
	if base == "" {
		base = "https://dashscope.aliyuncs.com"
	}
	return strings.TrimRight(base, "/")
}

// qwenChatURL 返回通义聊天补全接口地址
func qwenChatURL() string {
	return getQwenBaseURL() + "/compatible-mode/v1/chat/completions"
}
