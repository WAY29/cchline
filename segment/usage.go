package segment

import (
	"os"

	"github.com/WAY29/cchline/config"
)

type UsageSegment struct{}

func (s *UsageSegment) Collect(input *config.InputData) SegmentData {
	if input.TranscriptPath == "" {
		return SegmentData{}
	}

	inputTokens, outputTokens := parseUsageFromTranscript(input.TranscriptPath)

	if inputTokens == 0 && outputTokens == 0 {
		return SegmentData{}
	}

	return SegmentData{
		Primary: formatUsageDisplay(inputTokens, outputTokens),
	}
}

// parseUsageFromTranscript 解析 transcript 获取输入/输出 token
func parseUsageFromTranscript(transcriptPath string) (inputTokens, outputTokens int) {
	if transcriptPath == "" {
		return 0, 0
	}

	// 复用 context_window 的解析逻辑
	usage := getLastAssistantUsage(transcriptPath)
	if usage == nil {
		return 0, 0
	}

	// 输入 = input_tokens + cache_creation + cache_read
	inputTokens = usage.InputTokens + usage.CacheCreationInputTokens + usage.CacheReadInputTokens
	outputTokens = usage.OutputTokens

	return
}

// getLastAssistantUsage 获取最后一条 assistant 消息的 usage
func getLastAssistantUsage(transcriptPath string) *UsageData {
	file, err := os.Open(transcriptPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	lines := readAllLines(file)
	if len(lines) == 0 {
		return nil
	}

	// 从后往前找最后一条 assistant 消息
	for i := len(lines) - 1; i >= 0; i-- {
		msg := parseTranscriptLine(lines[i])
		if msg != nil && msg.Type == "assistant" && msg.Message != nil && msg.Message.Usage != nil {
			return msg.Message.Usage
		}
	}

	return nil
}

func formatUsageDisplay(inputTokens, outputTokens int) string {
	return "↓" + formatTokenCount(inputTokens) + " ↑" + formatTokenCount(outputTokens)
}
