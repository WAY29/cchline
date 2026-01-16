//go:build darwin

package tui

// macOS 使用 Cmd 键
const (
	// MoveUpKey 向上移动 segment 的按键
	MoveUpKey = "cmd+up"
	// MoveDownKey 向下移动 segment 的按键
	MoveDownKey = "cmd+down"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "Cmd+↑↓"
)
