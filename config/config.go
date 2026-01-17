package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LineBreakMarker 换行分隔符标记
const LineBreakMarker = "---"

// DefaultSegmentOrder 定义默认的 segment 显示顺序
var DefaultSegmentOrder = []string{
	"model",
	"directory",
	"output_style",
	LineBreakMarker,
	"context_window",
}

func defaultEnabledForSegment(name string) bool {
	switch name {
	case "model", "directory", "output_style", "context_window":
		return true
	default:
		return false
	}
}

// DefaultSegmentEnabledForOrder returns the default per-instance enabled flags for an order.
// It excludes line breaks and enables only: model, directory, output_style, context_window.
func DefaultSegmentEnabledForOrder(order []string) []bool {
	enabled := make([]bool, 0, NonBreakSegmentCount(order))
	for _, name := range order {
		if name == LineBreakMarker {
			continue
		}
		enabled = append(enabled, defaultEnabledForSegment(name))
	}
	return enabled
}

// NormalizeSegmentEnabled ensures enabled length matches NonBreakSegmentCount(order).
// If mismatch, it returns DefaultSegmentEnabledForOrder(order).
func NormalizeSegmentEnabled(order []string, enabled []bool) []bool {
	expected := NonBreakSegmentCount(order)
	if len(enabled) == expected {
		return enabled
	}
	return DefaultSegmentEnabledForOrder(order)
}

// InputData is the JSON structure passed from Claude Code via stdin
type InputData struct {
	Model          ModelInfo       `json:"model"`
	Workspace      WorkspaceInfo   `json:"workspace"`
	TranscriptPath string          `json:"transcript_path"`
	Cost           CostInfo        `json:"cost"`
	OutputStyle    OutputStyleInfo `json:"output_style"`
}

// ModelInfo contains model identification information
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// WorkspaceInfo contains workspace directory information
type WorkspaceInfo struct {
	CurrentDir string `json:"current_dir"`
}

// CostInfo contains API cost information
type CostInfo struct {
	TotalCostUSD float64 `json:"total_cost_usd"`
}

// OutputStyleInfo contains output style information
type OutputStyleInfo struct {
	Name string `json:"name"`
}

// SimpleConfig is the TOML configuration structure
type SimpleConfig struct {
	Theme          ThemeMode `toml:"theme"` // "default" or "nerd_font"
	Separator      string    `toml:"separator"`
	SegmentOrder   []string  `toml:"segment_order"`
	SegmentEnabled []bool    `toml:"segment_enabled"`
	// CCH Configuration
	CCHApiKey string `toml:"cch_api_key"`
	CCHURL    string `toml:"cch_url"`
}

// SegmentToggles contains enable/disable flags for each segment
type SegmentToggles struct {
	Model         bool `toml:"model"`
	Directory     bool `toml:"directory"`
	Git           bool `toml:"git"`
	ContextWindow bool `toml:"context_window"`
	Usage         bool `toml:"usage"`
	Cost          bool `toml:"cost"`
	Session       bool `toml:"session"`
	OutputStyle   bool `toml:"output_style"`
	Update        bool `toml:"update"`
	// CCH Segments
	CCHModel    bool `toml:"cch_model"`
	CCHProvider bool `toml:"cch_provider"`
	CCHCost     bool `toml:"cch_cost"`
	CCHRequests bool `toml:"cch_requests"`
	CCHLimits   bool `toml:"cch_limits"`
}

// NonBreakSegmentCount returns number of segments excluding line breaks.
func NonBreakSegmentCount(order []string) int {
	count := 0
	for _, s := range order {
		if s == LineBreakMarker {
			continue
		}
		count++
	}
	return count
}

// SegmentTogglePtr returns pointer to toggle field for a segment name.
func SegmentTogglePtr(toggles *SegmentToggles, name string) *bool {
	switch name {
	case "model":
		return &toggles.Model
	case "directory":
		return &toggles.Directory
	case "git":
		return &toggles.Git
	case "context_window":
		return &toggles.ContextWindow
	case "usage":
		return &toggles.Usage
	case "cost":
		return &toggles.Cost
	case "session":
		return &toggles.Session
	case "output_style":
		return &toggles.OutputStyle
	case "update":
		return &toggles.Update
	case "cch_model":
		return &toggles.CCHModel
	case "cch_provider":
		return &toggles.CCHProvider
	case "cch_cost":
		return &toggles.CCHCost
	case "cch_requests":
		return &toggles.CCHRequests
	case "cch_limits":
		return &toggles.CCHLimits
	default:
		return nil
	}
}

// BuildSegmentEnabledFromToggles expands global toggles into per-instance enabled flags.
func BuildSegmentEnabledFromToggles(order []string, toggles SegmentToggles) []bool {
	enabled := make([]bool, 0, NonBreakSegmentCount(order))
	for _, name := range order {
		if name == LineBreakMarker {
			continue
		}
		if p := SegmentTogglePtr(&toggles, name); p != nil {
			enabled = append(enabled, *p)
		} else {
			enabled = append(enabled, false)
		}
	}
	return enabled
}

// LoadConfig loads configuration from ~/.claude/cchline/config.toml
// Returns default configuration if file doesn't exist
func LoadConfig() (*SimpleConfig, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".claude", "cchline", "config.toml")

	// Default configuration
	config := &SimpleConfig{
		Theme:        ThemeModeNerdFont,
		Separator:    " | ",
		SegmentOrder: DefaultSegmentOrder,
	}
	config.SegmentEnabled = DefaultSegmentEnabledForOrder(config.SegmentOrder)

	data, err := os.ReadFile(configPath)
	if err != nil {
		// File doesn't exist, return default config
		return config, nil
	}

	type diskConfig struct {
		Theme          ThemeMode      `toml:"theme"`
		Separator      string         `toml:"separator"`
		SegmentOrder   []string       `toml:"segment_order"`
		SegmentEnabled []bool         `toml:"segment_enabled"`
		Segments       SegmentToggles `toml:"segments"` // legacy
		CCHApiKey      string         `toml:"cch_api_key"`
		CCHURL         string         `toml:"cch_url"`
	}

	var disk diskConfig
	meta, err := toml.Decode(string(data), &disk)
	if err != nil {
		// Config format incompatible, use default config
		return config, nil
	}

	// 确保有默认值
	if disk.Theme != "" {
		config.Theme = disk.Theme
	}
	if disk.Separator != "" {
		config.Separator = disk.Separator
	}
	// 如果 SegmentOrder 为空，使用默认顺序
	if len(disk.SegmentOrder) > 0 {
		config.SegmentOrder = disk.SegmentOrder
	} else if len(config.SegmentOrder) == 0 {
		config.SegmentOrder = make([]string, len(DefaultSegmentOrder))
		copy(config.SegmentOrder, DefaultSegmentOrder)
	}

	config.CCHApiKey = disk.CCHApiKey
	config.CCHURL = disk.CCHURL

	expectedEnabled := NonBreakSegmentCount(config.SegmentOrder)
	if len(disk.SegmentEnabled) == expectedEnabled {
		config.SegmentEnabled = disk.SegmentEnabled
	} else if meta.IsDefined("segments") {
		config.SegmentEnabled = BuildSegmentEnabledFromToggles(config.SegmentOrder, disk.Segments)
	} else {
		config.SegmentEnabled = DefaultSegmentEnabledForOrder(config.SegmentOrder)
	}

	return config, nil
}
