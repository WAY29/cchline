package segment

import (
	"strings"

	"github.com/WAY29/cchline/config"
)

type ModelSegment struct{}

func (s *ModelSegment) Collect(input *config.InputData) SegmentData {
	name := input.Model.DisplayName
	if name == "" {
		name = input.Model.ID
	}
	name = simplifyModelName(name)
	return SegmentData{Primary: name}
}

func simplifyModelName(name string) string {
	replacements := map[string]string{
		"claude-3-5-sonnet": "Sonnet 3.5",
		"claude-3-opus":     "Opus 3",
		"claude-3-haiku":    "Haiku 3",
		"claude-opus-4":     "Opus 4",
		"claude-sonnet-4":   "Sonnet 4",
		"gpt-4":             "GPT-4",
	}
	lowerName := strings.ToLower(name)
	for k, v := range replacements {
		if strings.Contains(lowerName, k) {
			return v
		}
	}
	return name
}
