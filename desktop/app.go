package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BenedictKing/ccx/desktop/internal/backend"
	"github.com/BenedictKing/ccx/desktop/internal/configservice"
	"github.com/BenedictKing/ccx/desktop/internal/updater"
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

type DesktopService struct {
	manager       *backend.Manager
	configService *configservice.Service
	app           *application.App
	mainWindow    application.Window
	updater       *updater.Updater
	versionInfo   VersionInfo
	notifications *notifications.NotificationService
}

type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	GitCommit string `json:"gitCommit"`
}

type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	Notes          string `json:"notes"`
	DownloadURL    string `json:"downloadUrl"`
	Sha256URL      string `json:"sha256Url"`
	Size           int64  `json:"size"`
}

type EnvFileState struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Exists  bool   `json:"exists"`
}

func NewDesktopService(manager *backend.Manager) *DesktopService {
	configService, _ := configservice.New(manager.DataDir())
	return &DesktopService{manager: manager, configService: configService}
}

func (s *DesktopService) setApp(app *application.App) {
	s.app = app
}

func (s *DesktopService) setMainWindow(window application.Window) {
	s.mainWindow = window
}

func (s *DesktopService) setVersion(v VersionInfo) {
	s.versionInfo = v
	s.updater = updater.New(v.Version)
}

func (s *DesktopService) setNotifications(svc *notifications.NotificationService) {
	s.notifications = svc
}

// Notify 推送一条系统通知。失败仅记录日志，不向上抛错。
//
//wails:ignore
func (s *DesktopService) Notify(title, body string) {
	if s.notifications == nil {
		return
	}
	id := fmt.Sprintf("ccx-%d", time.Now().UnixNano())
	err := s.notifications.SendNotification(notifications.NotificationOptions{
		ID:    id,
		Title: title,
		Body:  body,
	})
	if err != nil {
		log.Printf("[Desktop-Notify] 推送通知失败: %v", err)
	}
}

// CopyText 把文本写入系统剪贴板。
func (s *DesktopService) CopyText(text string) error {
	if s.app == nil {
		return fmt.Errorf("应用未初始化")
	}
	if !s.app.Clipboard.SetText(text) {
		return fmt.Errorf("写入剪贴板失败")
	}
	return nil
}

// WebURL 返回当前网关 Web UI 的访问地址（即使服务未启动，也基于配置端口拼接）。
func (s *DesktopService) WebURL() string {
	return s.manager.WebURL()
}

// GetVersion 返回构建时注入的版本信息。
func (s *DesktopService) GetVersion() VersionInfo {
	return s.versionInfo
}

func (s *DesktopService) GetStatus() backend.Status {
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	return s.manager.Status(ctx)
}

func (s *DesktopService) GetProxyAccessKey() (string, error) {
	return s.manager.EnsureProxyAccessKey()
}

func (s *DesktopService) GetEnvFile() (EnvFileState, error) {
	path := filepath.Join(s.manager.DataDir(), ".env")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return EnvFileState{Path: path, Exists: false}, nil
		}
		return EnvFileState{}, err
	}
	return EnvFileState{Path: path, Content: string(content), Exists: true}, nil
}

func (s *DesktopService) SaveEnvFile(content string) error {
	path := filepath.Join(s.manager.DataDir(), ".env")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o600)
}

func (s *DesktopService) StartService() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return s.manager.Start(ctx)
}

func (s *DesktopService) StopService() error {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	return s.manager.Stop(ctx)
}

func (s *DesktopService) RestartService() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.manager.Restart(ctx)
}

func (s *DesktopService) GetLogs() []string {
	return s.manager.Logs()
}

func (s *DesktopService) GetAgentConfigStatus(platform string) (configservice.AgentConfigStatus, error) {
	if s.configService == nil {
		return configservice.AgentConfigStatus{}, fmt.Errorf("配置服务未初始化")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	status := s.manager.Status(ctx)
	return s.configService.GetStatus(platform, status.Port)
}

func (s *DesktopService) ApplyAgentConfig(req configservice.ApplyAgentConfigRequest) error {
	if s.configService == nil {
		return fmt.Errorf("配置服务未初始化")
	}
	platform := req.Platform
	if platform == "" {
		return fmt.Errorf("agent 平台不能为空")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	status := s.manager.Status(ctx)
	var key string
	if platform == configservice.PlatformCodex || (platform == configservice.PlatformClaude && (req.Provider == "" || req.Provider == configservice.ProviderCCX)) {
		if !status.Running {
			return fmt.Errorf("请先启动 CCX 服务")
		}
		var err error
		key, err = s.manager.EnsureProxyAccessKey()
		if err != nil {
			return err
		}
	}
	return s.configService.Apply(req, status.Port, key)
}

func (s *DesktopService) RestoreAgentConfig(platform string) error {
	if s.configService == nil {
		return fmt.Errorf("配置服务未初始化")
	}
	return s.configService.Restore(platform)
}

func (s *DesktopService) ShowStatusTab() error {
	s.showWindow()
	if s.app != nil {
		s.app.Event.Emit("desktop:show-tab", "status")
	}
	return nil
}

func (s *DesktopService) ShowAgentTab() error {
	s.showWindow()
	if s.app != nil {
		s.app.Event.Emit("desktop:show-tab", "agent")
	}
	return nil
}

func (s *DesktopService) ShowWebUITab() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.manager.Start(ctx); err != nil {
		return err
	}
	if err := s.manager.WaitHealthy(ctx, 15*time.Second); err != nil {
		return err
	}
	s.showWindow()
	if s.app != nil {
		s.app.Event.Emit("desktop:show-tab", "web")
	}
	return nil
}

func (s *DesktopService) OpenWebUIInBrowser() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.manager.Start(ctx); err != nil {
		return err
	}
	if err := s.manager.WaitHealthy(ctx, 15*time.Second); err != nil {
		return err
	}
	return browser.OpenURL(s.manager.WebURL())
}

func (s *DesktopService) GetAutostartStatus() (bool, error) {
	if s.app == nil {
		return false, fmt.Errorf("应用未初始化")
	}
	return s.app.Autostart.IsEnabled()
}

func (s *DesktopService) SetAutostart(enabled bool) error {
	if s.app == nil {
		return fmt.Errorf("应用未初始化")
	}
	if enabled {
		return s.app.Autostart.Enable()
	}
	return s.app.Autostart.Disable()
}

func (s *DesktopService) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_ = s.manager.Stop(ctx)
}

// CheckUpdate 查询是否有新版本可用。
func (s *DesktopService) CheckUpdate() (UpdateInfo, error) {
	info := UpdateInfo{CurrentVersion: s.versionInfo.Version}
	if s.updater == nil {
		return info, fmt.Errorf("updater 未初始化")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	release, err := s.updater.Check(ctx)
	if err != nil {
		return info, err
	}
	if release == nil {
		return info, nil
	}
	info.Available = true
	info.LatestVersion = release.Version
	info.Notes = release.Notes
	info.DownloadURL = release.DownloadURL
	info.Sha256URL = release.Sha256URL
	info.Size = release.Size
	if s.app != nil {
		s.app.Event.Emit("update:available", info)
	}
	return info, nil
}

// DownloadAndInstall 下载、校验并触发安装。整个流程通过 update:progress 事件推送进度。
func (s *DesktopService) DownloadAndInstall(info UpdateInfo) error {
	if s.updater == nil {
		return fmt.Errorf("updater 未初始化")
	}
	if info.DownloadURL == "" {
		return fmt.Errorf("缺少下载地址")
	}
	release := &updater.Release{
		Version:     info.LatestVersion,
		DownloadURL: info.DownloadURL,
		Sha256URL:   info.Sha256URL,
		Size:        info.Size,
	}
	go s.runUpdate(release)
	return nil
}

func (s *DesktopService) runUpdate(release *updater.Release) {
	// 订阅进度并转发到前端
	if s.app != nil {
		go func() {
			for p := range s.updater.Subscribe() {
				s.app.Event.Emit("update:progress", p)
				if p.Phase == updater.PhaseDone || p.Phase == updater.PhaseError {
					return
				}
			}
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	localPath, err := s.updater.Download(ctx, release)
	if err != nil {
		if s.app != nil {
			s.app.Event.Emit("update:progress", updater.Progress{
				Phase: updater.PhaseError,
				Error: err.Error(),
			})
		}
		return
	}

	if err := s.updater.Verify(ctx, localPath, release.Sha256URL); err != nil {
		if s.app != nil {
			s.app.Event.Emit("update:progress", updater.Progress{
				Phase: updater.PhaseError,
				Error: err.Error(),
			})
		}
		return
	}

	if err := s.updater.Install(localPath); err != nil {
		if s.app != nil {
			s.app.Event.Emit("update:progress", updater.Progress{
				Phase: updater.PhaseError,
				Error: err.Error(),
			})
		}
	}
}

// CancelUpdate 中止正在进行的下载/安装。
func (s *DesktopService) CancelUpdate() error {
	if s.updater != nil {
		s.updater.Cancel()
	}
	return nil
}

func (s *DesktopService) showWindow() {
	if s.mainWindow == nil {
		return
	}
	if s.mainWindow.IsMinimised() {
		s.mainWindow.UnMinimise()
	}
	s.mainWindow.Show()
	s.mainWindow.Focus()
}
