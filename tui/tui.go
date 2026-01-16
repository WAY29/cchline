package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/WAY29/cchline/config"
	"github.com/charmbracelet/bubbles/textinput"
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
	isTextInput bool // 是否为文本输入项
	textKey     string // 文本输入的配置键名 ("cch_url" 或 "cch_api_key")
}

// Model TUI 模型
type Model struct {
	config       *config.SimpleConfig
	cursor       int
	items        []menuItem
	quitting     bool
	width        int
	height       int
	debugKey     string // 调试：显示最后按下的按键
	editing      bool // 是否处于编辑模式
	textInputs   map[string]textinput.Model // 文本输入组件
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

	// CCH SETTINGS 设置
	items = append(items, menuItem{label: "CCH SETTINGS", key: "", isHeader: true})
	items = append(items, menuItem{label: "CCH URL", key: "cch_url", isTextInput: true, textKey: "cch_url"})
	items = append(items, menuItem{label: "API Key", key: "cch_api_key", isTextInput: true, textKey: "cch_api_key"})

	// 初始化文本输入组件
	textInputs := make(map[string]textinput.Model)

	// CCH URL 输入
	cchURLInput := textinput.New()
	cchURLInput.Placeholder = "https://example.com"
	cchURLInput.SetValue(cfg.CCHURL)
	textInputs["cch_url"] = cchURLInput

	// API Key 输入（密码模式）
	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "Enter API Key"
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '*'
	apiKeyInput.SetValue(cfg.CCHApiKey)
	textInputs["cch_api_key"] = apiKeyInput

	// 找到第一个非 header 项
	cursor := 0
	for i, item := range items {
		if !item.isHeader {
			cursor = i
			break
		}
	}

	return Model{
		config:     cfg,
		cursor:     cursor,
		items:      items,
		width:      50,
		height:     20,
		textInputs: textInputs,
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

		// 如果处于编辑模式，处理文本输入
		if m.editing {
			switch msg.String() {
			case "enter":
				// 保存当前编辑的值
				m.saveTextInputValue()
				m.editing = false
				return m, nil
			case "esc":
				// 取消编辑，不保存
				m.editing = false
				return m, nil
			default:
				// 将按键传递给 textinput
				item := m.items[m.cursor]
				if item.isTextInput {
					input := m.textInputs[item.textKey]
					var cmd tea.Cmd
					input, cmd = input.Update(msg)
					m.textInputs[item.textKey] = input
					return m, cmd
				}
			}
		}

		// 非编辑模式的按键处理
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
			// 检查是否是文本输入项
			item := m.items[m.cursor]
			if item.isTextInput {
				m.editing = true
				// 聚焦到对应的 textinput
				input := m.textInputs[item.textKey]
				input.Focus()
				m.textInputs[item.textKey] = input
			} else {
				m.toggleItem()
			}
		}
	}

	return m, nil
}

// saveTextInputValue 保存文本输入的值到配置
func (m *Model) saveTextInputValue() {
	item := m.items[m.cursor]
	if !item.isTextInput {
		return
	}

	input := m.textInputs[item.textKey]
	value := input.Value()

	switch item.textKey {
	case "cch_url":
		m.config.CCHURL = value
	case "cch_api_key":
		m.config.CCHApiKey = value
	}
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
	// 只能移动 segment 项（非 header、非 separator、非 theme、非 textInput）
	if item.isHeader || item.isSeparator || item.key == "theme" || item.isTextInput {
		return
	}

	// 找到 SEGMENTS header 和 CCH SETTINGS header 的位置
	segmentsStart := -1
	segmentsEnd := -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEGMENTS" {
			segmentsStart = i + 1
		} else if it.isHeader && it.label == "CCH SETTINGS" {
			segmentsEnd = i
			break
		}
	}
	if segmentsStart < 0 {
		return
	}
	if segmentsEnd < 0 {
		segmentsEnd = len(m.items)
	}

	// 计算当前 segment 在 segment 列表中的相对位置
	segmentIndex := m.cursor - segmentsStart
	newIndex := segmentIndex + delta

	// 边界检查：只能在 SEGMENTS 区域内移动
	segmentCount := segmentsEnd - segmentsStart
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
	// 找到 SEGMENTS header 和 CCH SETTINGS header 的位置
	segmentsStart := -1
	segmentsEnd := -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEGMENTS" {
			segmentsStart = i + 1
		} else if it.isHeader && it.label == "CCH SETTINGS" {
			segmentsEnd = i
			break
		}
	}
	if segmentsStart < 0 {
		return
	}
	if segmentsEnd < 0 {
		segmentsEnd = len(m.items)
	}

	// 重建 SegmentOrder（只包含 SEGMENTS 区域内的项）
	var order []string
	for i := segmentsStart; i < segmentsEnd; i++ {
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
		if item.isTextInput {
			// 文本输入项显示
			input := m.textInputs[item.textKey]
			var displayValue string

			if m.editing && m.cursor == i {
				// 编辑模式：显示 textinput 的 View()
				displayValue = input.View()
			} else {
				// 非编辑模式：显示当前值
				value := input.Value()
				if item.textKey == "cch_api_key" && value != "" {
					// API Key 显示为 ****
					displayValue = valueStyle.Render("****")
				} else if value == "" {
					displayValue = disabledStyle.Render("(empty)")
				} else {
					displayValue = valueStyle.Render(value)
				}
			}

			if m.cursor == i {
				line = fmt.Sprintf("%s%s  %s", cursor, selectedStyle.Render(item.label), displayValue)
			} else {
				line = fmt.Sprintf("%s%s  %s", cursor, normalStyle.Render(item.label), displayValue)
			}
		} else if item.key == "theme" {
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
	var help string
	if m.editing {
		help = helpBarStyle.Render(
			keyStyle.Render("Enter") + " Save  " +
				keyStyle.Render("Esc") + " Cancel",
		)
	} else {
		help = helpBarStyle.Render(
			keyStyle.Render("↑↓") + " Navigate  " +
				keyStyle.Render("Space") + " Toggle  " +
				keyStyle.Render(ReorderKeyHint) + " Reorder  " +
				keyStyle.Render("Esc") + " Save & Exit",
		)
	}

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
