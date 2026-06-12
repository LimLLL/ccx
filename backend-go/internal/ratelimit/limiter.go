package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ChannelLimiter 是单个渠道的限速器，组合了令牌桶、并发信号量和动态 cooldown。
// 零值表示不限速（所有字段为 0 时 Acquire 立即成功）。
type ChannelLimiter struct {
	mu sync.Mutex

	// --- 令牌桶 ---
	// tokens 是当前可用令牌数（浮点，允许小步累积）。
	// rate 是每秒填充速率（= RPM / 60）；burst 是桶容量（最大令牌数）。
	// rate=0 || burst=0 时令牌桶不生效（始终放行）。
	tokens float64
	rate   float64
	burst  float64
	// lastRefill 是上一次补充令牌的时间，用于按时间差计算累积。
	lastRefill time.Time

	// --- 并发信号量 ---
	// maxConcurrent=0 表示不限并发。
	sem chan struct{}

	// --- 动态 cooldown ---
	// cooldownUntil 非零且在当前时间之前时，acquire 直接快速失败。
	cooldownUntil time.Time
}

// Config 是 ChannelLimiter 的创建/更新配置。
type Config struct {
	// RPM 是每分钟请求数上限。0=不限。
	RPM int
	// Burst 是令牌桶容量（允许的突发请求数）。0=不限（取 RPM 值的默认 burst）。
	Burst int
	// MaxConcurrent 是最大并发上游请求数。0=不限。
	MaxConcurrent int
	// AutoFromHeaders 是否自动从上游响应头解析限流信息。默认 false。
	AutoFromHeaders bool
}

// errors
var (
	ErrInCooldown  = fmt.Errorf("rate limited: channel is in cooldown")
	ErrAcquireBusy = fmt.Errorf("rate limited: max concurrent requests reached")
	ErrBucketEmpty = fmt.Errorf("rate limited: token bucket exhausted")
)

// NewChannelLimiter 创建一个新的 ChannelLimiter。now 用于初始化令牌桶时间基准。
func NewChannelLimiter(cfg Config, now time.Time) *ChannelLimiter {
	l := &ChannelLimiter{
		lastRefill: now,
	}
	l.applyConfig(cfg)
	return l
}

// applyConfig 将配置应用到 limiter；保留运行态的 tokens/cooldown。
func (l *ChannelLimiter) applyConfig(cfg Config) {
	if cfg.RPM > 0 {
		l.rate = float64(cfg.RPM) / 60.0
		if cfg.Burst > 0 {
			l.burst = float64(cfg.Burst)
		} else {
			// 默认 burst = 1（强制平滑限速，避免瞬间突发导致上游 429）
			l.burst = 1
		}
		// 初始填满
		if l.tokens > l.burst {
			l.tokens = l.burst
		} else if l.tokens == 0 {
			l.tokens = l.burst
		}
	} else {
		l.rate = 0
		l.burst = 0
		l.tokens = 0
	}

	if cfg.MaxConcurrent > 0 {
		// 重新分配信号量：如果有变更需重建
		if l.sem == nil || cap(l.sem) != cfg.MaxConcurrent {
			l.sem = make(chan struct{}, cfg.MaxConcurrent)
		}
	} else {
		l.sem = nil
	}
}

// UpdateConfig 热更新配置，不丢失运行态 cooldown 和当前令牌。
func (l *ChannelLimiter) UpdateConfig(cfg Config) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.applyConfig(cfg)
}

// InCooldown 返回当前是否处于动态 cooldown 期。
func (l *ChannelLimiter) InCooldown(now time.Time) (bool, time.Time) {
	if l == nil {
		return false, time.Time{}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.cooldownUntil.IsZero() || !now.Before(l.cooldownUntil) {
		return false, time.Time{}
	}
	return true, l.cooldownUntil
}

// Acquire 尝试获取一个请求许可。返回 release 函数（必须在请求完成后调用以释放并发信号量）。
// maxWait 是最长排队等待时间；ctx 支持客户端断开取消。
// 返回的 error 可能是 ErrInCooldown / ErrAcquireBusy / ErrBucketEmpty / context.Canceled。
func (l *ChannelLimiter) Acquire(ctx context.Context, maxWait time.Duration, now time.Time) (release func(), err error) {
	if l == nil {
		return func() {}, nil
	}

	// 1. 检查 cooldown
	if released, ok := l.tryCooldown(now); !ok {
		return released, ErrInCooldown
	}

	// 2. 获取令牌（等待期间不占用并发信号量，避免排队请求挤占并发槽位）
	if err := l.acquireToken(ctx, maxWait, now); err != nil {
		return func() {}, err
	}

	// 3. 获取并发信号量（拿到令牌后才占槽，确保信号量反映真实在途请求数）
	release, err = l.acquireSemaphore(ctx, maxWait)
	if err != nil {
		return func() {}, err
	}

	return release, nil
}

// tryCooldown 检查并返回是否可继续。不可继续时返回一个空 release + false。
func (l *ChannelLimiter) tryCooldown(now time.Time) (release func(), ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cooldownUntil.IsZero() || !now.Before(l.cooldownUntil) {
		return func() {}, true
	}
	return func() {}, false
}

// acquireSemaphore 获取并发信号量，支持 ctx 取消和 maxWait 超时。
func (l *ChannelLimiter) acquireSemaphore(ctx context.Context, maxWait time.Duration) (func(), error) {
	if l.sem == nil {
		return func() {}, nil
	}

	// 快速尝试（非阻塞）
	select {
	case l.sem <- struct{}{}:
		return l.makeSemaphoreRelease(), nil
	default:
	}

	// 需要等待
	deadline := time.After(maxWait)
	for {
		select {
		case <-ctx.Done():
			return func() {}, ctx.Err()
		case <-deadline:
			return func() {}, ErrAcquireBusy
		case l.sem <- struct{}{}:
			return l.makeSemaphoreRelease(), nil
		}
	}
}

func (l *ChannelLimiter) makeSemaphoreRelease() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			<-l.sem
		})
	}
}

// acquireToken 从令牌桶扣减一个令牌。需要等待时循环 sleep，支持 ctx/timeout 取消。
func (l *ChannelLimiter) acquireToken(ctx context.Context, maxWait time.Duration, now time.Time) error {
	l.mu.Lock()

	if l.rate <= 0 || l.burst <= 0 {
		// 令牌桶不生效
		l.mu.Unlock()
		return nil
	}

	// 按时间差补充令牌
	l.refillLocked(now)

	if l.tokens >= 1 {
		l.tokens--
		l.mu.Unlock()
		return nil
	}

	// 计算等待到下一个令牌的时间
	waitDuration := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
	l.mu.Unlock()

	if waitDuration > maxWait {
		return ErrBucketEmpty
	}

	deadline := time.After(maxWait)
	ticker := time.NewTicker(waitDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return ErrBucketEmpty
		case <-ticker.C:
			l.mu.Lock()
			l.refillLocked(time.Now())
			if l.tokens >= 1 {
				l.tokens--
				l.mu.Unlock()
				return nil
			}
			// 还不够，缩小等待间隔
			nextWait := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
			if nextWait < 10*time.Millisecond {
				nextWait = 10 * time.Millisecond
			}
			ticker.Reset(nextWait)
			l.mu.Unlock()
		}
	}
}

// refillLocked 按时间差补充令牌（需持有锁）。
func (l *ChannelLimiter) refillLocked(now time.Time) {
	if l.lastRefill.IsZero() {
		l.lastRefill = now
		return
	}
	elapsed := now.Sub(l.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}
	l.lastRefill = now
}

// Status 返回当前 limiter 状态快照（用于调试/日志）。
func (l *ChannelLimiter) Status(now time.Time) LimiterStatus {
	if l == nil {
		return LimiterStatus{}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refillLocked(now)

	inCooldown := !l.cooldownUntil.IsZero() && now.Before(l.cooldownUntil)
	semUsed := 0
	if l.sem != nil {
		semUsed = len(l.sem)
	}
	return LimiterStatus{
		Tokens:          l.tokens,
		Burst:           l.burst,
		RatePerSec:      l.rate,
		MaxConcurrent:   cap(l.sem),
		ActiveRequests:  semUsed,
		InCooldown:      inCooldown,
		CooldownUntil:   l.cooldownUntil,
		AutoFromHeaders: false, // 由 Manager 层设置
	}
}

// LimiterStatus 是 limiter 的只读快照。
type LimiterStatus struct {
	Tokens         float64
	Burst          float64
	RatePerSec     float64
	MaxConcurrent  int
	ActiveRequests int
	InCooldown     bool
	CooldownUntil  time.Time

	AutoFromHeaders bool
}
