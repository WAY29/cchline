//go:build !darwin && !windows

package tui

// Unix/Linux 使用 Alt+h/l
const (
	// MoveLeftKey 向左移动（行内）segment 的按键
	MoveLeftKey = "alt+h"
	// MoveRightKey 向右移动（行内）segment 的按键
	MoveRightKey = "alt+l"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "Alt+h/l"
)
