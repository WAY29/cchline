//go:build darwin

package tui

// macOS 使用 Option+h/l (在终端中被识别为 Alt)
const (
	// MoveLeftKey 向左移动（行内）segment 的按键
	MoveLeftKey = "alt+h"
	// MoveRightKey 向右移动（行内）segment 的按键
	MoveRightKey = "alt+l"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "⌥+h/l"
)
