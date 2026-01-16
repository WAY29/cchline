package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// DefaultSegmentOrder 定义默认的 segment 显示顺序
var DefaultSegmentOrder = []string{
	"model",
	"directory",
	"git",
	"context_window",
	"usage",
	"cost",
	"session",
	"output_style",
	"update",
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
	Theme        ThemeMode      `toml:"theme"`         // "default" or "nerd_font"
	Separator    string         `toml:"separator"`
	SegmentOrder []string       `toml:"segment_order"`
	Segments     SegmentToggles `toml:"segments"`
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
		Segments: SegmentToggles{
			Model:         true,
			Directory:     true,
			Git:           true,
			ContextWindow: true,
			Usage:         false,
			Cost:          false,
			Session:       false,
			OutputStyle:   false,
			Update:        false,
		},
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// File doesn't exist, return default config
		return config, nil
	}

	if err := toml.Unmarshal(data, config); err != nil {
		// Config format incompatible, use default config
		return config, nil
	}

	// 确保有默认值
	if config.Theme == "" {
		config.Theme = ThemeModeNerdFont
	}
	if config.Separator == "" {
		config.Separator = " | "
	}
	// 如果 SegmentOrder 为空，使用默认顺序
	if len(config.SegmentOrder) == 0 {
		config.SegmentOrder = make([]string, len(DefaultSegmentOrder))
		copy(config.SegmentOrder, DefaultSegmentOrder)
	}

	return config, nil
}
