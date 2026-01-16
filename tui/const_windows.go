//go:build windows

package tui

// Windows 使用 Alt 键
const (
	// MoveUpKey 向上移动 segment 的按键
	MoveUpKey = "alt+up"
	// MoveDownKey 向下移动 segment 的按键
	MoveDownKey = "alt+down"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "Alt+↑↓"
)
