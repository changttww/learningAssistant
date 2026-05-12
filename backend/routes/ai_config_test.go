package routes

import "testing"

func TestAIConfigDefaultsToDashScope(t *testing.T) {
	t.Setenv("QWEN_API_BASE_URL", "")
	t.Setenv("QWEN_CHAT_COMPLETIONS_PATH", "")
	t.Setenv("QWEN_CHAT_MODEL", "")
	t.Setenv("QWEN_FAST_MODEL", "")

	if got := qwenChatURL(); got != "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions" {
		t.Fatalf("unexpected default chat URL: %s", got)
	}
	if got := qwenChatModel(); got != "qwen-plus" {
		t.Fatalf("unexpected default chat model: %s", got)
	}
	if got := qwenFastModel(); got != "qwen-turbo" {
		t.Fatalf("unexpected default fast model: %s", got)
	}
}

func TestAIConfigSupportsOpenAICompatibleProvider(t *testing.T) {
	t.Setenv("QWEN_API_BASE_URL", "https://api.deepseek.com/")
	t.Setenv("QWEN_CHAT_COMPLETIONS_PATH", "/chat/completions")
	t.Setenv("QWEN_CHAT_MODEL", "deepseek-v4-pro")
	t.Setenv("QWEN_FAST_MODEL", "deepseek-v4-pro")

	if got := qwenChatURL(); got != "https://api.deepseek.com/chat/completions" {
		t.Fatalf("unexpected custom chat URL: %s", got)
	}
	if got := qwenChatModel(); got != "deepseek-v4-pro" {
		t.Fatalf("unexpected custom chat model: %s", got)
	}
	if got := qwenFastModel(); got != "deepseek-v4-pro" {
		t.Fatalf("unexpected custom fast model: %s", got)
	}
}
