// Package windowstate 负责窗口位置/尺寸/最大化状态的持久化。
//
// 数据存放在 {dataDir}/window-state.json。
package windowstate

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// State 描述一次窗口快照。
type State struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Maximised bool `json:"maximised"`
}

const fileName = "window-state.json"

const (
	minReasonableWidth  = 480
	minReasonableHeight = 320
	maxReasonableEdge   = 16384
)

// Load 从 dataDir 读取持久化状态。文件不存在或解析失败时返回 (zero, false, nil)。
func Load(dataDir string) (State, bool, error) {
	path := filepath.Join(dataDir, fileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return State{}, false, nil
		}
		return State{}, false, err
	}
	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		// 损坏的文件不应阻塞启动
		return State{}, false, nil
	}
	if !IsValid(s) {
		return State{}, false, nil
	}
	return s, true, nil
}

// Save 把状态原子写入 dataDir。
func Save(dataDir string, s State) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dataDir, "window-state-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName) // 成功时 rename 后已不在；失败时清理
	}()
	if err := json.NewEncoder(tmp).Encode(s); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	target := filepath.Join(dataDir, fileName)
	return os.Rename(tmpName, target)
}

// IsValid 做基本合理性校验，过滤掉异常值。
func IsValid(s State) bool {
	if s.Width < minReasonableWidth || s.Height < minReasonableHeight {
		return false
	}
	if s.Width > maxReasonableEdge || s.Height > maxReasonableEdge {
		return false
	}
	// 允许位置为 0,0 — 第一次显示在屏幕左上角是合法状态
	if s.X < -maxReasonableEdge || s.Y < -maxReasonableEdge {
		return false
	}
	if s.X > maxReasonableEdge || s.Y > maxReasonableEdge {
		return false
	}
	return true
}
