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
	"github.com/charmbracelet/x/ansi"
)

type installStatusMsg struct {
	status       installStatus
	installedVer *semVer
	currentVer   *semVer
}

func fetchInstallStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, installedVer, currentVer := getInstallStatus()
		return installStatusMsg{
			status:       status,
			installedVer: installedVer,
			currentVer:   currentVer,
		}
	}
}

// 样式定义
var (
	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("236")).
			Padding(0, 2)

	// 分组标题
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

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
			Padding(0, 2)

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
			Padding(0, 2)

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
	label        string
	key          string
	enabled      bool
	isHeader     bool
	isSeparator  bool   // 是否为分隔符选项
	isSelected   bool   // 分隔符是否被选中
	isTextInput  bool   // 是否为文本输入项
	textKey      string // 文本输入的配置键名 ("cch_url" 或 "cch_api_key")
	isSegmentRow bool   // 是否为 segments 行
	rowIndex     int    // segmentRows 行号（仅 isSegmentRow 时有效）
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
	confirmRow    int                        // 待确认删除的行号（仅 confirmAction == "delete_row" 时有效）
	statusMessage string                     // 操作结果消息

	segmentRows [][]segmentEntry // segments 按行存储
	segmentCol  int              // 当前行内选中的 segment 下标

	segmentPickerOpen   bool
	segmentPickerCursor int
	segmentPickerInput  textinput.Model

	installStatusLoading bool
	installStatusValue   installStatus
	installInstalledVer  *semVer
	installCurrentVer    *semVer
}

// NewModel 创建新的 TUI 模型
func NewModel(cfg *config.SimpleConfig) Model {
	cfg.SegmentEnabled = config.NormalizeSegmentEnabled(cfg.SegmentOrder, cfg.SegmentEnabled)

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

	segmentRows := segmentOrderToRows(cfg.SegmentOrder, cfg.SegmentEnabled)
	for i := range segmentRows {
		items = append(items, menuItem{
			isSegmentRow: true,
			rowIndex:     i,
		})
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

	segmentPickerInput := textinput.New()
	segmentPickerInput.Placeholder = "Type to filter"
	segmentPickerInput.Prompt = "Search: "
	segmentPickerInput.CharLimit = 64

	// 找到第一个非 header 项
	cursor := 0
	for i, item := range items {
		if !item.isHeader {
			cursor = i
			break
		}
	}

	return Model{
		config:      cfg,
		cursor:      cursor,
		items:       items,
		width:       0,
		height:      0,
		textInputs:  textInputs,
		confirmRow:  -1,
		segmentRows: segmentRows,
		segmentCol:  0,

		segmentPickerOpen:   false,
		segmentPickerCursor: 0,
		segmentPickerInput:  segmentPickerInput,

		installStatusLoading: true,
	}
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func safeWidth(w int) int {
	if w <= 0 {
		return 80
	}
	return w
}

func safeHeight(h int) int {
	if h <= 0 {
		return 24
	}
	return h
}

func limitANSILines(s string, maxLines, width int) string {
	if maxLines <= 0 || width <= 0 || s == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	overflow := len(lines) > maxLines
	if overflow {
		lines = lines[:maxLines]
	}

	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], width, "")
	}

	if overflow {
		last := maxLines - 1
		if width >= 3 {
			lines[last] = ansi.Truncate(lines[last], width-3, "") + "..."
		} else {
			lines[last] = strings.Repeat(".", width)
		}
	}

	return strings.Join(lines, "\n")
}

type layout struct {
	titleLines      []string
	previewBoxLines []string

	menuBoxLines []string

	installStatus string
	helpLines     []string
	statusLine    string
	debugLine     string
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func (m Model) buildHelpText() string {
	if m.confirmAction != "" {
		var actionText string
		if m.confirmAction == "install" {
			actionText = "安装 cchline 到 Claude Code?"
		} else if m.confirmAction == "delete_row" {
			actionText = fmt.Sprintf("删除第 %d 行?", m.confirmRow+1)
		} else {
			actionText = "从 Claude Code 卸载 cchline?"
		}
		return valueStyle.Render(actionText) + "  " +
			keyStyle.Render("y") + " 确认  " +
			keyStyle.Render("n") + " 取消"
	}

	if m.segmentPickerOpen {
		return keyStyle.Render("↑↓") + " Select  " +
			keyStyle.Render("Tab") + " Next  " +
			keyStyle.Render("Shift+Tab") + " Prev  " +
			keyStyle.Render("Enter") + " Apply  " +
			keyStyle.Render("Esc") + " Close\n" +
			keyStyle.Render("Type") + " Filter"
	}

	if m.editing {
		return keyStyle.Render("Enter") + " Save  " +
			keyStyle.Render("Ctrl+U") + " Clear  " +
			keyStyle.Render("Esc") + " Cancel"
	}

	return keyStyle.Render("↑↓") + " Navigate  " +
		keyStyle.Render("←→") + " Select  " +
		keyStyle.Render("Space") + " Toggle  " +
		keyStyle.Render("Tab") + " Cycle  " +
		keyStyle.Render("Shift+Tab") + " Reverse  " +
		keyStyle.Render("t") + " Pick\n" +
		keyStyle.Render("a") + " Add(model)  " +
		keyStyle.Render("A") + " Add(cch_model)  " +
		keyStyle.Render("d") + " Delete  " +
		keyStyle.Render(ReorderKeyHint) + " Move\n" +
		keyStyle.Render("n") + " New Line  " +
		keyStyle.Render("x") + " Del Line  " +
		keyStyle.Render("Enter") + " Edit  " +
		keyStyle.Render("i") + " Install  " +
		keyStyle.Render("u") + " Uninstall  " +
		keyStyle.Render("Esc") + " Exit"
}

func (m Model) buildMenuLines(innerWidth int) (lines []string, itemLineIndex []int) {
	itemLineIndex = make([]int, len(m.items))

	appendLine := func(s string) {
		lines = append(lines, ansi.Truncate(s, innerWidth, ""))
	}

	for i := 0; i < len(m.items); i++ {
		item := m.items[i]
		if item.isHeader {
			if len(lines) > 0 {
				lines = append(lines, "")
			}
			itemLineIndex[i] = len(lines)
			appendLine(sectionStyle.Render("  " + item.label))
			continue
		}

		// Collapse separator presets into a single line (plus a preview line).
		if item.isSeparator {
			start := i
			end := i
			for end+1 < len(m.items) && m.items[end+1].isSeparator {
				end++
			}

			lineIndex := len(lines)
			for j := start; j <= end; j++ {
				itemLineIndex[j] = lineIndex
			}

			tokens := make([]string, 0, end-start+1)
			for j := start; j <= end; j++ {
				it := m.items[j]
				cursorMark := " "
				if m.cursor == j {
					cursorMark = selectedStyle.Render("▸")
				}
				dot := disabledStyle.Render("○")
				if it.isSelected {
					dot = enabledStyle.Render("●")
				}
				label := normalStyle.Render(it.label)
				if m.cursor == j {
					label = selectedStyle.Render(it.label)
				}
				tokens = append(tokens, fmt.Sprintf("%s%s %s", cursorMark, dot, label))
			}

			appendLine("  " + strings.Join(tokens, "  "))

			preview := "A" + m.config.Separator + "B"
			appendLine("  " + valueStyle.Render(fmt.Sprintf("%q", m.config.Separator)) + "  " + normalStyle.Render(preview))

			i = end
			continue
		}

		itemLineIndex[i] = len(lines)

		cursor := "   "
		if m.cursor == i {
			cursor = " ▸ "
		}

		var line string
		if item.isTextInput {
			input := m.textInputs[item.textKey]
			var displayValue string

			if m.editing && m.cursor == i {
				displayValue = input.View()
			} else {
				value := input.Value()
				if item.textKey == "cch_api_key" && value != "" {
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
		} else if item.isSegmentRow {
			row := item.rowIndex
			var segs []segmentEntry
			if row >= 0 && row < len(m.segmentRows) {
				segs = m.segmentRows[row]
			}

			prefix := disabledStyle.Render(fmt.Sprintf("L%d:", row+1))
			if m.cursor != i {
				prefix = disabledStyle.Render(fmt.Sprintf("L%d:", row+1))
			}

			if len(segs) == 0 {
				line = fmt.Sprintf("%s%s %s", cursor, prefix, disabledStyle.Render("(empty)"))
			} else {
				tokens := make([]string, 0, len(segs))
				for idx, seg := range segs {
					icon := disabledStyle.Render("○")
					style := normalStyle
					if seg.enabled {
						icon = enabledStyle.Render("●")
					} else {
						style = disabledStyle
					}
					if m.cursor == i && idx == m.segmentCol {
						style = selectedStyle
					}
					tokens = append(tokens, icon+style.Render(seg.name))
				}
				line = fmt.Sprintf("%s%s %s", cursor, prefix, strings.Join(tokens, "  "))
			}
		} else {
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

		appendLine(line)
	}

	return lines, itemLineIndex
}

func (m Model) filteredSegmentChoices() []string {
	query := strings.TrimSpace(strings.ToLower(m.segmentPickerInput.Value()))
	if query == "" {
		return segmentCycle
	}

	parts := strings.Fields(query)
	var matches []string
	for _, name := range segmentCycle {
		lowerName := strings.ToLower(name)
		ok := true
		for _, p := range parts {
			if !strings.Contains(lowerName, p) {
				ok = false
				break
			}
		}
		if ok {
			matches = append(matches, name)
		}
	}
	return matches
}

func (m Model) buildSegmentPickerLines(innerWidth, innerHeight int) []string {
	appendLine := func(lines *[]string, s string) {
		*lines = append(*lines, ansi.Truncate(s, innerWidth, ""))
	}

	var lines []string
	appendLine(&lines, sectionStyle.Render("  SEGMENT PICKER"))
	appendLine(&lines, "  "+m.segmentPickerInput.View())

	listHeight := innerHeight - len(lines)
	if listHeight < 1 {
		if len(lines) > innerHeight {
			lines = lines[:innerHeight]
		}
		return lines
	}

	choices := m.filteredSegmentChoices()
	if len(choices) == 0 {
		appendLine(&lines, "  "+disabledStyle.Render("(no matches)"))
		for len(lines) < innerHeight {
			appendLine(&lines, "")
		}
		return lines
	}

	cursor := clampInt(m.segmentPickerCursor, 0, len(choices)-1)
	start := 0
	if len(choices) > listHeight {
		start = cursor - listHeight/2
		start = clampInt(start, 0, len(choices)-listHeight)
	}
	end := len(choices)
	if len(choices) > listHeight {
		end = start + listHeight
	}

	for i := start; i < end; i++ {
		mark := "   "
		style := normalStyle
		if i == cursor {
			mark = " ▸ "
			style = selectedStyle
		}
		appendLine(&lines, mark+style.Render(choices[i]))
	}

	for len(lines) < innerHeight {
		appendLine(&lines, "")
	}
	if len(lines) > innerHeight {
		lines = lines[:innerHeight]
	}
	return lines
}

func (m Model) buildLayout(width, height int) layout {
	termWidth := safeWidth(width)
	height = safeHeight(height)

	// Avoid printing into the last terminal column to prevent autowrap.
	totalWidth := termWidth
	if totalWidth > 1 {
		totalWidth--
	}

	// lipgloss Style.Width is applied *before* borders are added, so for bordered
	// blocks we must subtract the border sizes to achieve an overall total width.
	totalToContentWidth := func(s lipgloss.Style, total int) int {
		w := total - s.GetHorizontalBorderSize() - s.GetHorizontalMargins()
		if w < 0 {
			return 0
		}
		return w
	}

	contentToInnerWidth := func(s lipgloss.Style, content int) int {
		w := content - s.GetHorizontalPadding()
		if w < 1 {
			return 1
		}
		return w
	}

	menuContentWidth := totalToContentWidth(boxStyle, totalWidth)
	menuBoxStyle := boxStyle.Width(menuContentWidth)
	menuFrameY := menuBoxStyle.GetVerticalFrameSize()

	previewMaxLines := 3
	helpMaxLines := 3
	showStatus := m.statusMessage != ""
	showDebug := m.debugKey != ""

	var (
		title         string
		previewBox    string
		installStatus string
		helpBox       string
		statusLine    string
		debugLine     string
	)

	for {
		title = titleStyle.Width(totalWidth).Render(" CCHLine Configuration ")

		previewContentWidth := totalToContentWidth(previewStyle, totalWidth)
		previewBoxStyle := previewStyle.Width(previewContentWidth)
		previewInnerWidth := contentToInnerWidth(previewStyle, previewContentWidth)
		previewContent := limitANSILines(m.generatePreview(), previewMaxLines, previewInnerWidth)
		previewBox = previewBoxStyle.Render(previewContent)

		if m.installStatusLoading {
			installStatus = "  " + disabledStyle.Render("◌ 检测中")
		} else {
			switch m.installStatusValue {
			case installStatusInstalled:
				installStatus = "  " + enabledStyle.Render("● 已安装")
				if m.installInstalledVer != nil && m.installCurrentVer != nil {
					installStatus += " " + disabledStyle.Render(fmt.Sprintf("(installed %d.%d.%d, current %d.%d.%d)", m.installInstalledVer.major, m.installInstalledVer.minor, m.installInstalledVer.patch, m.installCurrentVer.major, m.installCurrentVer.minor, m.installCurrentVer.patch))
				}
			case installStatusOutdated:
				installStatus = "  " + valueStyle.Render("◐ 版本过旧")
				if m.installInstalledVer != nil && m.installCurrentVer != nil {
					installStatus += " " + disabledStyle.Render(fmt.Sprintf("(installed %d.%d.%d < current %d.%d.%d)", m.installInstalledVer.major, m.installInstalledVer.minor, m.installInstalledVer.patch, m.installCurrentVer.major, m.installCurrentVer.minor, m.installCurrentVer.patch))
				}
			case installStatusUnknown:
				installStatus = "  " + valueStyle.Render("◑ 版本未知")
			default:
				installStatus = "  " + disabledStyle.Render("○ 未安装")
			}
		}
		installStatus = ansi.Truncate(installStatus, totalWidth, "")

		helpBoxStyle := helpBarStyle.Width(totalWidth)
		helpInnerWidth := totalWidth - helpBarStyle.GetHorizontalPadding()
		if helpInnerWidth < 1 {
			helpInnerWidth = 1
		}
		helpText := limitANSILines(m.buildHelpText(), helpMaxLines, helpInnerWidth)
		helpBox = helpBoxStyle.Render(helpText)

		statusLine = ""
		if showStatus {
			statusLine = ansi.Truncate("  "+m.statusMessage, totalWidth, "")
		}

		debugLine = ""
		if showDebug {
			debugLine = ansi.Truncate(fmt.Sprintf("  DEBUG: Last key = %q", m.debugKey), totalWidth, "")
		}

		topHeight := lipgloss.Height(title) + 1 + lipgloss.Height(previewBox)
		bottomHeight := 1 + lipgloss.Height(helpBox)
		if showStatus {
			bottomHeight += 1
		}
		if showDebug {
			bottomHeight += 1
		}

		remaining := height - topHeight - bottomHeight
		if remaining >= menuFrameY+1 || (previewMaxLines <= 1 && helpMaxLines <= 1 && !showStatus && !showDebug) {
			break
		}

		if showDebug {
			showDebug = false
			continue
		}
		if showStatus {
			showStatus = false
			continue
		}
		if helpMaxLines > 1 {
			helpMaxLines--
			continue
		}
		if previewMaxLines > 1 {
			previewMaxLines--
			continue
		}
		break
	}

	topHeight := lipgloss.Height(title) + 1 + lipgloss.Height(previewBox)
	bottomHeight := 1 + lipgloss.Height(helpBox)
	if showStatus {
		bottomHeight += 1
	}
	if showDebug {
		bottomHeight += 1
	}

	available := height - topHeight - bottomHeight
	if available < 0 {
		available = 0
	}

	menuInnerHeight := available - menuFrameY
	if menuInnerHeight < 1 {
		menuInnerHeight = 1
	}

	menuInnerWidth := contentToInnerWidth(boxStyle, menuContentWidth)
	menuLines, itemLineIndex := m.buildMenuLines(menuInnerWidth)

	cursorLine := 0
	if m.cursor >= 0 && m.cursor < len(itemLineIndex) {
		cursorLine = itemLineIndex[m.cursor]
	}

	if m.segmentPickerOpen {
		menuLines = m.buildSegmentPickerLines(menuInnerWidth, menuInnerHeight)
		itemLineIndex = nil
		cursorLine = 0
	}

	start := 0
	if len(menuLines) > menuInnerHeight && menuInnerHeight > 0 {
		if cursorLine >= menuInnerHeight {
			start = cursorLine - menuInnerHeight + 1
		}
		start = clampInt(start, 0, len(menuLines)-menuInnerHeight)
	}

	end := len(menuLines)
	if len(menuLines) > menuInnerHeight {
		end = start + menuInnerHeight
	}

	visibleMenuLines := menuLines[start:end]
	menuBox := menuBoxStyle.Render(strings.Join(visibleMenuLines, "\n"))

	return layout{
		titleLines:      splitLines(title),
		previewBoxLines: splitLines(previewBox),
		menuBoxLines:    splitLines(menuBox),
		installStatus:   installStatus,
		helpLines:       splitLines(helpBox),
		statusLine:      statusLine,
		debugLine:       debugLine,
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

	enabledIdx := 0
	for _, name := range m.config.SegmentOrder {
		// 处理换行分隔符
		if name == config.LineBreakMarker {
			results = append(results, segment.SegmentResult{
				ID:   config.SegmentLineBreak,
				Data: segment.SegmentData{},
			})
			continue
		}

		instanceEnabled := true
		if enabledIdx < len(m.config.SegmentEnabled) {
			instanceEnabled = m.config.SegmentEnabled[enabledIdx]
		}
		enabledIdx++
		if !instanceEnabled {
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

// Init 初始化
func (m Model) Init() tea.Cmd {
	return fetchInstallStatusCmd()
}

// Update 更新
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case installStatusMsg:
		m.installStatusLoading = false
		m.installStatusValue = msg.status
		m.installInstalledVer = msg.installedVer
		m.installCurrentVer = msg.currentVer

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
				needsRefreshInstallStatus := false
				if m.confirmAction == "install" {
					err = m.installStatusLine()
					if err == nil {
						m.statusMessage = "✓ 安装成功"
						needsRefreshInstallStatus = true
					} else {
						m.statusMessage = "✗ " + err.Error()
					}
				} else if m.confirmAction == "uninstall" {
					err = m.uninstallStatusLine()
					if err == nil {
						m.statusMessage = "✓ 卸载成功"
						needsRefreshInstallStatus = true
					} else {
						m.statusMessage = "✗ " + err.Error()
					}
				} else if m.confirmAction == "delete_row" {
					m.deleteRowConfirmed()
					m.statusMessage = "✓ 已删除"
				}
				m.confirmAction = ""
				m.confirmRow = -1
				if needsRefreshInstallStatus {
					m.installStatusLoading = true
					return m, fetchInstallStatusCmd()
				}
				return m, nil
			case "n", "N", "esc":
				// 取消操作
				m.confirmAction = ""
				m.confirmRow = -1
				m.statusMessage = ""
				return m, nil
			default:
				// 忽略其他按键
				return m, nil
			}
		}

		// 非编辑模式的按键处理
		if m.segmentPickerOpen {
			switch msg.String() {
			case "esc":
				m.segmentPickerOpen = false
				m.segmentPickerInput.Blur()
				return m, nil

			case "enter":
				choices := m.filteredSegmentChoices()
				if len(choices) == 0 {
					return m, nil
				}
				m.segmentPickerCursor = clampInt(m.segmentPickerCursor, 0, len(choices)-1)
				m.setCurrentSegmentName(choices[m.segmentPickerCursor])
				m.segmentPickerOpen = false
				m.segmentPickerInput.Blur()
				return m, nil

			case "up", "k", "shift+tab", "backtab":
				choices := m.filteredSegmentChoices()
				if len(choices) == 0 {
					return m, nil
				}
				m.segmentPickerCursor--
				if m.segmentPickerCursor < 0 {
					m.segmentPickerCursor = len(choices) - 1
				}
				return m, nil

			case "down", "j", "tab":
				choices := m.filteredSegmentChoices()
				if len(choices) == 0 {
					return m, nil
				}
				m.segmentPickerCursor++
				if m.segmentPickerCursor >= len(choices) {
					m.segmentPickerCursor = 0
				}
				return m, nil

			default:
				var cmd tea.Cmd
				m.segmentPickerInput, cmd = m.segmentPickerInput.Update(msg)
				choices := m.filteredSegmentChoices()
				if len(choices) == 0 {
					m.segmentPickerCursor = 0
				} else {
					m.segmentPickerCursor = clampInt(m.segmentPickerCursor, 0, len(choices)-1)
				}
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			// 退出时自动保存
			m.saveConfig()
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.moveCursorVertical(-1)

		case "down", "j":
			m.moveCursorVertical(1)

		case "left", "h":
			m.moveCursorHorizontal(-1)

		case "right", "l":
			m.moveCursorHorizontal(1)

		case "a":
			m.insertSegmentAfterCurrent("model")

		case "A":
			m.insertSegmentAfterCurrent("cch_model")

		case "d":
			m.deleteCurrentSegment()

		case "tab":
			m.cycleCurrentSegment()

		case "shift+tab", "backtab":
			m.cycleCurrentSegmentReverse()

		case "t":
			item := m.items[m.cursor]
			if item.isSegmentRow {
				m.segmentPickerInput.SetValue("")
				m.segmentPickerInput.Focus()
				m.segmentPickerOpen = true
				m.segmentPickerCursor = 0

				row := item.rowIndex
				if row >= 0 && row < len(m.segmentRows) && m.segmentCol >= 0 && m.segmentCol < len(m.segmentRows[row]) {
					current := m.segmentRows[row][m.segmentCol].name
					for i := range segmentCycle {
						if segmentCycle[i] == current {
							m.segmentPickerCursor = i
							break
						}
					}
				}
				return m, nil
			}

		case MoveLeftKey:
			m.moveSegmentWithinRow(-1)

		case MoveRightKey:
			m.moveSegmentWithinRow(1)

		case "n":
			m.insertRowAfterCurrent()

		case "x":
			m.requestDeleteCurrentRow()

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
			} else if item.isSegmentRow {
				m.toggleCurrentSegmentInstanceEnabled()
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

	lay := m.buildLayout(m.width, m.height)

	var lines []string
	lines = append(lines, lay.titleLines...)
	lines = append(lines, lay.previewBoxLines...)
	lines = append(lines, lay.menuBoxLines...)
	lines = append(lines, lay.installStatus)
	lines = append(lines, lay.helpLines...)
	if lay.statusLine != "" {
		lines = append(lines, lay.statusLine)
	}
	if lay.debugLine != "" {
		lines = append(lines, lay.debugLine)
	}

	h := safeHeight(m.height)
	if len(lines) < h {
		lines = append(lines, make([]string, h-len(lines))...)
	} else if len(lines) > h {
		lines = lines[:h]
	}

	w := safeWidth(m.width)
	if w > 1 {
		w--
	}
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], w, "")
	}

	return strings.Join(lines, "\n")
}

// Run 运行 TUI
func Run(cfg *config.SimpleConfig) error {
	p := tea.NewProgram(NewModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
