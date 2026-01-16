package segment

import (
	"fmt"

	"github.com/WAY29/cchline/config"
)

// SegmentData 段渲染数据
type SegmentData struct {
	Primary   string            // 主要显示内容
	Secondary string            // 次要内容（可选）
	Metadata  map[string]string // 元数据
}

// SegmentResult 段结果
type SegmentResult struct {
	ID   config.SegmentID
	Data SegmentData
}

// Segment 接口定义
type Segment interface {
	Collect(input *config.InputData) SegmentData
}

// TranscriptEntry represents a single entry in the transcript file
type TranscriptEntry struct {
	Role    string     `json:"role"`
	Content string     `json:"content"`
	Usage   *UsageInfo `json:"usage,omitempty"`
}

// UsageInfo contains token usage information
type UsageInfo struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	PromptTokens             int `json:"prompt_tokens"`
	CompletionTokens         int `json:"completion_tokens"`
}

// formatTokenCount 格式化 token 数量
func formatTokenCount(count int) string {
	if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
	if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}
