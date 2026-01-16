//go:build !darwin && !windows

package tui

// Unix/Linux 使用 Alt+j/k
const (
	// MoveUpKey 向上移动 segment 的按键
	MoveUpKey = "alt+k"
	// MoveDownKey 向下移动 segment 的按键
	MoveDownKey = "alt+j"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "Alt+j/k"
)
