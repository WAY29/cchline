package segment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/WAY29/cchline/config"
)

type SessionSegment struct{}

func (s *SessionSegment) Collect(input *config.InputData) SegmentData {
	if input.TranscriptPath == "" {
		return SegmentData{}
	}

	startTime := getSessionStartTime(input.TranscriptPath)
	if startTime.IsZero() {
		return SegmentData{}
	}

	duration := time.Since(startTime)
	return SegmentData{
		Primary: formatDuration(duration),
	}
}

func getSessionStartTime(path string) time.Time {
	file, err := os.Open(path)
	if err != nil {
		return time.Time{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		var entry TranscriptEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return time.Time{}
		}
		// 如果有时间戳信息，可以从这里提取
		// 目前返回文件修改时间作为会话开始时间
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return fileInfo.ModTime()
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
