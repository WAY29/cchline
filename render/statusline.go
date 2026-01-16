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
	var output []string

	for _, seg := range segments {
		rendered := g.renderSegment(seg)
		if rendered != "" {
			output = append(output, rendered)
		}
	}

	return strings.Join(output, g.config.Separator)
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
