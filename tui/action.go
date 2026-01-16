package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/WAY29/cchline/config"
)

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

	// 边界检查，处理初始越界
	if newCursor < 0 || newCursor >= len(m.items) {
		// 直接越界，触发循环导航
		if delta < 0 {
			// 向上越界，跳到最后一个非 header 项
			newCursor = len(m.items) - 1
			for newCursor >= 0 && m.items[newCursor].isHeader {
				newCursor--
			}
		} else {
			// 向下越界，跳到第一个非 header 项
			newCursor = 0
			for newCursor < len(m.items) && m.items[newCursor].isHeader {
				newCursor++
			}
		}
		if newCursor >= 0 && newCursor < len(m.items) {
			m.cursor = newCursor
		}
		return
	}

	// 跳过 header
	for newCursor >= 0 && newCursor < len(m.items) && m.items[newCursor].isHeader {
		newCursor += delta
	}

	// 跳过 header 后越界，触发循环导航
	if newCursor < 0 {
		// 向上越界，跳到最后一个非 header 项
		newCursor = len(m.items) - 1
		for newCursor >= 0 && m.items[newCursor].isHeader {
			newCursor--
		}
	} else if newCursor >= len(m.items) {
		// 向下越界，跳到第一个非 header 项
		newCursor = 0
		for newCursor < len(m.items) && m.items[newCursor].isHeader {
			newCursor++
		}
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

// insertLineBreak 在当前 segment 后插入换行分隔符
func (m *Model) insertLineBreak() {
	item := m.items[m.cursor]
	// 只能在 segment 项后插入（非 header、非 separator、非 theme、非 textInput、非 lineBreak）
	if item.isHeader || item.isSeparator || item.key == "theme" || item.isTextInput || item.isLineBreak {
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
	if segmentsStart < 0 || m.cursor < segmentsStart || m.cursor >= segmentsEnd {
		return
	}

	// 在当前位置后插入换行分隔符
	newItem := menuItem{
		label:       "Line Break",
		key:         config.LineBreakMarker,
		isLineBreak: true,
	}

	// 插入到 cursor+1 位置
	insertPos := m.cursor + 1
	m.items = append(m.items[:insertPos], append([]menuItem{newItem}, m.items[insertPos:]...)...)

	// 更新配置中的 SegmentOrder
	m.updateSegmentOrder()
}

// deleteLineBreak 删除当前换行分隔符
func (m *Model) deleteLineBreak() {
	item := m.items[m.cursor]
	// 只能删除换行分隔符
	if !item.isLineBreak {
		return
	}

	// 删除当前项
	m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)

	// 调整光标位置
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	// 跳过 header
	for m.cursor >= 0 && m.items[m.cursor].isHeader {
		m.cursor--
	}

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

// getExecutablePath 获取当前可执行文件的完整路径
func getExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

// getSettingsPath 获取 Claude settings.json 路径
func getSettingsPath() string {
	return filepath.Join(os.Getenv("HOME"), ".claude", "settings.json")
}

// installStatusLine 安装 statusLine 到 settings.json
func (m *Model) installStatusLine() error {
	exePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}

	settingsPath := getSettingsPath()

	// 读取现有配置
	var settings map[string]any
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]any)
		} else {
			return fmt.Errorf("读取 settings.json 失败: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("解析 settings.json 失败: %w", err)
		}
	}

	// 添加 statusLine 配置
	settings["statusLine"] = map[string]any{
		"type":    "command",
		"command": exePath,
		"padding": 0,
	}

	// 写回文件
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 settings.json 失败: %w", err)
	}

	if err := os.WriteFile(settingsPath, newData, 0644); err != nil {
		return fmt.Errorf("写入 settings.json 失败: %w", err)
	}

	return nil
}

// uninstallStatusLine 从 settings.json 移除 statusLine
func (m *Model) uninstallStatusLine() error {
	settingsPath := getSettingsPath()

	// 读取现有配置
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，无需卸载
		}
		return fmt.Errorf("读取 settings.json 失败: %w", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("解析 settings.json 失败: %w", err)
	}

	// 删除 statusLine 字段
	delete(settings, "statusLine")

	// 写回文件
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 settings.json 失败: %w", err)
	}

	if err := os.WriteFile(settingsPath, newData, 0644); err != nil {
		return fmt.Errorf("写入 settings.json 失败: %w", err)
	}

	return nil
}

// isInstalled 检查 cchline 是否已安装到 Claude Code
func isInstalled() bool {
	settingsPath := getSettingsPath()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return false
	}

	_, exists := settings["statusLine"]
	return exists
}
