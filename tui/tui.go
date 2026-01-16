package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/WAY29/cchline/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 样式定义
var (
	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			MarginBottom(1)

	// 分组标题
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			MarginTop(1)

	// 选中项样式
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	// 普通项样式
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// 启用状态
	enabledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("78"))

	// 禁用状态
	disabledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// 值样式
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	// 帮助栏样式
	helpBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			MarginTop(1)

	// 快捷键样式
	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	// 边框样式
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)
)

// SeparatorPreset 分隔符预设
type SeparatorPreset struct {
	Label string
	Value string
}

// SeparatorPresets 预设分隔符列表
var SeparatorPresets = []SeparatorPreset{
	{Label: "Pipe", Value: " | "},
	{Label: "Dot", Value: " · "},
	{Label: "Bar", Value: " ⁞ "},
	{Label: "Arrow", Value: " → "},
	{Label: "Chevron", Value: " ❯ "},
}

// menuItem 菜单项
type menuItem struct {
	label       string
	key         string
	enabled     bool
	isHeader    bool
	isSeparator bool // 是否为分隔符选项
	isSelected  bool // 分隔符是否被选中
}

// Model TUI 模型
type Model struct {
	config   *config.SimpleConfig
	cursor   int
	items    []menuItem
	quitting bool
	width    int
	height   int
	debugKey string // 调试：显示最后按下的按键
}

// NewModel 创建新的 TUI 模型
func NewModel(cfg *config.SimpleConfig) Model {
	items := []menuItem{
		// Theme 设置
		{label: "THEME", key: "", isHeader: true},
		{label: "Theme", key: "theme", enabled: cfg.Theme == config.ThemeModeNerdFont},

		// Separator 设置
		{label: "SEPARATOR", key: "", isHeader: true},
	}

	// 添加分隔符预设选项
	for _, preset := range SeparatorPresets {
		items = append(items, menuItem{
			label:       preset.Label,
			key:         preset.Value,
			isSeparator: true,
			isSelected:  cfg.Separator == preset.Value,
		})
	}

	// Segments 设置
	items = append(items, menuItem{label: "SEGMENTS", key: "", isHeader: true})

	// segment 名称到标签的映射
	segmentLabels := map[string]string{
		"model":          "Model",
		"directory":      "Directory",
		"git":            "Git",
		"context_window": "Context Window",
		"usage":          "Usage",
		"cost":           "Cost",
		"session":        "Session",
		"output_style":   "Output Style",
		"update":         "Update",
		// CCH Segments
		"cch_model":    "CCH Model",
		"cch_provider": "CCH Provider",
		"cch_cost":     "CCH Cost",
		"cch_requests": "CCH Requests",
		"cch_limits":   "CCH Limits",
	}

	// segment 名称到启用状态的映射
	segmentEnabled := map[string]bool{
		"model":          cfg.Segments.Model,
		"directory":      cfg.Segments.Directory,
		"git":            cfg.Segments.Git,
		"context_window": cfg.Segments.ContextWindow,
		"usage":          cfg.Segments.Usage,
		"cost":           cfg.Segments.Cost,
		"session":        cfg.Segments.Session,
		"output_style":   cfg.Segments.OutputStyle,
		"update":         cfg.Segments.Update,
		// CCH Segments
		"cch_model":    cfg.Segments.CCHModel,
		"cch_provider": cfg.Segments.CCHProvider,
		"cch_cost":     cfg.Segments.CCHCost,
		"cch_requests": cfg.Segments.CCHRequests,
		"cch_limits":   cfg.Segments.CCHLimits,
	}

	// 按照配置的顺序添加 segments
	for _, name := range cfg.SegmentOrder {
		if label, exists := segmentLabels[name]; exists {
			items = append(items, menuItem{
				label:   label,
				key:     name,
				enabled: segmentEnabled[name],
			})
		}
	}

	// 找到第一个非 header 项
	cursor := 0
	for i, item := range items {
		if !item.isHeader {
			cursor = i
			break
		}
	}

	return Model{
		config: cfg,
		cursor: cursor,
		items:  items,
		width:  50,
		height: 20,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// 调试：记录按键
		m.debugKey = msg.String()

		switch msg.String() {
		case "ctrl+c", "esc":
			// 退出时自动保存
			m.saveConfig()
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.moveCursor(-1)

		case "down", "j":
			m.moveCursor(1)

		case MoveUpKey:
			m.moveSegment(-1)

		case MoveDownKey:
			m.moveSegment(1)

		case "enter", " ":
			m.toggleItem()
		}
	}

	return m, nil
}

// moveCursor 移动光标
func (m *Model) moveCursor(delta int) {
	newCursor := m.cursor + delta

	// 边界检查
	if newCursor < 0 || newCursor >= len(m.items) {
		return
	}

	// 跳过 header
	for newCursor >= 0 && newCursor < len(m.items) && m.items[newCursor].isHeader {
		newCursor += delta
	}

	if newCursor >= 0 && newCursor < len(m.items) {
		m.cursor = newCursor
	}
}

// moveSegment 移动 segment 顺序
func (m *Model) moveSegment(delta int) {
	item := m.items[m.cursor]
	// 只能移动 segment 项（非 header、非 separator、非 theme）
	if item.isHeader || item.isSeparator || item.key == "theme" {
		return
	}

	// 找到 SEGMENTS header 的位置
	segmentsStart := -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEGMENTS" {
			segmentsStart = i + 1
			break
		}
	}
	if segmentsStart < 0 {
		return
	}

	// 计算当前 segment 在 segment 列表中的相对位置
	segmentIndex := m.cursor - segmentsStart
	newIndex := segmentIndex + delta

	// 边界检查
	segmentCount := len(m.items) - segmentsStart
	if newIndex < 0 || newIndex >= segmentCount {
		return
	}

	// 交换 items
	targetCursor := m.cursor + delta
	m.items[m.cursor], m.items[targetCursor] = m.items[targetCursor], m.items[m.cursor]
	m.cursor = targetCursor

	// 更新配置中的 SegmentOrder
	m.updateSegmentOrder()
}

// updateSegmentOrder 更新配置中的 segment 顺序
func (m *Model) updateSegmentOrder() {
	// 找到 SEGMENTS header 的位置
	segmentsStart := -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEGMENTS" {
			segmentsStart = i + 1
			break
		}
	}
	if segmentsStart < 0 {
		return
	}

	// 重建 SegmentOrder
	var order []string
	for i := segmentsStart; i < len(m.items); i++ {
		if !m.items[i].isHeader {
			order = append(order, m.items[i].key)
		}
	}
	m.config.SegmentOrder = order
}

// toggleItem 切换选项
func (m *Model) toggleItem() {
	item := &m.items[m.cursor]
	if item.isHeader {
		return
	}

	// 处理分隔符选择（单选）
	if item.isSeparator {
		// 取消其他分隔符的选中状态
		for i := range m.items {
			if m.items[i].isSeparator {
				m.items[i].isSelected = false
			}
		}
		item.isSelected = true
		m.config.Separator = item.key
		return
	}

	item.enabled = !item.enabled

	// 更新配置
	switch item.key {
	case "theme":
		if item.enabled {
			m.config.Theme = config.ThemeModeNerdFont
		} else {
			m.config.Theme = config.ThemeModeDefault
		}
	case "model":
		m.config.Segments.Model = item.enabled
	case "directory":
		m.config.Segments.Directory = item.enabled
	case "git":
		m.config.Segments.Git = item.enabled
	case "context_window":
		m.config.Segments.ContextWindow = item.enabled
	case "usage":
		m.config.Segments.Usage = item.enabled
	case "cost":
		m.config.Segments.Cost = item.enabled
	case "session":
		m.config.Segments.Session = item.enabled
	case "output_style":
		m.config.Segments.OutputStyle = item.enabled
	case "update":
		m.config.Segments.Update = item.enabled
	// CCH Segments
	case "cch_model":
		m.config.Segments.CCHModel = item.enabled
	case "cch_provider":
		m.config.Segments.CCHProvider = item.enabled
	case "cch_cost":
		m.config.Segments.CCHCost = item.enabled
	case "cch_requests":
		m.config.Segments.CCHRequests = item.enabled
	case "cch_limits":
		m.config.Segments.CCHLimits = item.enabled
	}
}

// saveConfig 保存配置
func (m *Model) saveConfig() error {
	configDir := filepath.Join(os.Getenv("HOME"), ".claude", "cchline")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.toml")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(m.config)
}

// View 渲染视图
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// 标题
	title := titleStyle.Render(" CCHLine Configuration ")

	// 内容
	var content string
	for i, item := range m.items {
		if item.isHeader {
			// 分组标题
			content += sectionStyle.Render("  "+item.label) + "\n"
			continue
		}

		// 光标指示器
		cursor := "   "
		if m.cursor == i {
			cursor = " ▸ "
		}

		// 渲染行
		var line string
		if item.key == "theme" {
			// Theme 特殊显示
			var themeValue string
			if item.enabled {
				themeValue = valueStyle.Render("nerd_font")
			} else {
				themeValue = valueStyle.Render("default")
			}

			if m.cursor == i {
				line = fmt.Sprintf("%s%s  %s", cursor, selectedStyle.Render(item.label), themeValue)
			} else {
				line = fmt.Sprintf("%s%s  %s", cursor, normalStyle.Render(item.label), themeValue)
			}
		} else if item.isSeparator {
			// 分隔符选项显示
			var status string
			if item.isSelected {
				status = enabledStyle.Render("●")
			} else {
				status = disabledStyle.Render("○")
			}
			// 显示分隔符预览
			preview := valueStyle.Render(fmt.Sprintf("%q", item.key))

			if m.cursor == i {
				line = fmt.Sprintf("%s%s %s  %s", cursor, status, selectedStyle.Render(item.label), preview)
			} else {
				line = fmt.Sprintf("%s%s %s  %s", cursor, status, normalStyle.Render(item.label), preview)
			}
		} else {
			// Segment 开关显示
			var status string
			if item.enabled {
				status = enabledStyle.Render("●")
			} else {
				status = disabledStyle.Render("○")
			}

			if m.cursor == i {
				line = fmt.Sprintf("%s%s %s", cursor, status, selectedStyle.Render(item.label))
			} else {
				line = fmt.Sprintf("%s%s %s", cursor, status, normalStyle.Render(item.label))
			}
		}

		content += line + "\n"
	}

	// 包装内容
	box := boxStyle.Render(content)

	// 帮助栏
	help := helpBarStyle.Render(
		keyStyle.Render("↑↓") + " Navigate  " +
			keyStyle.Render("Space") + " Toggle  " +
			keyStyle.Render(ReorderKeyHint) + " Reorder  " +
			keyStyle.Render("Esc") + " Save & Exit",
	)

	// 调试信息
	debug := fmt.Sprintf("\n  DEBUG: Last key = %q", m.debugKey)

	return fmt.Sprintf("\n%s\n%s\n%s%s\n", title, box, help, debug)
}

// Run 运行 TUI
func Run(cfg *config.SimpleConfig) error {
	p := tea.NewProgram(NewModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
