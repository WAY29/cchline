package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func (m *Model) separatorBounds() (headerIndex, start, end int, ok bool) {
	headerIndex = -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEPARATOR" {
			headerIndex = i
			break
		}
	}
	if headerIndex < 0 {
		return -1, 0, 0, false
	}

	start = headerIndex + 1
	if start >= len(m.items) || !m.items[start].isSeparator {
		return -1, 0, 0, false
	}

	end = start
	for end+1 < len(m.items) && m.items[end+1].isSeparator {
		end++
	}

	return headerIndex, start, end, true
}

func (m *Model) selectedSeparatorIndex(start, end int) int {
	for i := start; i <= end; i++ {
		if m.items[i].isSeparator && m.items[i].isSelected {
			return i
		}
	}
	return start
}

func (m *Model) moveCursorHorizontal(delta int) {
	item := m.items[m.cursor]
	if item.isSegmentRow {
		row := item.rowIndex
		if row < 0 || row >= len(m.segmentRows) {
			m.segmentCol = 0
			return
		}
		segs := m.segmentRows[row]
		if len(segs) == 0 {
			m.segmentCol = 0
			return
		}
		next := m.segmentCol + delta
		if next < 0 {
			next = len(segs) - 1
		} else if next >= len(segs) {
			next = 0
		}
		m.segmentCol = next
		return
	}

	_, start, end, ok := m.separatorBounds()
	if !ok {
		return
	}
	if m.cursor < start || m.cursor > end {
		return
	}

	next := m.cursor + delta
	if next < start {
		next = end
	} else if next > end {
		next = start
	}
	m.cursor = next
}

func (m *Model) moveCursorVertical(delta int) {
	if delta == 0 {
		return
	}

	headerIndex, start, end, ok := m.separatorBounds()
	if ok && m.cursor >= start && m.cursor <= end {
		if delta < 0 {
			// Jump above the SEPARATOR block.
			newCursor := headerIndex - 1
			for newCursor >= 0 && m.items[newCursor].isHeader {
				newCursor--
			}
			if newCursor < 0 {
				// Wrap to last selectable item.
				newCursor = len(m.items) - 1
				for newCursor >= 0 && m.items[newCursor].isHeader {
					newCursor--
				}
			}
			if newCursor >= 0 && newCursor < len(m.items) {
				m.cursor = newCursor
			}
			m.clampSegmentColForCursor()
			return
		}

		// delta > 0: Jump below the SEPARATOR block.
		newCursor := end + 1
		for newCursor < len(m.items) && m.items[newCursor].isHeader {
			newCursor++
		}
		if newCursor >= len(m.items) {
			// Wrap to first selectable item.
			newCursor = 0
			for newCursor < len(m.items) && m.items[newCursor].isHeader {
				newCursor++
			}
		}
		if newCursor >= 0 && newCursor < len(m.items) {
			m.cursor = newCursor
		}
		m.clampSegmentColForCursor()
		return
	}

	previous := m.cursor
	m.moveCursor(delta)

	// When entering the SEPARATOR block vertically, land on the selected preset.
	if ok && m.cursor >= start && m.cursor <= end {
		if previous < start && delta > 0 {
			m.cursor = m.selectedSeparatorIndex(start, end)
		}
		if previous > end && delta < 0 {
			m.cursor = m.selectedSeparatorIndex(start, end)
		}
	}

	m.clampSegmentColForCursor()
}

var segmentCycle = []string{
	"model",
	"directory",
	"git",
	"context_window",
	"usage",
	"cost",
	"session",
	"output_style",
	"update",
	"cch_model",
	"cch_provider",
	"cch_cost",
	"cch_requests",
	"cch_limits",
}

type segmentEntry struct {
	name    string
	enabled bool
}

func segmentOrderToRows(order []string, enabled []bool) [][]segmentEntry {
	if len(order) == 0 {
		return [][]segmentEntry{{}}
	}

	var rows [][]segmentEntry
	var current []segmentEntry
	enabledIdx := 0
	for _, seg := range order {
		if seg == config.LineBreakMarker {
			rows = append(rows, current)
			current = nil
			continue
		}
		entryEnabled := true
		if enabledIdx < len(enabled) {
			entryEnabled = enabled[enabledIdx]
		}
		enabledIdx++
		current = append(current, segmentEntry{name: seg, enabled: entryEnabled})
	}
	rows = append(rows, current)
	if len(rows) == 0 {
		rows = [][]segmentEntry{{}}
	}
	return rows
}

func rowsToSegmentOrder(rows [][]segmentEntry) ([]string, []bool) {
	if len(rows) == 0 {
		return nil, nil
	}
	if len(rows) == 1 && len(rows[0]) == 0 {
		return nil, nil
	}

	var order []string
	var enabled []bool
	for i, row := range rows {
		for _, entry := range row {
			order = append(order, entry.name)
			enabled = append(enabled, entry.enabled)
		}
		if i != len(rows)-1 {
			order = append(order, config.LineBreakMarker)
		}
	}
	return order, enabled
}

func (m *Model) segmentRowItemsRange() (start, end int, ok bool) {
	segmentsHeader := -1
	cchHeader := -1
	for i, it := range m.items {
		if it.isHeader && it.label == "SEGMENTS" {
			segmentsHeader = i
			continue
		}
		if it.isHeader && it.label == "CCH SETTINGS" {
			cchHeader = i
			break
		}
	}
	if segmentsHeader < 0 {
		return 0, 0, false
	}
	start = segmentsHeader + 1
	if cchHeader < 0 {
		end = len(m.items)
	} else {
		end = cchHeader
	}
	if start > end {
		return 0, 0, false
	}
	return start, end, true
}

func (m *Model) rebuildSegmentRowItems() {
	start, end, ok := m.segmentRowItemsRange()
	if !ok {
		return
	}

	prefix := append([]menuItem(nil), m.items[:start]...)
	suffix := append([]menuItem(nil), m.items[end:]...)

	rows := m.segmentRows
	newItems := make([]menuItem, 0, len(prefix)+len(rows)+len(suffix))
	newItems = append(newItems, prefix...)
	for i := range rows {
		newItems = append(newItems, menuItem{
			isSegmentRow: true,
			rowIndex:     i,
		})
	}
	newItems = append(newItems, suffix...)
	m.items = newItems
}

func (m *Model) currentSegmentRowIndex() (int, bool) {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return 0, false
	}
	item := m.items[m.cursor]
	if !item.isSegmentRow {
		return 0, false
	}
	if item.rowIndex < 0 || item.rowIndex >= len(m.segmentRows) {
		return 0, false
	}
	return item.rowIndex, true
}

func (m *Model) setCursorToSegmentRow(row int) {
	start, _, ok := m.segmentRowItemsRange()
	if !ok {
		return
	}
	if len(m.segmentRows) == 0 {
		return
	}
	row = clampInt(row, 0, len(m.segmentRows)-1)
	m.cursor = start + row
	m.clampSegmentColForCursor()
}

func (m *Model) clampSegmentColForCursor() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]
	if len(entries) == 0 {
		m.segmentCol = 0
		return
	}
	m.segmentCol = clampInt(m.segmentCol, 0, len(entries)-1)
}

func (m *Model) syncSegmentOrder() {
	order, enabled := rowsToSegmentOrder(m.segmentRows)
	m.config.SegmentOrder = order
	m.config.SegmentEnabled = enabled
	m.config.Segments = config.DeriveSegmentToggles(order, enabled)
}

func (m *Model) toggleCurrentSegmentInstanceEnabled() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]
	if len(entries) == 0 {
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	entries[col].enabled = !entries[col].enabled
	m.segmentRows[row] = entries
	m.syncSegmentOrder()
}

func (m *Model) insertSegmentAfterCurrent(seg string) {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}

	entries := m.segmentRows[row]
	insertPos := 0
	if len(entries) == 0 {
		insertPos = 0
	} else {
		insertPos = clampInt(m.segmentCol+1, 0, len(entries))
	}

	newEntry := segmentEntry{name: seg, enabled: true}
	entries = append(entries[:insertPos], append([]segmentEntry{newEntry}, entries[insertPos:]...)...)
	m.segmentRows[row] = entries
	m.segmentCol = insertPos
	m.syncSegmentOrder()
}

func (m *Model) deleteCurrentSegment() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]
	if len(entries) == 0 {
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	entries = append(entries[:col], entries[col+1:]...)
	if len(entries) == 0 {
		m.segmentRows = append(m.segmentRows[:row], m.segmentRows[row+1:]...)
		if len(m.segmentRows) == 0 {
			m.segmentRows = [][]segmentEntry{{}}
			row = 0
		} else if row >= len(m.segmentRows) {
			row = len(m.segmentRows) - 1
		}
		m.segmentCol = 0
		m.rebuildSegmentRowItems()
		m.setCursorToSegmentRow(row)
	} else {
		m.segmentRows[row] = entries
		if col >= len(entries) {
			col = len(entries) - 1
		}
		m.segmentCol = col
	}

	m.syncSegmentOrder()
}

func (m *Model) cycleCurrentSegment() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]

	if len(entries) == 0 {
		m.segmentRows[row] = []segmentEntry{{name: segmentCycle[0], enabled: true}}
		m.segmentCol = 0
		m.syncSegmentOrder()
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	current := entries[col].name
	idx := -1
	for i := range segmentCycle {
		if segmentCycle[i] == current {
			idx = i
			break
		}
	}
	next := segmentCycle[0]
	if idx >= 0 {
		next = segmentCycle[(idx+1)%len(segmentCycle)]
	}

	entries[col].name = next
	entries[col].enabled = true
	m.segmentRows[row] = entries
	m.syncSegmentOrder()
}

func (m *Model) cycleCurrentSegmentReverse() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]

	if len(entries) == 0 {
		m.segmentRows[row] = []segmentEntry{{name: segmentCycle[0], enabled: true}}
		m.segmentCol = 0
		m.syncSegmentOrder()
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	current := entries[col].name
	idx := -1
	for i := range segmentCycle {
		if segmentCycle[i] == current {
			idx = i
			break
		}
	}

	prev := segmentCycle[0]
	if idx >= 0 {
		prev = segmentCycle[(idx-1+len(segmentCycle))%len(segmentCycle)]
	}

	entries[col].name = prev
	entries[col].enabled = true
	m.segmentRows[row] = entries
	m.syncSegmentOrder()
}

func (m *Model) setCurrentSegmentName(name string) {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}

	entries := m.segmentRows[row]
	if len(entries) == 0 {
		m.segmentRows[row] = []segmentEntry{{name: name, enabled: true}}
		m.segmentCol = 0
		m.syncSegmentOrder()
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	entries[col].name = name
	entries[col].enabled = true
	m.segmentRows[row] = entries
	m.syncSegmentOrder()
}

func (m *Model) moveSegmentWithinRow(delta int) {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	entries := m.segmentRows[row]
	if len(entries) < 2 {
		return
	}

	col := clampInt(m.segmentCol, 0, len(entries)-1)
	target := col + delta
	if target < 0 || target >= len(entries) {
		return
	}

	entries[col], entries[target] = entries[target], entries[col]
	m.segmentRows[row] = entries
	m.segmentCol = target
	m.syncSegmentOrder()
}

func (m *Model) insertRowAfterCurrent() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	insertPos := row + 1

	newRow := []segmentEntry{{name: "model", enabled: true}}

	m.segmentRows = append(m.segmentRows[:insertPos], append([][]segmentEntry{newRow}, m.segmentRows[insertPos:]...)...)
	m.rebuildSegmentRowItems()
	m.segmentCol = 0
	m.setCursorToSegmentRow(insertPos)
	m.syncSegmentOrder()
}

func (m *Model) requestDeleteCurrentRow() {
	row, ok := m.currentSegmentRowIndex()
	if !ok {
		return
	}
	m.confirmAction = "delete_row"
	m.confirmRow = row
	m.statusMessage = ""
}

func (m *Model) deleteRowConfirmed() {
	row := m.confirmRow
	if row < 0 || row >= len(m.segmentRows) {
		return
	}

	m.segmentRows = append(m.segmentRows[:row], m.segmentRows[row+1:]...)
	if len(m.segmentRows) == 0 {
		m.segmentRows = [][]segmentEntry{{}}
		row = 0
	} else if row >= len(m.segmentRows) {
		row = len(m.segmentRows) - 1
	}

	m.rebuildSegmentRowItems()
	m.segmentCol = 0
	m.setCursorToSegmentRow(row)
	m.syncSegmentOrder()
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

	statusLine, exists := settings["statusLine"]
	if !exists {
		return false
	}

	// 提取 command 字段
	statusLineMap, ok := statusLine.(map[string]any)
	if !ok {
		return false
	}

	command, ok := statusLineMap["command"].(string)
	if !ok || command == "" {
		return false
	}

	// 执行 command -v 检查输出是否包含 cchline
	out, err := exec.Command(command, "-v").Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), "cchline")
}
