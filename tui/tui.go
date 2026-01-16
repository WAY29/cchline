package tui

import (
	"fmt"
	"strings"

	"github.com/WAY29/cchline/config"
	"github.com/WAY29/cchline/render"
	"github.com/WAY29/cchline/segment"
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

	// 预览区域样式
	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(0, 2).
			MarginBottom(1)

	// 预览标签样式
	previewLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)
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
	isSeparator bool   // 是否为分隔符选项
	isSelected  bool   // 分隔符是否被选中
	isTextInput bool   // 是否为文本输入项
	textKey     string // 文本输入的配置键名 ("cch_url" 或 "cch_api_key")
	isLineBreak bool   // 是否为换行分隔符
}

// Model TUI 模型
type Model struct {
	config        *config.SimpleConfig
	cursor        int
	items         []menuItem
	quitting      bool
	width         int
	height        int
	debugKey      string                     // 调试：显示最后按下的按键
	editing       bool                       // 是否处于编辑模式
	textInputs    map[string]textinput.Model // 文本输入组件
	confirmAction string                     // 待确认的操作: "install" 或 "uninstall"
	statusMessage string                     // 操作结果消息
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
		// 处理换行分隔符
		if name == config.LineBreakMarker {
			items = append(items, menuItem{
				label:       "Line Break",
				key:         config.LineBreakMarker,
				isLineBreak: true,
			})
			continue
		}

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

// generatePreview 生成状态栏预览
func (m Model) generatePreview() string {
	// Mock 数据映射
	mockData := map[string]string{
		"model":          "Opus 4.5",
		"directory":      "myapp",
		"git":            "main *",
		"context_window": "15.6%",
		"usage":          "↓31K ↑5K",
		"cost":           "$0.15",
		"session":        "1h23m",
		"output_style":   "default",
		"update":         "",
		"cch_model":      "claude-3-opus",
		"cch_provider":   "anthropic",
		"cch_cost":       "$1.50/$10",
		"cch_requests":   "123 reqs",
		"cch_limits":     "5h:$0",
	}

	// 根据当前配置构建 SegmentResult 列表
	var results []segment.SegmentResult

	for _, name := range m.config.SegmentOrder {
		// 处理换行分隔符
		if name == config.LineBreakMarker {
			results = append(results, segment.SegmentResult{
				ID:   config.SegmentLineBreak,
				Data: segment.SegmentData{},
			})
			continue
		}

		// 检查 segment 是否启用
		if !m.isSegmentEnabled(name) {
			continue
		}

		// 获取 mock 数据
		mockValue, exists := mockData[name]
		if !exists || mockValue == "" {
			continue
		}

		// 构建 SegmentResult
		results = append(results, segment.SegmentResult{
			ID:   config.SegmentID(name),
			Data: segment.SegmentData{Primary: mockValue},
		})
	}

	// 使用 StatusLineGenerator 生成预览
	generator := render.NewStatusLineGenerator(m.config)
	return generator.Generate(results)
}

// isSegmentEnabled 检查 segment 是否启用
func (m Model) isSegmentEnabled(name string) bool {
	switch name {
	case "model":
		return m.config.Segments.Model
	case "directory":
		return m.config.Segments.Directory
	case "git":
		return m.config.Segments.Git
	case "context_window":
		return m.config.Segments.ContextWindow
	case "usage":
		return m.config.Segments.Usage
	case "cost":
		return m.config.Segments.Cost
	case "session":
		return m.config.Segments.Session
	case "output_style":
		return m.config.Segments.OutputStyle
	case "update":
		return m.config.Segments.Update
	case "cch_model":
		return m.config.Segments.CCHModel
	case "cch_provider":
		return m.config.Segments.CCHProvider
	case "cch_cost":
		return m.config.Segments.CCHCost
	case "cch_requests":
		return m.config.Segments.CCHRequests
	case "cch_limits":
		return m.config.Segments.CCHLimits
	default:
		return false
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
			case "ctrl+u":
				// 清空当前文本输入
				item := m.items[m.cursor]
				if item.isTextInput {
					input := m.textInputs[item.textKey]
					input.SetValue("")
					m.textInputs[item.textKey] = input
				}
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

		// 确认模式的按键处理
		if m.confirmAction != "" {
			switch msg.String() {
			case "y", "Y":
				// 执行确认的操作
				var err error
				if m.confirmAction == "install" {
					err = m.installStatusLine()
					if err == nil {
						m.statusMessage = "✓ 安装成功"
					} else {
						m.statusMessage = "✗ " + err.Error()
					}
				} else if m.confirmAction == "uninstall" {
					err = m.uninstallStatusLine()
					if err == nil {
						m.statusMessage = "✓ 卸载成功"
					} else {
						m.statusMessage = "✗ " + err.Error()
					}
				}
				m.confirmAction = ""
				return m, nil
			case "n", "N", "esc":
				// 取消操作
				m.confirmAction = ""
				m.statusMessage = ""
				return m, nil
			default:
				// 忽略其他按键
				return m, nil
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

		case "a":
			m.insertLineBreak()

		case "d":
			m.deleteLineBreak()

		case "i":
			// 安装确认
			m.confirmAction = "install"
			m.statusMessage = ""

		case "u":
			// 卸载确认
			m.confirmAction = "uninstall"
			m.statusMessage = ""

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

// View 渲染视图
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// 标题
	title := titleStyle.Render(" CCHLine Configuration ")

	// 预览区域
	previewLabel := previewLabelStyle.Render("  Preview:")
	previewContent := m.generatePreview()
	preview := previewStyle.Render(previewContent)

	// 内容
	var content strings.Builder
	for i, item := range m.items {
		if item.isHeader {
			// 分组标题
			content.WriteString(sectionStyle.Render("  "+item.label) + "\n")
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
			separatorPreview := valueStyle.Render(fmt.Sprintf("%q", item.key))

			if m.cursor == i {
				line = fmt.Sprintf("%s%s %s  %s", cursor, status, selectedStyle.Render(item.label), separatorPreview)
			} else {
				line = fmt.Sprintf("%s%s %s  %s", cursor, status, normalStyle.Render(item.label), separatorPreview)
			}
		} else if item.isLineBreak {
			// 换行分隔符显示
			if m.cursor == i {
				line = fmt.Sprintf("%s%s", cursor, selectedStyle.Render("↵ Line Break"))
			} else {
				line = fmt.Sprintf("%s%s", cursor, disabledStyle.Render("↵ Line Break"))
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

		content.WriteString(line + "\n")
	}

	// 包装内容
	box := boxStyle.Render(content.String())

	// 安装状态
	var installStatus string
	if isInstalled() {
		installStatus = "  " + enabledStyle.Render("● 已安装")
	} else {
		installStatus = "  " + disabledStyle.Render("○ 未安装")
	}

	// 帮助栏
	var help string
	if m.confirmAction != "" {
		// 确认模式
		var actionText string
		if m.confirmAction == "install" {
			actionText = "安装 cchline 到 Claude Code?"
		} else {
			actionText = "从 Claude Code 卸载 cchline?"
		}
		help = helpBarStyle.Render(
			valueStyle.Render(actionText) + "  " +
				keyStyle.Render("y") + " 确认  " +
				keyStyle.Render("n") + " 取消",
		)
	} else if m.editing {
		help = helpBarStyle.Render(
			keyStyle.Render("Enter") + " Save  " +
				keyStyle.Render("Ctrl+U") + " Clear  " +
				keyStyle.Render("Esc") + " Cancel",
		)
	} else {
		help = helpBarStyle.Render(
			keyStyle.Render("↑↓") + " Navigate  " +
				keyStyle.Render("Space") + " Toggle  " +
				keyStyle.Render("Enter") + " Edit Text\n" +
				keyStyle.Render(ReorderKeyHint) + " Move  " +
				keyStyle.Render("a") + " Add Break  " +
				keyStyle.Render("d") + " Del Break\n" +
				keyStyle.Render("i") + " Install  " +
				keyStyle.Render("u") + " Uninstall  " +
				keyStyle.Render("Esc") + " Exit",
		)
	}

	// 状态消息
	var status string
	if m.statusMessage != "" {
		status = "\n  " + m.statusMessage
	}

	// 调试信息
	debug := fmt.Sprintf("\n  DEBUG: Last key = %q", m.debugKey)

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n%s%s%s\n", title, previewLabel, preview, box, installStatus, help, status, debug)
}

// Run 运行 TUI
func Run(cfg *config.SimpleConfig) error {
	p := tea.NewProgram(NewModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
