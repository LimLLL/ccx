package converters

import (
	"context"
	"strings"
	"testing"
)

func TestConvertOpenAIChatToResponses_Stream(t *testing.T) {
	ctx := context.Background()

	// 模拟 OpenAI Chat Completions SSE 流
	sseLines := []string{
		`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":" world!"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`,
		`data: [DONE]`,
	}

	originalReq := []byte(`{"model":"gpt-4o","input":"Hi"}`)

	var state any
	var allEvents []string

	for _, line := range sseLines {
		events := ConvertOpenAIChatToResponses(ctx, "gpt-4o", originalReq, nil, []byte(line), &state)
		allEvents = append(allEvents, events...)
	}

	// 验证事件序列
	if len(allEvents) == 0 {
		t.Fatal("should produce events")
	}

	// 检查是否有 response.created 事件
	hasCreated := false
	hasInProgress := false
	hasCompleted := false
	hasTextDelta := false

	for _, ev := range allEvents {
		if strings.Contains(ev, "response.created") {
			hasCreated = true
		}
		if strings.Contains(ev, "response.in_progress") {
			hasInProgress = true
		}
		if strings.Contains(ev, "response.completed") {
			hasCompleted = true
		}
		if strings.Contains(ev, "response.output_text.delta") {
			hasTextDelta = true
		}
	}

	if !hasCreated {
		t.Error("should have response.created event")
	}
	if !hasInProgress {
		t.Error("should have response.in_progress event")
	}
	if !hasCompleted {
		t.Error("should have response.completed event")
	}
	if !hasTextDelta {
		t.Error("should have response.output_text.delta event")
	}
}

func TestConvertOpenAIChatToResponses_StreamReasoningContent(t *testing.T) {
	ctx := context.Background()
	sseLines := []string{
		`data: {"id":"chatcmpl-ds","object":"chat.completion.chunk","created":1234567890,"model":"deepseek-v4-pro","choices":[{"index":0,"delta":{"reasoning_content":"chat reasoning"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-ds","object":"chat.completion.chunk","created":1234567890,"model":"deepseek-v4-pro","choices":[{"index":0,"delta":{"content":"chat text"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-ds","object":"chat.completion.chunk","created":1234567890,"model":"deepseek-v4-pro","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
		`data: [DONE]`,
	}

	var state any
	var allEvents []string
	for _, line := range sseLines {
		events := ConvertOpenAIChatToResponses(ctx, "deepseek-v4-pro", []byte(`{"model":"deepseek-v4-pro","input":"hello"}`), nil, []byte(line), &state)
		allEvents = append(allEvents, events...)
	}

	joined := strings.Join(allEvents, "\n")
	if !strings.Contains(joined, `"type":"reasoning"`) {
		t.Fatalf("expected reasoning item, got %v", allEvents)
	}
	if !strings.Contains(joined, `"text":"chat reasoning"`) {
		t.Fatalf("expected reasoning summary text, got %v", allEvents)
	}
	if !strings.Contains(joined, `"delta":"chat text"`) {
		t.Fatalf("expected text delta after reasoning, got %v", allEvents)
	}
}

func TestConvertOpenAIChatToResponses_ToolCall(t *testing.T) {
	ctx := context.Background()

	// 模拟带 tool_call 的 SSE 流
	sseLines := []string{
		`data: {"id":"chatcmpl-456","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":null,"tool_calls":[{"index":0,"id":"call_abc","type":"function","function":{"name":"get_weather","arguments":""}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-456","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"loc"}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-456","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"ation\": \"NYC\"}"}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-456","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`,
		`data: [DONE]`,
	}

	originalReq := []byte(`{"model":"gpt-4o","input":"What's the weather?","tools":[{"name":"get_weather"}]}`)

	var state any
	var allEvents []string

	for _, line := range sseLines {
		events := ConvertOpenAIChatToResponses(ctx, "gpt-4o", originalReq, nil, []byte(line), &state)
		allEvents = append(allEvents, events...)
	}

	// 验证是否有 function_call 相关事件
	hasFuncAdded := false
	hasFuncDelta := false
	hasFuncDone := false

	for _, ev := range allEvents {
		if strings.Contains(ev, "response.output_item.added") && strings.Contains(ev, "function_call") {
			hasFuncAdded = true
		}
		if strings.Contains(ev, "response.function_call_arguments.delta") {
			hasFuncDelta = true
		}
		if strings.Contains(ev, "response.function_call_arguments.done") {
			hasFuncDone = true
		}
	}

	if !hasFuncAdded {
		t.Error("should have function_call output_item.added event")
	}
	if !hasFuncDelta {
		t.Error("should have function_call_arguments.delta event")
	}
	if !hasFuncDone {
		t.Error("should have function_call_arguments.done event")
	}
}

func TestConvertOpenAIChatToResponses_CustomToolCall(t *testing.T) {
	ctx := context.Background()
	sseLines := []string{
		`data: {"id":"chatcmpl-custom","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","tool_calls":[{"index":0,"id":"call_patch","type":"function","function":{"name":"apply_patch_add_file","arguments":""}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-custom","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"path\":\"docs/test.md\",\"content\":\"# Test\\n\"}"}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-custom","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`,
		`data: [DONE]`,
	}
	originalReq := []byte(`{"model":"gpt-4o","input":"edit","tools":[{"type":"custom","name":"apply_patch"}],"transformer_metadata":{"codex_tool_compat_enabled":true}}`)

	var state any
	var allEvents []string
	for _, line := range sseLines {
		allEvents = append(allEvents, ConvertOpenAIChatToResponses(ctx, "gpt-4o", originalReq, nil, []byte(line), &state)...)
	}
	joined := strings.Join(allEvents, "\n")
	if strings.Contains(joined, `"type":"function_call"`) {
		t.Fatalf("custom tool stream leaked function_call events: %s", joined)
	}
	if strings.Contains(joined, "response.function_call_arguments.") {
		t.Fatalf("custom tool stream leaked function argument events: %s", joined)
	}
	if !strings.Contains(joined, `"type":"custom_tool_call"`) {
		t.Fatalf("missing custom_tool_call item: %s", joined)
	}
	if !strings.Contains(joined, "response.custom_tool_call_input.delta") {
		t.Fatalf("missing custom tool input delta: %s", joined)
	}
	if !strings.Contains(joined, `*** Add File: docs/test.md`) || strings.Contains(joined, `\n+\n*** End Patch`) {
		t.Fatalf("unexpected patch input: %s", joined)
	}
}

func TestConvertOpenAIChatToResponses_ToolSearchRemainsFunctionCall(t *testing.T) {
	ctx := context.Background()
	sseLines := []string{
		`data: {"id":"chatcmpl-tool-search","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","tool_calls":[{"index":0,"id":"call_search","type":"function","function":{"name":"tool_search","arguments":""}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-tool-search","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"query\":\"sub-agent\"}"}}]},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-tool-search","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`,
		`data: [DONE]`,
	}
	originalReq := []byte(`{"model":"gpt-4o","input":"find tools","tools":[{"type":"tool_search","execution":"client","parameters":{"type":"object","properties":{"query":{"type":"string"}},"required":["query"]}}],"transformer_metadata":{"codex_tool_compat_enabled":true}}`)

	var state any
	var allEvents []string
	for _, line := range sseLines {
		allEvents = append(allEvents, ConvertOpenAIChatToResponses(ctx, "gpt-4o", originalReq, nil, []byte(line), &state)...)
	}
	joined := strings.Join(allEvents, "\n")
	if strings.Contains(joined, `"type":"custom_tool_call"`) {
		t.Fatalf("tool_search must not be remapped to custom_tool_call: %s", joined)
	}
	if !strings.Contains(joined, `"type":"function_call"`) || !strings.Contains(joined, `"name":"tool_search"`) {
		t.Fatalf("missing function_call tool_search item: %s", joined)
	}
	if !strings.Contains(joined, "response.function_call_arguments.done") {
		t.Fatalf("missing function call arguments done event: %s", joined)
	}
}

func TestConvertOpenAIChatToResponses_Stream_ClaudeCacheTotalTokens(t *testing.T) {
	ctx := context.Background()
	sseLines := []string{
		`data: {"id":"msg-claude-cache","object":"chat.completion.chunk","created":1234567890,"model":"claude","choices":[{"index":0,"delta":{"role":"assistant","content":"ok"},"finish_reason":null}]}`,
		`data: {"id":"msg-claude-cache","object":"chat.completion.chunk","created":1234567890,"model":"claude","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"input_tokens":100,"output_tokens":20,"cache_creation_input_tokens":10,"cache_read_input_tokens":30}}`,
		`data: [DONE]`,
	}

	originalReq := []byte(`{"model":"claude","input":"hi"}`)
	var state any
	var allEvents []string
	for _, line := range sseLines {
		allEvents = append(allEvents, ConvertOpenAIChatToResponses(ctx, "claude", originalReq, nil, []byte(line), &state)...)
	}

	usage := extractResponseCompletedUsage(t, allEvents)
	if got := int(usage["total_tokens"].(float64)); got != 160 {
		t.Fatalf("total_tokens = %d, want 160", got)
	}
	if got := int(usage["cache_creation_input_tokens"].(float64)); got != 10 {
		t.Fatalf("cache_creation_input_tokens = %d, want 10", got)
	}
	if got := int(usage["cache_read_input_tokens"].(float64)); got != 30 {
		t.Fatalf("cache_read_input_tokens = %d, want 30", got)
	}
}

func TestConvertOpenAIChatToResponses_Stream_OpenAICacheDetailsNormalizesInput(t *testing.T) {
	ctx := context.Background()
	sseLines := []string{
		`data: {"id":"chatcmpl-openai-cache","object":"chat.completion.chunk","created":1234567890,"model":"gpt-5.5","choices":[{"index":0,"delta":{"role":"assistant","content":"ok"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-openai-cache","object":"chat.completion.chunk","created":1234567890,"model":"gpt-5.5","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":38451,"completion_tokens":1275,"total_tokens":39726,"prompt_tokens_details":{"cached_tokens":36608}}}`,
		`data: [DONE]`,
	}

	originalReq := []byte(`{"model":"gpt-5.5","input":"hi"}`)
	var state any
	var allEvents []string
	for _, line := range sseLines {
		allEvents = append(allEvents, ConvertOpenAIChatToResponses(ctx, "gpt-5.5", originalReq, nil, []byte(line), &state)...)
	}

	usage := extractResponseCompletedUsage(t, allEvents)
	if got := int(usage["input_tokens"].(float64)); got != 1843 {
		t.Fatalf("input_tokens = %d, want 1843", got)
	}
	if got := int(usage["total_tokens"].(float64)); got != 39726 {
		t.Fatalf("total_tokens = %d, want 39726", got)
	}
	if _, exists := usage["cache_read_input_tokens"]; exists {
		t.Fatalf("cache_read_input_tokens should not be emitted for OpenAI cache details: %#v", usage)
	}
	details, ok := usage["input_tokens_details"].(map[string]interface{})
	if !ok {
		t.Fatalf("input_tokens_details missing: %#v", usage)
	}
	if got := int(details["cached_tokens"].(float64)); got != 36608 {
		t.Fatalf("cached_tokens = %d, want 36608", got)
	}
}
