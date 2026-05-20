// Package updater 提供基于 GitHub Releases 的桌面端自动更新能力。
//
// 工作流：
//  1. Check        查询 GitHub 最新 release tag，对比当前版本
//  2. Download     下载平台对应资产到临时目录
//  3. Verify       校验 SHA256
//  4. Install      调用平台特定安装器（DMG / NSIS / AppImage）
//
// 进度通过 Subscribe() 返回的 channel 推送，UI 层订阅即可。
package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultRepo    = "BenedictKing/ccx"
	apiBase        = "https://api.github.com"
	userAgent      = "ccx-desktop-updater"
	httpTimeout    = 30 * time.Second
	downloadBuffer = 32 * 1024
)

// Phase 标识更新过程的阶段。
type Phase string

const (
	PhaseIdle        Phase = ""
	PhaseChecking    Phase = "checking"
	PhaseDownloading Phase = "downloading"
	PhaseVerifying   Phase = "verifying"
	PhaseInstalling  Phase = "installing"
	PhaseDone        Phase = "done"
	PhaseError       Phase = "error"
)

// Release 描述一个候选的更新发布。
type Release struct {
	Version     string `json:"version"`
	Notes       string `json:"notes"`
	DownloadURL string `json:"downloadUrl"`
	Sha256URL   string `json:"sha256Url"`
	Size        int64  `json:"size"`
}

// Progress 是 updater 向订阅者推送的状态。
type Progress struct {
	Phase      Phase   `json:"phase"`
	Percent    float64 `json:"percent"`
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Error      string  `json:"error,omitempty"`
}

// Updater 是非线程安全的，全部入口方法通过 mutex 串行化。
type Updater struct {
	currentVersion string
	repo           string
	httpClient     *http.Client

	mu         sync.Mutex
	progressCh chan Progress
	cancel     context.CancelFunc
	inProgress bool
}

// New 构造一个 Updater。currentVersion 通常来自 main.Version（不含 v 前缀）。
func New(currentVersion string) *Updater {
	return &Updater{
		currentVersion: strings.TrimPrefix(currentVersion, "v"),
		repo:           defaultRepo,
		httpClient:     &http.Client{Timeout: httpTimeout},
		progressCh:     make(chan Progress, 32),
	}
}

// Subscribe 返回进度通道。多次调用返回同一个 channel。
func (u *Updater) Subscribe() <-chan Progress {
	return u.progressCh
}

// Cancel 中止当前 Download/Install 流程。
func (u *Updater) Cancel() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cancel != nil {
		u.cancel()
	}
}

// Check 查询 GitHub 最新 release。返回 nil 表示已是最新。
func (u *Updater) Check(ctx context.Context) (*Release, error) {
	u.emit(Progress{Phase: PhaseChecking})

	if u.currentVersion == "" || u.currentVersion == "dev" {
		return nil, nil
	}

	url := fmt.Sprintf("%s/repos/%s/releases/latest", apiBase, u.repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("查询 GitHub Releases 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub Releases 返回状态 %d", resp.StatusCode)
	}

	var ghRelease struct {
		TagName    string `json:"tag_name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Assets     []struct {
			Name               string `json:"name"`
			Size               int64  `json:"size"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghRelease); err != nil {
		return nil, fmt.Errorf("解析 release JSON 失败: %w", err)
	}
	if ghRelease.Draft || ghRelease.Prerelease {
		return nil, nil
	}

	latest := strings.TrimPrefix(ghRelease.TagName, "v")
	cmp, err := compareSemver(latest, u.currentVersion)
	if err != nil {
		return nil, fmt.Errorf("版本号格式错误: %w", err)
	}
	if cmp <= 0 {
		return nil, nil
	}

	assetName := PlatformAssetName(latest)
	var downloadURL, sha256URL string
	var size int64
	for _, a := range ghRelease.Assets {
		switch a.Name {
		case assetName:
			downloadURL = a.BrowserDownloadURL
			size = a.Size
		case assetName + ".sha256":
			sha256URL = a.BrowserDownloadURL
		}
	}
	if downloadURL == "" {
		return nil, fmt.Errorf("release %s 缺少平台资产 %s", latest, assetName)
	}

	return &Release{
		Version:     latest,
		Notes:       ghRelease.Body,
		DownloadURL: downloadURL,
		Sha256URL:   sha256URL,
		Size:        size,
	}, nil
}

// Download 把 release 资产下载到临时目录，返回本地文件路径。
func (u *Updater) Download(ctx context.Context, r *Release) (string, error) {
	if r == nil || r.DownloadURL == "" {
		return "", errors.New("无效的 release")
	}

	u.mu.Lock()
	if u.inProgress {
		u.mu.Unlock()
		return "", errors.New("已有更新任务在执行")
	}
	subCtx, cancel := context.WithCancel(ctx)
	u.cancel = cancel
	u.inProgress = true
	u.mu.Unlock()
	defer func() {
		u.mu.Lock()
		u.inProgress = false
		u.cancel = nil
		u.mu.Unlock()
	}()

	u.emit(Progress{Phase: PhaseDownloading, Total: r.Size})

	req, err := http.NewRequestWithContext(subCtx, http.MethodGet, r.DownloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败 HTTP %d", resp.StatusCode)
	}

	total := r.Size
	if total <= 0 {
		total = resp.ContentLength
	}

	assetName := PlatformAssetName(r.Version)
	tmpDir := filepath.Join(os.TempDir(), "ccx-updater")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(tmpDir, assetName)
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, downloadBuffer)
	var downloaded int64
	lastEmit := time.Now()
	for {
		select {
		case <-subCtx.Done():
			return "", subCtx.Err()
		default:
		}
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return "", werr
			}
			downloaded += int64(n)
			if time.Since(lastEmit) > 200*time.Millisecond {
				lastEmit = time.Now()
				percent := 0.0
				if total > 0 {
					percent = float64(downloaded) / float64(total) * 100
				}
				u.emit(Progress{
					Phase:      PhaseDownloading,
					Percent:    percent,
					Downloaded: downloaded,
					Total:      total,
				})
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", readErr
		}
	}

	u.emit(Progress{
		Phase:      PhaseDownloading,
		Percent:    100,
		Downloaded: downloaded,
		Total:      total,
	})

	return dest, nil
}

// Verify 校验文件 SHA256。sha256URL 为空时跳过（仅记录警告）。
func (u *Updater) Verify(ctx context.Context, localPath, sha256URL string) error {
	u.emit(Progress{Phase: PhaseVerifying})
	if sha256URL == "" {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sha256URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("拉取 sha256 失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("拉取 sha256 HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	expected := strings.ToLower(strings.TrimSpace(strings.Split(string(raw), " ")[0]))
	if expected == "" {
		return errors.New("sha256 内容为空")
	}

	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if got != expected {
		return fmt.Errorf("sha256 不匹配: 期望 %s, 实际 %s", expected, got)
	}
	return nil
}

// Install 调用平台特定的安装器。具体实现见 updater_{darwin,windows,linux}.go。
func (u *Updater) Install(localPath string) error {
	u.emit(Progress{Phase: PhaseInstalling})
	if err := platformInstall(localPath); err != nil {
		u.emit(Progress{Phase: PhaseError, Error: err.Error()})
		return err
	}
	u.emit(Progress{Phase: PhaseDone, Percent: 100})
	return nil
}

// CurrentVersion 返回 Updater 携带的当前版本。
func (u *Updater) CurrentVersion() string {
	return u.currentVersion
}

func (u *Updater) emit(p Progress) {
	select {
	case u.progressCh <- p:
	default:
		// drop oldest，避免订阅者迟钝时阻塞 updater
		select {
		case <-u.progressCh:
		default:
		}
		select {
		case u.progressCh <- p:
		default:
		}
	}
}

// PlatformAssetName 返回当前平台对应的资产文件名。
//
// 与 CI 命名约定必须保持一致：
//
//	CCX-Desktop-{version}-{os}-{arch}.{ext}
func PlatformAssetName(version string) string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("CCX-Desktop-%s-darwin-%s.dmg", version, runtime.GOARCH)
	case "windows":
		return fmt.Sprintf("CCX-Desktop-%s-windows-%s-setup.exe", version, runtime.GOARCH)
	case "linux":
		return fmt.Sprintf("CCX-Desktop-%s-linux-%s.AppImage", version, runtime.GOARCH)
	default:
		return fmt.Sprintf("CCX-Desktop-%s-%s-%s.bin", version, runtime.GOOS, runtime.GOARCH)
	}
}

// compareSemver 简易 semver 比较：返回 -1/0/1。
// 接受 "1.2.3" 或 "1.2.3-rc1"；预发布段按字典序比较。
func compareSemver(a, b string) (int, error) {
	parse := func(s string) ([]int, string, error) {
		pre := ""
		if idx := strings.IndexAny(s, "-+"); idx >= 0 {
			pre = s[idx:]
			s = s[:idx]
		}
		parts := strings.Split(s, ".")
		nums := make([]int, len(parts))
		for i, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, "", fmt.Errorf("非法 semver 段 %q", p)
			}
			nums[i] = n
		}
		return nums, pre, nil
	}
	an, ap, err := parse(a)
	if err != nil {
		return 0, err
	}
	bn, bp, err := parse(b)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(an) || i < len(bn); i++ {
		var av, bv int
		if i < len(an) {
			av = an[i]
		}
		if i < len(bn) {
			bv = bn[i]
		}
		if av != bv {
			if av > bv {
				return 1, nil
			}
			return -1, nil
		}
	}
	// numeric 相等：无预发布 > 有预发布
	switch {
	case ap == "" && bp == "":
		return 0, nil
	case ap == "" && bp != "":
		return 1, nil
	case ap != "" && bp == "":
		return -1, nil
	}
	if ap == bp {
		return 0, nil
	}
	if ap > bp {
		return 1, nil
	}
	return -1, nil
}
