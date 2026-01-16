//go:build darwin

package tui

// macOS 使用 Option+j/k (在终端中被识别为 Alt)
const (
	// MoveUpKey 向上移动 segment 的按键
	MoveUpKey = "alt+k"
	// MoveDownKey 向下移动 segment 的按键
	MoveDownKey = "alt+j"
	// ReorderKeyHint 帮助栏显示的按键提示
	ReorderKeyHint = "⌥+j/k"
)
