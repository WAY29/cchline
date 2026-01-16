package segment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WAY29/cchline/config"
)

type ContextWindowSegment struct{}

func (s *ContextWindowSegment) Collect(input *config.InputData) SegmentData {
	contextLimit := getModelContextLimit(input.Model.ID)

	// 解析 transcript 获取 token 使用量
	totalTokens := parseTranscriptUsage(input.TranscriptPath)

	if totalTokens == 0 {
		// 无数据时显示 "-"
		return SegmentData{
			Primary: "- · - tokens",
			Metadata: map[string]string{
				"tokens": "-",
				"limit":  fmt.Sprintf("%d", contextLimit),
			},
		}
	}

	// 计算使用率
	percentage := float64(totalTokens) / float64(contextLimit) * 100

	// 格式化百分比
	var percentageStr string
	if percentage == float64(int(percentage)) {
		percentageStr = fmt.Sprintf("%.0f%%", percentage)
	} else {
		percentageStr = fmt.Sprintf("%.1f%%", percentage)
	}

	// 格式化 token 数量
	tokensStr := formatTokenCount(totalTokens)

	return SegmentData{
		Primary: fmt.Sprintf("%s · %s tokens", percentageStr, tokensStr),
		Metadata: map[string]string{
			"tokens": fmt.Sprintf("%d", totalTokens),
			"limit":  fmt.Sprintf("%d", contextLimit),
		},
	}
}

// TranscriptMessage 表示 transcript 中的消息结构
type TranscriptMessage struct {
	Type     string          `json:"type"`
	UUID     string          `json:"uuid,omitempty"`
	LeafUUID string          `json:"leafUuid,omitempty"`
	Message  *MessageContent `json:"message,omitempty"`
}

// MessageContent 消息内容
type MessageContent struct {
	Usage *UsageData `json:"usage,omitempty"`
}

// UsageData token 使用数据
type UsageData struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// displayTokens 计算显示用的 token 总数
func (u *UsageData) displayTokens() int {
	// 总 token = input + output + cache_creation + cache_read
	return u.InputTokens + u.OutputTokens + u.CacheCreationInputTokens + u.CacheReadInputTokens
}

// parseTranscriptUsage 解析 transcript 文件获取 token 使用量
func parseTranscriptUsage(transcriptPath string) int {
	if transcriptPath == "" {
		return 0
	}

	// 尝试从当前 transcript 文件解析
	if usage := tryParseTranscriptFile(transcriptPath); usage > 0 {
		return usage
	}

	// 如果文件不存在，尝试从项目历史中查找
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		if usage := tryFindUsageFromProjectHistory(transcriptPath); usage > 0 {
			return usage
		}
	}

	return 0
}

// tryParseTranscriptFile 尝试解析 transcript 文件
func tryParseTranscriptFile(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	// 增加 buffer 大小以处理大行
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return 0
	}

	// 检查最后一行是否是 summary
	lastLine := strings.TrimSpace(lines[len(lines)-1])
	var entry TranscriptMessage
	if err := json.Unmarshal([]byte(lastLine), &entry); err == nil {
		if entry.Type == "summary" && entry.LeafUUID != "" {
			// summary 情况：通过 leafUuid 查找 usage
			projectDir := filepath.Dir(path)
			return findUsageByLeafUUID(entry.LeafUUID, projectDir)
		}
	}

	// 正常情况：从后往前找最后一条 assistant 消息
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		var msg TranscriptMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		if msg.Type == "assistant" && msg.Message != nil && msg.Message.Usage != nil {
			return msg.Message.Usage.displayTokens()
		}
	}

	return 0
}

// findUsageByLeafUUID 通过 leafUUID 在项目目录中查找 usage
func findUsageByLeafUUID(leafUUID string, projectDir string) int {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return 0
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		path := filepath.Join(projectDir, entry.Name())
		if usage := searchUUIDInFile(path, leafUUID); usage > 0 {
			return usage
		}
	}

	return 0
}

// searchUUIDInFile 在文件中搜索指定 UUID 的 usage
func searchUUIDInFile(path string, targetUUID string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var msg TranscriptMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		if msg.UUID == targetUUID {
			if msg.Type == "assistant" && msg.Message != nil && msg.Message.Usage != nil {
				return msg.Message.Usage.displayTokens()
			}
		}
	}

	return 0
}

// tryFindUsageFromProjectHistory 从项目历史中查找 usage
func tryFindUsageFromProjectHistory(transcriptPath string) int {
	projectDir := filepath.Dir(transcriptPath)

	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return 0
	}

	// 收集所有 jsonl 文件
	type fileInfo struct {
		path    string
		modTime int64
	}
	var files []fileInfo

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		path := filepath.Join(projectDir, entry.Name())
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		files = append(files, fileInfo{path: path, modTime: info.ModTime().Unix()})
	}

	if len(files) == 0 {
		return 0
	}

	// 按修改时间排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime > files[j].modTime
	})

	// 从最新的文件开始查找
	for _, f := range files {
		if usage := tryParseTranscriptFile(f.path); usage > 0 {
			return usage
		}
	}

	return 0
}

var modelContextLimits = map[string]int{
	"claude-3-opus":      200000,
	"claude-3-sonnet":    200000,
	"claude-3-haiku":     200000,
	"claude-3-5-sonnet":  200000,
	"claude-3-5-haiku":   200000,
	"claude-opus-4":      200000,
	"claude-sonnet-4":    200000,
	"gpt-4":              128000,
	"gpt-4-turbo":        128000,
	"gpt-4o":             128000,
}

func getModelContextLimit(modelID string) int {
	modelID = strings.ToLower(modelID)
	for pattern, limit := range modelContextLimits {
		if strings.Contains(modelID, pattern) {
			return limit
		}
	}
	return 200000
}

// 辅助函数供其他 segment 使用

// readAllLines 读取文件所有行
func readAllLines(file *os.File) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// parseTranscriptLine 解析单行 transcript
func parseTranscriptLine(line string) *TranscriptMessage {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var msg TranscriptMessage
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return nil
	}
	return &msg
}
