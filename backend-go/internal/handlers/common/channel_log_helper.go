// Package common 提供 handlers 模块的公共功能
package common

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/BenedictKing/ccx/internal/metrics"
	"github.com/BenedictKing/ccx/internal/utils"
)

// GenerateRequestID 生成唯一的请求标识
func GenerateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// CreatePendingLog 创建 pending 状态的日志条目（请求开始时调用）
func CreatePendingLog(
	channelLogStore *metrics.ChannelLogStore,
	metricsKey string,
	channelIndex int,
	model, originalModel string,
	apiKey, baseURL, interfaceType, operation string,
	requestSource string,
) string {
	if channelLogStore == nil || metricsKey == "" {
		return ""
	}
	if requestSource == "" {
		requestSource = metrics.RequestSourceProxy
	}

	requestID := GenerateRequestID()
	now := time.Now()

	channelLogStore.Record(metricsKey, &metrics.ChannelLog{
		RequestID:     requestID,
		ChannelIndex:  channelIndex, // 记录创建时的渠道索引，便于内部排查
		Timestamp:     now,
		StartTime:     now,
		Model:         model,
		OriginalModel: originalModel,
		Operation:     operation,
		StatusCode:    0,
		DurationMs:    0,
		Success:       false,
		KeyMask:       utils.MaskAPIKey(apiKey),
		BaseURL:       baseURL,
		ErrorInfo:     "",
		IsRetry:       false,
		InterfaceType: interfaceType,
		RequestSource: requestSource,
		Status:        metrics.StatusPending,
	})

	return requestID
}

// UpdateLogStatus 更新日志状态（连接建立、首字节、流式传输等）
func UpdateLogStatus(
	channelLogStore *metrics.ChannelLogStore,
	metricsKey string,
	requestID string,
	status string,
) {
	if channelLogStore == nil || metricsKey == "" || requestID == "" {
		return
	}

	now := time.Now()
	channelLogStore.Update(metricsKey, requestID, func(log *metrics.ChannelLog) {
		log.Status = status
		switch status {
		case metrics.StatusConnecting:
			log.ConnectedAt = &now
		case metrics.StatusFirstByte:
			log.FirstByteAt = &now
		case metrics.StatusStreaming:
			if log.FirstByteAt == nil {
				log.FirstByteAt = &now
			}
		}
	})
}

// CompleteLog 完成日志记录（请求结束时调用）
func CompleteLog(
	channelLogStore *metrics.ChannelLogStore,
	metricsKey string,
	requestID string,
	statusCode int,
	success bool,
	errorInfo string,
	isRetry bool,
) {
	if channelLogStore == nil || metricsKey == "" || requestID == "" {
		return
	}

	errorInfo = normalizeChannelLogErrorInfo(errorInfo)
	if len(errorInfo) > 200 {
		errorInfo = errorInfo[:200]
	}

	status := getStatusFromResult(success, errorInfo)
	now := time.Now()
	updateStatus, actualMetricsKey := channelLogStore.Update(metricsKey, requestID, func(log *metrics.ChannelLog) {
		log.StatusCode = statusCode
		log.Success = success
		log.ErrorInfo = errorInfo
		log.IsRetry = isRetry
		log.CompletedAt = &now
		log.DurationMs = now.Sub(log.StartTime).Milliseconds()
		log.Status = status
	})

	// 仅在确认是环形缓冲淘汰时补写终态日志；若渠道已删除则不补写，避免污染其他渠道。
	if updateStatus == metrics.UpdateMissingEvicted && actualMetricsKey != "" {
		channelLogStore.Record(actualMetricsKey, &metrics.ChannelLog{
			RequestID:   requestID,
			Timestamp:   now,
			StatusCode:  statusCode,
			Success:     success,
			ErrorInfo:   errorInfo,
			IsRetry:     isRetry,
			Status:      status,
			StartTime:   now,
			CompletedAt: &now,
			DurationMs:  0,
		})
	}
}

func getStatusFromResult(success bool, errorInfo string) string {
	if success {
		return metrics.StatusCompleted
	}
	if strings.EqualFold(strings.TrimSpace(errorInfo), "client canceled") {
		return metrics.StatusCancelled
	}
	return metrics.StatusFailed
}

// RecordChannelLog 统一记录渠道尝试日志。
// 约束：凡是会进入渠道图表统计的尝试，都应该调用此函数或等价逻辑写入日志，保证图表与日志口径一致。
// 注意：此函数用于兼容旧代码，新代码应使用 CreatePendingLog + UpdateLogStatus + CompleteLog 组合
func RecordChannelLog(
	channelLogStore *metrics.ChannelLogStore,
	metricsKey string,
	channelIndex int,
	model, originalModel string,
	statusCode int,
	durationMs int64,
	success bool,
	apiKey, baseURL, errorInfo, interfaceType string,
	isRetry bool,
) {
	RecordChannelLogWithSource(
		channelLogStore,
		metricsKey,
		channelIndex,
		model,
		originalModel,
		statusCode,
		durationMs,
		success,
		apiKey,
		baseURL,
		errorInfo,
		interfaceType,
		isRetry,
		metrics.RequestSourceProxy,
	)
}

// RecordChannelLogWithSource 记录带来源标识的渠道尝试日志。
// 注意：此函数用于兼容旧代码，新代码应使用 CreatePendingLog + UpdateLogStatus + CompleteLog 组合
func RecordChannelLogWithSource(
	channelLogStore *metrics.ChannelLogStore,
	metricsKey string,
	channelIndex int,
	model, originalModel string,
	statusCode int,
	durationMs int64,
	success bool,
	apiKey, baseURL, errorInfo, interfaceType string,
	isRetry bool,
	requestSource string,
) {
	if channelLogStore == nil || metricsKey == "" {
		return
	}
	errorInfo = normalizeChannelLogErrorInfo(errorInfo)
	if len(errorInfo) > 200 {
		errorInfo = errorInfo[:200]
	}
	if requestSource == "" {
		requestSource = metrics.RequestSourceProxy
	}

	now := time.Now()
	startTime := now.Add(-time.Duration(durationMs) * time.Millisecond)
	requestID := GenerateRequestID()

	var status string
	if success {
		status = metrics.StatusCompleted
	} else {
		status = metrics.StatusFailed
	}

	channelLogStore.Record(metricsKey, &metrics.ChannelLog{
		RequestID:     requestID,
		ChannelIndex:  channelIndex, // 记录创建时的渠道索引，便于内部排查
		Timestamp:     now,
		StartTime:     startTime,
		Model:         model,
		OriginalModel: originalModel,
		StatusCode:    statusCode,
		DurationMs:    durationMs,
		Success:       success,
		KeyMask:       utils.MaskAPIKey(apiKey),
		BaseURL:       baseURL,
		ErrorInfo:     errorInfo,
		IsRetry:       isRetry,
		InterfaceType: interfaceType,
		RequestSource: requestSource,
		Status:        status,
		CompletedAt:   &now,
	})
}

func normalizeChannelLogErrorInfo(errorInfo string) string {
	trimmed := strings.TrimSpace(errorInfo)
	switch {
	case strings.HasPrefix(trimmed, ErrEmptyStreamResponse.Error()):
		diagnostic := strings.TrimSpace(strings.TrimPrefix(trimmed, ErrEmptyStreamResponse.Error()))
		diagnostic = strings.TrimSpace(strings.TrimPrefix(diagnostic, ":"))
		if diagnostic != "" {
			return "空流响应：上游 HTTP 200 返回 SSE 流后结束，但未检测到文本或语义内容（" + diagnostic + "）"
		}
		return "空流响应：上游 HTTP 200 返回 SSE 流后结束，但未检测到文本或语义内容"
	case strings.HasPrefix(trimmed, ErrEmptyNonStreamResponse.Error()):
		return "空响应：上游 HTTP 200 返回非流式响应，但未检测到文本或语义内容"
	case strings.HasPrefix(trimmed, ErrStreamFirstContentTimeout.Error()):
		return "流式首内容超时：上游 HTTP 200 后未在配置窗口内返回有效内容"
	case strings.HasPrefix(trimmed, ErrStreamStalled.Error()):
		return "流式断流：首个有效内容后未在配置窗口内继续返回上游活动"
	default:
		return errorInfo
	}
}
