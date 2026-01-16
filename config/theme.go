package config

import (
	"os"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	statuslineRendererOnce sync.Once
	statuslineRenderer     *lipgloss.Renderer
)

func getStatuslineRenderer() *lipgloss.Renderer {
	statuslineRendererOnce.Do(func() {
		// 强制启用颜色输出，即使 stdout 不是 TTY
		statuslineRenderer = lipgloss.NewRenderer(
			os.Stdout,
			termenv.WithUnsafe(),
			termenv.WithProfile(termenv.ANSI),
		)
	})
	return statuslineRenderer
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
	// CCH Segments
	SegmentCCHModel    SegmentID = "cch_model"
	SegmentCCHProvider SegmentID = "cch_provider"
	SegmentCCHCost     SegmentID = "cch_cost"
	SegmentCCHRequests SegmentID = "cch_requests"
	SegmentCCHLimits   SegmentID = "cch_limits"
	// Line Break
	SegmentLineBreak SegmentID = "---"
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
	IconColor *lipgloss.Color
	TextColor *lipgloss.Color
	BgColor   *lipgloss.Color
	Bold      bool
}

// 预定义颜色 (ANSI 16 色)
var (
	colorLightCyan    = lipgloss.Color("14") // 14
	colorLightYellow  = lipgloss.Color("11") // 11
	colorLightGreen   = lipgloss.Color("10") // 10
	colorLightBlue    = lipgloss.Color("12") // 12
	colorLightMagenta = lipgloss.Color("13") // 13
	colorYellow       = lipgloss.Color("3")  // 3
	colorGreen        = lipgloss.Color("2")  // 2
	colorCyan         = lipgloss.Color("6")  // 6
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
	// CCH Icons
	DefaultIconCCHModel    = "🔮"
	DefaultIconCCHProvider = "🏢"
	DefaultIconCCHCost     = "💵"
	DefaultIconCCHRequests = "📈"
	DefaultIconCCHLimits   = "🚦"
)

// Nerd Font 主题图标 (Unicode 码点)
const (
	NerdFontIconModel       = "\ue26d"     // nf-md-creation
	NerdFontIconDirectory   = "\U000F024B" // nf-md-folder
	NerdFontIconGit         = "\U000F02A2" // nf-md-git
	NerdFontIconContext     = "\uf49b"     // nf-md-layers_triple
	NerdFontIconUsage       = "\U000F0A9E" // nf-md-chart_bar
	NerdFontIconCost        = "\uf155"     // nf-md-currency_usd
	NerdFontIconSession     = "\U000F19BB" // nf-md-clock_outline
	NerdFontIconOutputStyle = "\U000F12F5" // nf-md-flag_variant
	NerdFontIconUpdate      = "\uf021"     // nf-fa-refresh
	// CCH Icons
	NerdFontIconCCHModel    = "\U000F02A1" // nf-md-ghost
	NerdFontIconCCHProvider = "\U000F0F74" // nf-md-server
	NerdFontIconCCHCost     = "\U000F01E0" // nf-md-cash
	NerdFontIconCCHRequests = "\U000F0127" // nf-md-chart_line
	NerdFontIconCCHLimits   = "\U000F0A1B" // nf-md-gauge
)

// segmentColors 定义各 Segment 的颜色配置 (与图标无关)
var segmentColors = map[SegmentID]struct {
	IconColor *lipgloss.Color
	TextColor *lipgloss.Color
	Bold      bool
}{
	SegmentModel:         {&colorLightCyan, &colorLightCyan, true},
	SegmentDirectory:     {&colorLightYellow, &colorLightGreen, true},
	SegmentGit:           {&colorLightBlue, &colorLightBlue, true},
	SegmentContextWindow: {&colorLightMagenta, &colorLightMagenta, true},
	SegmentUsage:         {&colorLightCyan, &colorLightCyan, false},
	SegmentCost:          {&colorYellow, &colorYellow, true},
	SegmentSession:       {&colorGreen, &colorGreen, true},
	SegmentOutputStyle:   {&colorCyan, &colorCyan, true},
	SegmentUpdate:        {&colorLightYellow, &colorLightYellow, false},
	// CCH Segments
	SegmentCCHModel:    {&colorLightMagenta, &colorLightMagenta, true},
	SegmentCCHProvider: {&colorLightBlue, &colorLightBlue, true},
	SegmentCCHCost:     {&colorYellow, &colorYellow, true},
	SegmentCCHRequests: {&colorLightGreen, &colorLightGreen, false},
	SegmentCCHLimits:   {&colorLightCyan, &colorLightCyan, false},
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
	// CCH Segments
	SegmentCCHModel:    DefaultIconCCHModel,
	SegmentCCHProvider: DefaultIconCCHProvider,
	SegmentCCHCost:     DefaultIconCCHCost,
	SegmentCCHRequests: DefaultIconCCHRequests,
	SegmentCCHLimits:   DefaultIconCCHLimits,
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
	// CCH Segments
	SegmentCCHModel:    NerdFontIconCCHModel,
	SegmentCCHProvider: NerdFontIconCCHProvider,
	SegmentCCHCost:     NerdFontIconCCHCost,
	SegmentCCHRequests: NerdFontIconCCHRequests,
	SegmentCCHLimits:   NerdFontIconCCHLimits,
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
func ApplyColor(text string, c *lipgloss.Color) string {
	if c == nil {
		return text
	}
	return getStatuslineRenderer().NewStyle().Foreground(*c).Render(text)
}

// ApplyBackground applies background color to text
func ApplyBackground(text string, c *lipgloss.Color) string {
	if c == nil {
		return text
	}
	return getStatuslineRenderer().NewStyle().Background(*c).Render(text)
}
