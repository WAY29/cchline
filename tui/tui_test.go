package tui

import (
	"strings"
	"testing"

	"github.com/WAY29/cchline/config"
	"github.com/charmbracelet/x/ansi"
)

func TestViewFitsTerminalAndPreviewTruncates(t *testing.T) {
	cfg := &config.SimpleConfig{
		Theme:     config.ThemeModeNerdFont,
		Separator: " | ",
		SegmentOrder: []string{
			"model",
			config.LineBreakMarker,
			"directory",
			config.LineBreakMarker,
			"git",
			config.LineBreakMarker,
			"context_window",
		},
		Segments: config.SegmentToggles{
			Model:         true,
			Directory:     true,
			Git:           true,
			ContextWindow: true,
			Usage:         true,
			Cost:          true,
			Session:       true,
			OutputStyle:   true,
			Update:        true,
			CCHModel:      true,
			CCHProvider:   true,
			CCHCost:       true,
			CCHRequests:   true,
			CCHLimits:     true,
		},
	}

	m := NewModel(cfg)
	m.width = 44
	m.height = 20
	m.debugKey = "up"

	s := m.View()
	lines := strings.Split(s, "\n")
	if len(lines) != m.height {
		t.Fatalf("expected %d lines, got %d", m.height, len(lines))
	}

	expectedMaxWidth := m.width
	if expectedMaxWidth > 1 {
		expectedMaxWidth--
	}
	for i, line := range lines {
		if w := ansi.StringWidth(line); w > expectedMaxWidth {
			t.Fatalf("line %d exceeds width: %d > %d", i, w, expectedMaxWidth)
		}
	}

	if !strings.Contains(s, "CCHLine Configuration") {
		t.Fatalf("expected title to be visible")
	}

	if strings.Contains(s, "Preview:") {
		t.Fatalf("did not expect Preview label line")
	}

	if !strings.Contains(s, "...") {
		t.Fatalf("expected preview to include ellipsis when truncated")
	}

	// Right border should be visible in at least one box line.
	// We check for common right border glyphs from RoundedBorder.
	rightBorderSeen := false
	for _, line := range lines {
		plain := ansi.Strip(line)
		if strings.Contains(plain, "│") || strings.Contains(plain, "╮") || strings.Contains(plain, "╯") {
			rightBorderSeen = true
			break
		}
	}
	if !rightBorderSeen {
		t.Fatalf("expected right border glyphs to be present in output")
	}
}
