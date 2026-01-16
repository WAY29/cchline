package render

import (
	"fmt"
	"strings"

	"github.com/WAY29/cchline/config"
	"github.com/WAY29/cchline/segment"
)

// StatusLineGenerator generates the status line from segment results
type StatusLineGenerator struct {
	config *config.SimpleConfig
}

// NewStatusLineGenerator creates a new status line generator
func NewStatusLineGenerator(cfg *config.SimpleConfig) *StatusLineGenerator {
	return &StatusLineGenerator{config: cfg}
}

// Generate generates the status line from segment results
func (g *StatusLineGenerator) Generate(segments []segment.SegmentResult) string {
	var lines [][]string
	var currentLine []string

	for _, seg := range segments {
		// 遇到换行标记，保存当前行并开始新行
		if seg.ID == config.SegmentLineBreak {
			if len(currentLine) > 0 {
				lines = append(lines, currentLine)
				currentLine = nil
			}
			continue
		}

		rendered := g.renderSegment(seg)
		if rendered != "" {
			currentLine = append(currentLine, rendered)
		}
	}

	// 添加最后一行
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	// 将每行用分隔符连接，行之间用换行符连接
	var result []string
	for _, line := range lines {
		result = append(result, strings.Join(line, g.config.Separator))
	}

	return strings.Join(result, "\n")
}

// renderSegment renders a single segment with theme and colors
func (g *StatusLineGenerator) renderSegment(seg segment.SegmentResult) string {
	// Get theme configuration based on theme mode
	theme := config.GetSegmentTheme(seg.ID, g.config.Theme)

	// Apply colors to icon and text
	icon := config.ApplyColor(theme.Icon, theme.IconColor)
	text := config.ApplyColor(seg.Data.Primary, theme.TextColor)

	// Apply background color if present
	if theme.BgColor != nil {
		return config.ApplyBackground(fmt.Sprintf(" %s %s ", icon, text), theme.BgColor)
	}

	return fmt.Sprintf("%s %s", icon, text)
}
