package config

import (
	"github.com/fatih/color"
)

func init() {
	// 强制启用颜色输出，即使 stdout 不是 TTY
	color.NoColor = false
}

// SegmentID 段类型枚举
type SegmentID string

const (
	SegmentModel         SegmentID = "model"
	SegmentDirectory     SegmentID = "directory"
	SegmentGit           SegmentID = "git"
	SegmentContextWindow SegmentID = "context_window"
	SegmentUsage         SegmentID = "usage"
	SegmentCost          SegmentID = "cost"
	SegmentSession       SegmentID = "session"
	SegmentOutputStyle   SegmentID = "output_style"
	SegmentUpdate        SegmentID = "update"
)

// ThemeMode 主题模式
type ThemeMode string

const (
	ThemeModeDefault  ThemeMode = "default"
	ThemeModeNerdFont ThemeMode = "nerd_font"
)

// SegmentTheme 单个 Segment 的主题配置
type SegmentTheme struct {
	Icon      string
	IconColor *color.Color
	TextColor *color.Color
	BgColor   *color.Color
	Bold      bool
}

// 预定义颜色 (ANSI 16 色)
var (
	colorLightCyan    = color.New(color.FgHiCyan)    // 14
	colorLightYellow  = color.New(color.FgHiYellow)  // 11
	colorLightGreen   = color.New(color.FgHiGreen)   // 10
	colorLightBlue    = color.New(color.FgHiBlue)    // 12
	colorLightMagenta = color.New(color.FgHiMagenta) // 13
	colorYellow       = color.New(color.FgYellow)    // 3
	colorGreen        = color.New(color.FgGreen)     // 2
	colorCyan         = color.New(color.FgCyan)      // 6
)

// Default 主题图标 (Emoji)
const (
	DefaultIconModel       = "🤖"
	DefaultIconDirectory   = "📁"
	DefaultIconGit         = "🌿"
	DefaultIconContext     = "⚡️"
	DefaultIconUsage       = "📊"
	DefaultIconCost        = "💰"
	DefaultIconSession     = "⏱️"
	DefaultIconOutputStyle = "🎯"
	DefaultIconUpdate      = "🔄"
)

// Nerd Font 主题图标 (Unicode 码点)
const (
	NerdFontIconModel       = "\ue26d"     // nf-md-creation
	NerdFontIconDirectory   = "\U000F024B" // nf-md-folder
	NerdFontIconGit         = "\U000F02A2" // nf-md-git
	NerdFontIconContext     = "\uf49b"     // nf-md-layers_triple
	NerdFontIconUsage       = "\U000F0A9E" // nf-md-chart_bar
	NerdFontIconCost        = "\ueec1"     // nf-md-currency_usd
	NerdFontIconSession     = "\U000F19BB" // nf-md-clock_outline
	NerdFontIconOutputStyle = "\U000F12F5" // nf-md-flag_variant
	NerdFontIconUpdate      = "\uf021"     // nf-fa-refresh
)

// segmentColors 定义各 Segment 的颜色配置 (与图标无关)
var segmentColors = map[SegmentID]struct {
	IconColor *color.Color
	TextColor *color.Color
	Bold      bool
}{
	SegmentModel:         {colorLightCyan, colorLightCyan, true},
	SegmentDirectory:     {colorLightYellow, colorLightGreen, true},
	SegmentGit:           {colorLightBlue, colorLightBlue, true},
	SegmentContextWindow: {colorLightMagenta, colorLightMagenta, true},
	SegmentUsage:         {colorLightCyan, colorLightCyan, false},
	SegmentCost:          {colorYellow, colorYellow, true},
	SegmentSession:       {colorGreen, colorGreen, true},
	SegmentOutputStyle:   {colorCyan, colorCyan, true},
	SegmentUpdate:        {colorLightYellow, colorLightYellow, false},
}

// defaultIcons Default 主题图标映射
var defaultIcons = map[SegmentID]string{
	SegmentModel:         DefaultIconModel,
	SegmentDirectory:     DefaultIconDirectory,
	SegmentGit:           DefaultIconGit,
	SegmentContextWindow: DefaultIconContext,
	SegmentUsage:         DefaultIconUsage,
	SegmentCost:          DefaultIconCost,
	SegmentSession:       DefaultIconSession,
	SegmentOutputStyle:   DefaultIconOutputStyle,
	SegmentUpdate:        DefaultIconUpdate,
}

// nerdFontIcons Nerd Font 主题图标映射
var nerdFontIcons = map[SegmentID]string{
	SegmentModel:         NerdFontIconModel,
	SegmentDirectory:     NerdFontIconDirectory,
	SegmentGit:           NerdFontIconGit,
	SegmentContextWindow: NerdFontIconContext,
	SegmentUsage:         NerdFontIconUsage,
	SegmentCost:          NerdFontIconCost,
	SegmentSession:       NerdFontIconSession,
	SegmentOutputStyle:   NerdFontIconOutputStyle,
	SegmentUpdate:        NerdFontIconUpdate,
}

// GetSegmentTheme 根据主题模式获取 Segment 主题配置
func GetSegmentTheme(id SegmentID, mode ThemeMode) SegmentTheme {
	// 获取颜色配置
	colors, ok := segmentColors[id]
	if !ok {
		return SegmentTheme{}
	}

	// 根据模式选择图标
	var icon string
	switch mode {
	case ThemeModeNerdFont:
		icon = nerdFontIcons[id]
	default:
		icon = defaultIcons[id]
	}

	return SegmentTheme{
		Icon:      icon,
		IconColor: colors.IconColor,
		TextColor: colors.TextColor,
		Bold:      colors.Bold,
	}
}

// ApplyColor applies foreground color to text
func ApplyColor(text string, c *color.Color) string {
	if c == nil {
		return text
	}
	return c.Sprint(text)
}

// ApplyBackground applies background color to text
func ApplyBackground(text string, c *color.Color) string {
	if c == nil {
		return text
	}
	return c.Sprint(text)
}
