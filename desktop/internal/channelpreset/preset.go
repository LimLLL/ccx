package channelpreset

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

const (
	ProviderDeepSeek     = "deepseek"
	ProviderMiMo         = "mimo"
	ProviderCompshare    = "compshare"
	ProviderRunAPI       = "runapi"
	ProviderKimi         = "kimi"
	ProviderGLM          = "glm"
	ProviderMiniMax      = "minimax"
	ProviderDashScope    = "dashscope"
	ProviderTencentLkeap = "tencent-lkeap"
	ProviderKimiCode     = "kimi-code"
	ProviderVolcArk      = "volc-ark"
	ProviderQianfan      = "qianfan"
	ProviderOriginRouter = "originrouter"
	ProviderOpenRouter   = "openrouter"
	ProviderModelScope   = "modelscope"
	ProviderOpenCodeZen  = "opencode-zen"
	ProviderOpenCodeGo   = "opencode-go"

	TargetMessages  = "messages"
	TargetChat      = "chat"
	TargetResponses = "responses"
)

type ProviderPreset struct {
	ID                  string          `json:"id"`
	Order               int             `json:"order"`
	Label               string          `json:"label"`
	Description         string          `json:"description"`
	DirectAgent         bool            `json:"directAgent"`
	NativeMessages      bool            `json:"nativeMessages"`
	ChatCompatible      bool            `json:"chatCompatible"`
	ResponsesCompatible bool            `json:"responsesCompatible"`
	Plans               []ProviderPlan  `json:"plans"`
	Targets             []ChannelTarget `json:"targets"`
	DefaultTarget       string          `json:"defaultTarget"`
}

type ProviderPlan struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	BaseURL     string `json:"baseUrl"`
	Description string `json:"description"`
	Recommended bool   `json:"recommended"`
	Custom      bool   `json:"custom"`
}

// Protocol 返回 plan 所属的协议类型。
func (p ProviderPlan) Protocol() string {
	id := strings.ToLower(p.ID)
	label := strings.ToLower(p.Label)
	baseURL := strings.TrimRight(strings.ToLower(p.BaseURL), "/")
	if strings.Contains(id, "anthropic") || strings.Contains(label, "anthropic") || strings.Contains(baseURL, "anthropic") || strings.HasSuffix(baseURL, "/messages") {
		return "anthropic"
	}
	return "openai"
}

type ChannelTarget struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Recommended bool   `json:"recommended"`
}

type CreateChannelRequest struct {
	Provider    string `json:"provider"`
	Target      string `json:"target"`
	PlanID      string `json:"planId"`
	BaseURL     string `json:"baseUrl"`
	APIKey      string `json:"apiKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateChannelResult struct {
	Provider string `json:"provider"`
	Target   string `json:"target"`
	Name     string `json:"name"`
	BaseURL  string `json:"baseUrl"`
	Message  string `json:"message"`
}

type ChannelPayload struct {
	Name                          string            `json:"name"`
	Description                   string            `json:"description,omitempty"`
	Website                       string            `json:"website,omitempty"`
	ServiceType                   string            `json:"serviceType"`
	BaseURL                       string            `json:"baseUrl"`
	APIKeys                       []string          `json:"apiKeys"`
	ModelMapping                  map[string]string `json:"modelMapping,omitempty"`
	ReasoningMapping              map[string]string `json:"reasoningMapping,omitempty"`
	ReasoningParamStyle           string            `json:"reasoningParamStyle,omitempty"`
	PassbackReasoningContent      bool              `json:"passbackReasoningContent,omitempty"`
	PassbackThinkingBlocks        bool              `json:"passbackThinkingBlocks,omitempty"`
	NoVision                      bool              `json:"noVision,omitempty"`
	NoVisionModels                []string          `json:"noVisionModels,omitempty"`
	VisionFallbackModel           string            `json:"visionFallbackModel,omitempty"`
	SupportedModels               []string          `json:"supportedModels,omitempty"`
	StripEmptyTextBlocks          bool              `json:"stripEmptyTextBlocks,omitempty"`
	InjectDummyThoughtSignature   bool              `json:"injectDummyThoughtSignature,omitempty"`
	StripThoughtSignature         bool              `json:"stripThoughtSignature,omitempty"`
	CodexNativeToolPassthrough    bool              `json:"codexNativeToolPassthrough,omitempty"`
	CodexToolCompat               bool              `json:"codexToolCompat,omitempty"`
	StripCodexClientTools         bool              `json:"stripCodexClientTools,omitempty"`
	StripImageGenerationTool      bool              `json:"stripImageGenerationTool,omitempty"`
	NormalizeNonstandardChatRoles bool              `json:"normalizeNonstandardChatRoles,omitempty"`
	NormalizeMetadataUserId       *bool             `json:"normalizeMetadataUserId,omitempty"`
	RateLimitRPM                  int               `json:"rateLimitRpm,omitempty"`
	RateLimitBurst                int               `json:"rateLimitBurst,omitempty"`
	RateLimitMaxConcurrent        int               `json:"rateLimitMaxConcurrent,omitempty"`
	RateLimitAutoFromHeaders      bool              `json:"rateLimitAutoFromHeaders,omitempty"`
	Priority                      int               `json:"priority,omitempty"`
	Status                        string            `json:"status,omitempty"`
}

var providerConsoleURLs = map[string]string{
	ProviderDeepSeek:     "https://platform.deepseek.com/usage",
	ProviderMiMo:         "https://platform.xiaomimimo.com/console/balance",
	ProviderCompshare:    "https://console.compshare.cn/light-gpu/model-manage",
	ProviderRunAPI:       "https://runapi.co/console",
	ProviderKimi:         "https://platform.moonshot.cn/console/account",
	ProviderGLM:          "https://open.bigmodel.cn/coding-plan/personal/overview",
	ProviderMiniMax:      "https://platform.minimaxi.com/user-center/payment/balance",
	ProviderDashScope:    "https://bailian.console.aliyun.com/cn-beijing?tab=model#/api-key",
	ProviderTencentLkeap: "https://console.cloud.tencent.com/lkeap/token-plan",
	ProviderKimiCode:     "https://www.kimi.com/code/console",
	ProviderVolcArk:      "https://console.volcengine.com/ark",
	ProviderQianfan:      "https://console.bce.baidu.com/qianfan/resource/subscribe",
	ProviderOriginRouter: "https://easytransnote.com/ai/console/#key",
	ProviderOpenRouter:   "https://openrouter.ai/keys",
	ProviderModelScope:   "https://modelscope.cn/my/myaccesstoken",
	ProviderOpenCodeZen:  "https://opencode.ai/",
	ProviderOpenCodeGo:   "https://opencode.ai/",
}

func Presets() []ProviderPreset {
	return []ProviderPreset{
		{
			ID:                  ProviderDeepSeek,
			Order:               10,
			Label:               "DeepSeek",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传三种用法。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.deepseek.com/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.deepseek.com/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderMiMo,
			Order:               20,
			Label:               "MiMo",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传；内置按量与 token plan 入口。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "按量 (Anthropic)", BaseURL: "https://api.xiaomimimo.com/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "按量 (OpenAI)", BaseURL: "https://api.xiaomimimo.com/v1", Description: "Chat / Responses 通用入口"},
				{ID: "token-cn", Label: "Token Plan - 中国 (OpenAI)", BaseURL: "https://token-plan-cn.xiaomimimo.com/v1", Description: "中国区 Token Plan Chat / Responses 通用入口"},
				{ID: "token-sgp", Label: "Token Plan - 新加坡 (OpenAI)", BaseURL: "https://token-plan-sgp.xiaomimimo.com/v1", Description: "新加坡区 Token Plan Chat / Responses 通用入口"},
				{ID: "token-ams", Label: "Token Plan - 欧洲 (OpenAI)", BaseURL: "https://token-plan-ams.xiaomimimo.com/v1", Description: "欧洲区 Token Plan Chat / Responses 通用入口"},
				{ID: "token-cn-anthropic", Label: "Token Plan - 中国 (Anthropic)", BaseURL: "https://token-plan-cn.xiaomimimo.com/anthropic", Description: "中国区 Token Plan Claude Messages 原生入口"},
				{ID: "token-sgp-anthropic", Label: "Token Plan - 新加坡 (Anthropic)", BaseURL: "https://token-plan-sgp.xiaomimimo.com/anthropic", Description: "新加坡区 Token Plan Claude Messages 原生入口"},
				{ID: "token-ams-anthropic", Label: "Token Plan - 欧洲 (Anthropic)", BaseURL: "https://token-plan-ams.xiaomimimo.com/anthropic", Description: "欧洲区 Token Plan Claude Messages 原生入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderCompshare,
			Order:               30,
			Label:               "优云智算套餐",
			Description:         "优云智算是 UCloud 旗下 AI 云平台，提供高性价比国内 AI 模型 Agent Plan 套餐，支持包月订阅或按量付费（49 元/月起），同时提供官方海外模型稳定接入，支持 Claude Code、Codex 与 API 集成，具备企业级高并发、7×24 技术支持和自助开票能力；通过推广链接注册可领取 5 元平台试用金。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://cp.compshare.cn", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://cp.compshare.cn/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderRunAPI,
			Order:               40,
			Label:               "RunAPI",
			Description:         "RunAPI 是高效稳定的API OpenRouter平替平台，一个 API Key 即可访问 OpenAI、Claude、Gemini、DeepSeek、Grok 等 150+ 主流模型，低至 1 折，极其稳定，可以无缝兼容 Claude Code、OpenClaw 等工具。RunAPI 为 CCX用户提供专属福利：注册联系管理员即可领取￥7的免费额度",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://runapi.co/v1", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://runapi.co/v1", Description: "OpenAI Chat 兼容入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderKimi,
			Order:               50,
			Label:               "Kimi / Moonshot",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传三种用法。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.moonshot.cn/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.moonshot.cn/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderGLM,
			Order:               60,
			Label:               "GLM / BigModel",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传三种用法。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://open.bigmodel.cn/api/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "coding", Label: "Coding Plan (OpenAI)", BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4#", Description: "Coding Plan Chat / Responses 通用入口"},
				{ID: "openai-chat", Label: "通用 (OpenAI)", BaseURL: "https://open.bigmodel.cn/api/paas/v4#", Description: "通用 Chat / Responses 入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderMiniMax,
			Order:               70,
			Label:               "MiniMax",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传三种用法。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.minimaxi.com/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.minimax.chat/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderDashScope,
			Order:               80,
			Label:               "阿里云 DashScope",
			Description:         "Messages 原生透传、Codex Responses、Chat 渠道透传三种用法。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "按量 (Anthropic)", BaseURL: "https://dashscope.aliyuncs.com/apps/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "按量 (OpenAI)", BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1", Description: "Chat / Responses 通用入口"},
				{ID: "coding-anthropic", Label: "Coding Plan (Anthropic)", BaseURL: "https://coding.dashscope.aliyuncs.com/apps/anthropic", Description: "Coding Plan Claude Messages 原生入口"},
				{ID: "coding-openai-chat", Label: "Coding Plan (OpenAI)", BaseURL: "https://coding.dashscope.aliyuncs.com/v1", Description: "Coding Plan Chat / Responses 通用入口"},
				{ID: "token-plan-anthropic", Label: "Token Plan (Anthropic)", BaseURL: "https://token-plan.cn-beijing.maas.aliyuncs.com/apps/anthropic", Description: "Token Plan Claude Messages 原生入口"},
				{ID: "token-plan-openai-chat", Label: "Token Plan (OpenAI)", BaseURL: "https://token-plan.cn-beijing.maas.aliyuncs.com/compatible-mode/v1", Description: "Token Plan Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderOpenCodeZen,
			Order:               90,
			Label:               "OpenCode Zen",
			Description:         "按量付费精选模型网关，支持 Messages、Chat、Responses 三种协议。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://opencode.ai/zen/v1/messages", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://opencode.ai/zen/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderOpenCodeGo,
			Order:               100,
			Label:               "OpenCode Go",
			Description:         "低成本开源编程模型订阅服务（$5/月起），支持 Messages、Chat、Responses 三种协议。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://opencode.ai/zen/go/v1/messages", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://opencode.ai/zen/go/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderTencentLkeap,
			Order:               85,
			Label:               "腾讯云 TokenHub",
			Description:         "腾讯云大模型 TokenHub Token Plan 覆盖腾讯混元与国产主流模型，原生支持 Anthropic/OpenAI 双协议，适配 Claude Code、Codex、OpenCode、Cursor、Cline 等主流工具。",
			DirectAgent:         false,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.lkeap.cloud.tencent.com/plan/anthropic", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.lkeap.cloud.tencent.com/plan/v3", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderKimiCode,
			Order:               55,
			Label:               "Kimi Code Plan",
			Description:         "Kimi 会员 Coding Plan 编程套餐，原生支持 Anthropic/OpenAI 双协议，适配 Claude Code、OpenCode、Roo Code 等主流编码工具。",
			DirectAgent:         false,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.kimi.com/coding", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.kimi.com/coding/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderVolcArk,
			Order:               86,
			Label:               "火山方舟 Coding Plan",
			Description:         "火山引擎方舟 Coding Plan（豆包）订阅套餐，覆盖 Doubao、DeepSeek、Kimi、GLM、MiniMax 等模型，原生支持 Anthropic/OpenAI 双协议，适配 Claude Code、Codex、OpenCode、Cursor 等工具。",
			DirectAgent:         false,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://ark.cn-beijing.volces.com/api/coding", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://ark.cn-beijing.volces.com/api/coding/v3", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderQianfan,
			Order:               87,
			Label:               "百度千帆 Coding Plan",
			Description:         "百度智能云千帆 Coding Plan 订阅套餐，覆盖文心、DeepSeek、GLM、Kimi、MiniMax 等模型，原生支持 Anthropic/OpenAI 双协议，适配 Claude Code、Codex、OpenCode、Cursor 等工具。",
			DirectAgent:         false,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://qianfan.baidubce.com/anthropic/coding", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://qianfan.baidubce.com/v2/coding", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderOpenRouter,
			Order:               35,
			Label:               "OpenRouter",
			Description:         "全球最大模型聚合平台，300+ 模型一个 Key 通吃，原生支持 Messages、Chat、Responses 三种协议。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://openrouter.ai/api", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://openrouter.ai/api/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderModelScope,
			Order:               82,
			Label:               "ModelScope 魔搭",
			Description:         "阿里达摩院开源模型社区，提供 Qwen、DeepSeek、Llama 等模型免费推理 API，原生支持 Messages、Chat、Responses 三种协议。",
			DirectAgent:         true,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api-inference.modelscope.cn", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api-inference.modelscope.cn/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
			},
			DefaultTarget: TargetMessages,
		},
		{
			ID:                  ProviderOriginRouter,
			Order:               200,
			Label:               "极易云 OriginRouter",
			Description:         "极易云 OriginRouter Coding Plan 编程套餐，原生支持 Anthropic/OpenAI 双协议，兼容 Claude Code、OpenCode、Cline 等工具。",
			DirectAgent:         false,
			NativeMessages:      true,
			ChatCompatible:      true,
			ResponsesCompatible: true,
			Plans: []ProviderPlan{
				{ID: "anthropic", Label: "Anthropic-compatible", BaseURL: "https://api.easytransnote.com/coding", Description: "Claude Messages 原生入口", Recommended: true},
				{ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://api.easytransnote.com/coding/v1", Description: "Chat / Responses 通用入口"},
			},
			Targets: []ChannelTarget{
				{Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
				{Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
				{Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
			},
			DefaultTarget: TargetMessages,
		},
	}
}

func BuildPayload(req CreateChannelRequest) (ChannelPayload, error) {
	preset, ok := FindPreset(req.Provider)
	if !ok {
		return ChannelPayload{}, fmt.Errorf("不支持的 provider: %s", req.Provider)
	}
	target := strings.TrimSpace(req.Target)
	if target == "" {
		target = preset.DefaultTarget
	}
	if !supportsTarget(preset, target) {
		return ChannelPayload{}, fmt.Errorf("%s 不支持添加到 %s 渠道", preset.Label, target)
	}
	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		return ChannelPayload{}, fmt.Errorf("API Key 不能为空")
	}
	planID := strings.TrimSpace(req.PlanID)
	if planID == "" {
		planID = bestPlanForTarget(preset, target)
	}
	baseURL, err := ResolveBaseURL(preset, planID, req.BaseURL)
	if err != nil {
		return ChannelPayload{}, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = defaultChannelName(preset.ID, target)
	}
	payload := ChannelPayload{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Website:     providerConsoleURLs[preset.ID],
		BaseURL:     baseURL,
		APIKeys:     []string{apiKey},
		Priority:    1,
		Status:      "active",
	}
	applyTargetDefaults(&payload, preset.ID, target)
	return payload, nil
}

// targetMatchesURL 判断 URL 是否与 target 协议兼容。
// messages target 使用 Anthropic 协议，需要 anthropic 路径；
// chat/responses target 使用 OpenAI 协议，不应使用 anthropic 路径。

// bestPlanForTarget 根据 target 自动选择最合适的 plan。
// 当 provider 有多个 plan（如 DeepSeek 的 /anthropic 和 /v1）时，
// 确保 chat/responses 选择 OpenAI-compatible 入口，messages 选择 Anthropic 入口。
func bestPlanForTarget(preset ProviderPreset, target string) string {
	if len(preset.Plans) == 0 {
		return ""
	}
	if len(preset.Plans) == 1 {
		return preset.Plans[0].ID
	}
	for _, plan := range preset.Plans {
		if !plan.Custom {
			isAnthropic := plan.Protocol() == "anthropic"
			if target == TargetMessages && isAnthropic {
				return plan.ID
			}
			if (target == TargetChat || target == TargetResponses) && !isAnthropic {
				return plan.ID
			}
		}
	}
	for _, plan := range preset.Plans {
		if plan.Recommended {
			return plan.ID
		}
	}
	return preset.Plans[0].ID
}

func FindPreset(provider string) (ProviderPreset, bool) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	for _, preset := range Presets() {
		if preset.ID == provider {
			return preset, true
		}
	}
	return ProviderPreset{}, false
}

// FilterPlansForTarget 按 target 协议过滤 plans，只返回兼容的入口。
// messages 使用 Anthropic 协议，chat/responses 使用 OpenAI 协议。
// 自定义 plan 始终保留。
func FilterPlansForTarget(plans []ProviderPlan, target string) []ProviderPlan {
	if len(plans) <= 1 {
		return plans
	}
	wantAnthropic := target == TargetMessages
	var filtered []ProviderPlan
	for _, plan := range plans {
		isAnthropic := plan.Protocol() == "anthropic"
		if wantAnthropic == isAnthropic {
			filtered = append(filtered, plan)
		}
	}
	if len(filtered) == 0 {
		return plans
	}
	return filtered
}

func ResolveBaseURL(preset ProviderPreset, planID string, customBaseURL string) (string, error) {
	customBaseURL = strings.TrimSpace(customBaseURL)
	planID = strings.TrimSpace(planID)
	if planID == "" {
		for _, plan := range preset.Plans {
			if plan.Recommended {
				planID = plan.ID
				break
			}
		}
	}
	if planID == "" && len(preset.Plans) > 0 {
		planID = preset.Plans[0].ID
	}
	for _, plan := range preset.Plans {
		if plan.ID != planID {
			continue
		}
		if plan.Custom {
			if customBaseURL == "" {
				return "", fmt.Errorf("自定义 plan 需要填写 Base URL")
			}
			return customBaseURL, nil
		}
		if customBaseURL != "" {
			return customBaseURL, nil
		}
		if plan.BaseURL == "" {
			return "", fmt.Errorf("plan %s 缺少 Base URL", plan.ID)
		}
		return plan.BaseURL, nil
	}
	if customBaseURL != "" {
		return customBaseURL, nil
	}
	return "", fmt.Errorf("未找到 provider %s 的 plan: %s", preset.ID, planID)
}

func supportsTarget(preset ProviderPreset, target string) bool {
	for _, item := range preset.Targets {
		if item.Type == target {
			return true
		}
	}
	return false
}

func defaultChannelName(provider string, target string) string {
	return fmt.Sprintf("desktop-%s-%s", provider, target)
}

type channelTargetConfig struct {
	ModelMapping                  map[string]string
	ReasoningMapping              map[string]string
	ReasoningParamStyle           string
	PassbackReasoningContent      bool
	PassbackThinkingBlocks        bool
	NoVision                      bool
	NoVisionModels                []string
	VisionFallbackModel           string
	StripEmptyTextBlocks          bool
	CodexNativeToolPassthrough    bool
	CodexToolCompat               *bool
	StripCodexClientTools         *bool
	NormalizeNonstandardChatRoles bool
	RateLimitRPM                  int // 主动限速默认 RPM（0=不设默认）
}

var channelTargetConfigs = map[string]map[string]channelTargetConfig{
	TargetMessages: {
		ProviderDeepSeek: {
			ModelMapping: map[string]string{
				"haiku":  "deepseek-v4-flash",
				"opus":   "deepseek-v4-pro",
				"sonnet": "deepseek-v4-pro",
			},
			ReasoningParamStyle:      "reasoning",
			PassbackReasoningContent: true,
			PassbackThinkingBlocks:   true,
			NoVision:                 true,
		},
		ProviderMiMo: {
			ModelMapping: map[string]string{
				"haiku":  "mimo-v2.5-pro",
				"opus":   "mimo-v2.5-pro",
				"sonnet": "mimo-v2.5-pro",
			},
			ReasoningParamStyle:      "reasoning",
			PassbackReasoningContent: true,
			NoVisionModels:           []string{"mimo-v2.5-pro"},
			VisionFallbackModel:      "mimo-v2.5",
			RateLimitRPM:             80, // MiMo 官方 RPM=100（账号级、跨 key 共享），留 20% 余量
		},
		ProviderCompshare: {
			ModelMapping: map[string]string{
				"haiku":  "deepseek-v4-flash",
				"opus":   "glm-5.1",
				"sonnet": "glm-5.1",
			},
			ReasoningParamStyle:      "reasoning",
			PassbackReasoningContent: true,
			NoVisionModels:           []string{"deepseek-v4-flash"},
			VisionFallbackModel:      "MiniMax-M2.7",
		},
		ProviderRunAPI: {},
		ProviderKimi: {
			ModelMapping: map[string]string{
				"haiku":  "kimi-k2.6",
				"opus":   "kimi-k2.6",
				"sonnet": "kimi-k2.6",
			},
		},
		ProviderGLM: {
			ModelMapping: map[string]string{
				"haiku":  "glm-5.1",
				"opus":   "glm-5.1",
				"sonnet": "glm-5.1",
			},
		},
		ProviderMiniMax: {
			ModelMapping: map[string]string{
				"haiku":  "MiniMax-M2.7",
				"opus":   "MiniMax-M2.7",
				"sonnet": "MiniMax-M2.7",
			},
			PassbackReasoningContent: true,
		},
		ProviderDashScope: {
			ModelMapping: map[string]string{
				"haiku":  "glm-5.1",
				"opus":   "glm-5.1",
				"sonnet": "glm-5.1",
			},
		},
		ProviderOpenCodeZen: {
			ModelMapping: map[string]string{
				"haiku":  "glm-5.1",
				"opus":   "glm-5.1",
				"sonnet": "glm-5.1",
			},
		},
		ProviderOpenCodeGo: {
			ModelMapping: map[string]string{
				"haiku":  "glm-5.1",
				"opus":   "glm-5.1",
				"sonnet": "glm-5.1",
			},
		},
		ProviderModelScope: {
			ModelMapping: map[string]string{
				"haiku":  "deepseek-ai/DeepSeek-V4-Flash",
				"sonnet": "ZhipuAI/GLM-5.1",
				"opus":   "ZhipuAI/GLM-5.1",
			},
		},
	},
	TargetChat: {
		ProviderDeepSeek: {
			ReasoningParamStyle: "reasoning",
			NoVision:            true,
		},
		ProviderMiMo: {
			ReasoningParamStyle: "reasoning",
			NoVisionModels:      []string{"mimo-v2.5-pro"},
			VisionFallbackModel: "mimo-v2.5",
			RateLimitRPM:        80,
		},
		ProviderCompshare: {
			ReasoningParamStyle: "reasoning",
			NoVisionModels:      []string{"deepseek-v4-flash"},
			VisionFallbackModel: "MiniMax-M2.7",
		},
		ProviderRunAPI:      {},
		ProviderOpenRouter:  {},
		ProviderModelScope: {
			NormalizeNonstandardChatRoles: true,
		},
		ProviderMiniMax:     {},
		ProviderDashScope:   {},
		ProviderOpenCodeZen: {},
		ProviderOpenCodeGo:  {},
	},
	TargetResponses: {
		ProviderDeepSeek: {
			ModelMapping: map[string]string{
				"gpt":               "deepseek-v4-pro",
				"mini":              "deepseek-v4-flash",
				"codex-auto-review": "deepseek-v4-flash",
			},
			ReasoningMapping:              map[string]string{"gpt": "max"},
			ReasoningParamStyle:           "reasoning",
			CodexToolCompat:               boolRef(false),
			StripCodexClientTools:         boolRef(false),
			CodexNativeToolPassthrough:    true,
			NormalizeNonstandardChatRoles: true,
			NoVision:                      true,
		},
		ProviderMiMo: {
			ModelMapping:          map[string]string{"gpt-5": "mimo-v2.5-pro", "codex-auto-review": "mimo-v2.5"},
			ReasoningParamStyle:   "reasoning",
			CodexToolCompat:       boolRef(false),
			StripCodexClientTools: boolRef(false),
			NoVisionModels:        []string{"mimo-v2.5-pro"},
			VisionFallbackModel:   "mimo-v2.5",
			RateLimitRPM:          80,
		},
		ProviderCompshare: {
			ModelMapping: map[string]string{
				"gpt":               "glm-5.1",
				"mini":              "deepseek-v4-flash",
				"codex-auto-review": "deepseek-v4-flash",
			},
			ReasoningParamStyle:           "reasoning",
			CodexToolCompat:               boolRef(false),
			StripCodexClientTools:         boolRef(false),
			CodexNativeToolPassthrough:    true,
			NormalizeNonstandardChatRoles: true,
			NoVisionModels:                []string{"deepseek-v4-flash"},
			VisionFallbackModel:           "MiniMax-M2.7",
		},
		ProviderRunAPI: {
			CodexToolCompat:       boolRef(false),
			StripCodexClientTools: boolRef(false),
		},
		ProviderMiniMax: {
			ModelMapping:                  map[string]string{"gpt-5": "MiniMax-M2.7", "codex-auto-review": "MiniMax-M2.7"},
			CodexToolCompat:               boolRef(false),
			StripCodexClientTools:         boolRef(false),
			CodexNativeToolPassthrough:    true,
			NormalizeNonstandardChatRoles: true,
		},
		ProviderDashScope: {
			ModelMapping: map[string]string{
				"gpt-5.5":           "glm-5.1",
				"gpt-5.4":           "deepseek-v4-pro",
				"gpt-5.4-mini":      "deepseek-v4-flash",
				"codex-auto-review": "deepseek-v4-flash",
			},
			ReasoningMapping: map[string]string{
				"gpt-5.5":      "high",
				"gpt-5.4":      "max",
				"gpt-5.4-mini": "high",
			},
		},
		ProviderOpenCodeZen: {
			ModelMapping: map[string]string{"gpt-5": "glm-5.1", "codex-auto-review": "glm-5.1"},
		},
		ProviderOpenCodeGo: {
			ModelMapping: map[string]string{"gpt-5": "glm-5.1", "codex-auto-review": "glm-5.1"},
		},
		ProviderKimi: {
			ModelMapping: map[string]string{"gpt-5": "kimi-k2.6", "codex-auto-review": "kimi-k2.6"},
		},
		ProviderGLM: {
			ModelMapping: map[string]string{"gpt-5": "glm-5.1", "codex-auto-review": "glm-5.1"},
		},
		ProviderOpenRouter: {
			CodexToolCompat:       boolRef(false),
			StripCodexClientTools: boolRef(false),
		},
		ProviderModelScope: {
			ModelMapping: map[string]string{
				"gpt":               "ZhipuAI/GLM-5.1",
				"mini":              "deepseek-ai/DeepSeek-V4-Flash",
				"codex-auto-review": "deepseek-ai/DeepSeek-V4-Flash",
			},
			CodexToolCompat:               boolRef(false),
			StripCodexClientTools:         boolRef(false),
			NormalizeNonstandardChatRoles: true,
		},
	},
}

func boolRef(value bool) *bool {
	return &value
}

func applyChannelTargetConfig(payload *ChannelPayload, config channelTargetConfig) {
	payload.ModelMapping = maps.Clone(config.ModelMapping)
	payload.ReasoningMapping = maps.Clone(config.ReasoningMapping)
	payload.NoVisionModels = slices.Clone(config.NoVisionModels)
	payload.ReasoningParamStyle = config.ReasoningParamStyle
	payload.PassbackReasoningContent = payload.PassbackReasoningContent || config.PassbackReasoningContent
	payload.PassbackThinkingBlocks = payload.PassbackThinkingBlocks || config.PassbackThinkingBlocks
	payload.NoVision = payload.NoVision || config.NoVision
	payload.VisionFallbackModel = config.VisionFallbackModel
	payload.StripEmptyTextBlocks = payload.StripEmptyTextBlocks || config.StripEmptyTextBlocks
	payload.CodexNativeToolPassthrough = payload.CodexNativeToolPassthrough || config.CodexNativeToolPassthrough
	payload.NormalizeNonstandardChatRoles = payload.NormalizeNonstandardChatRoles || config.NormalizeNonstandardChatRoles
	if config.CodexToolCompat != nil {
		payload.CodexToolCompat = *config.CodexToolCompat
	}
	if config.StripCodexClientTools != nil {
		payload.StripCodexClientTools = *config.StripCodexClientTools
	}
	if config.RateLimitRPM > 0 {
		payload.RateLimitRPM = config.RateLimitRPM
	}
}

func applyTargetDefaults(payload *ChannelPayload, provider string, target string) {
	switch target {
	case TargetMessages:
		payload.ServiceType = "claude"
		payload.StripEmptyTextBlocks = true
		payload.StripThoughtSignature = true
		if provider == ProviderRunAPI {
			payload.NormalizeMetadataUserId = boolRef(false)
		}
	case TargetChat:
		payload.ServiceType = "openai"
		payload.NormalizeNonstandardChatRoles = true
	case TargetResponses:
		payload.ServiceType = "openai"
		payload.CodexToolCompat = true
		payload.StripCodexClientTools = true
		if provider == ProviderRunAPI {
			payload.ServiceType = "responses"
			payload.CodexToolCompat = false
			payload.StripCodexClientTools = false
		}
	}

	configs, ok := channelTargetConfigs[target]
	if !ok {
		return
	}
	config, ok := configs[provider]
	if !ok {
		return
	}
	applyChannelTargetConfig(payload, config)
}
