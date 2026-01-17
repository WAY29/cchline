package tui

import (
	"strings"
	"testing"

	"github.com/WAY29/cchline/config"
	"github.com/charmbracelet/x/ansi"
)

func TestSegmentOrderRowsRoundTrip(t *testing.T) {
	cases := [][]string{
		{"model", "directory", config.LineBreakMarker, "context_window"},
		{"model", config.LineBreakMarker, config.LineBreakMarker, "context_window"},
		{},
	}

	for _, order := range cases {
		wantEnabled := make([]bool, config.NonBreakSegmentCount(order))
		for i := range wantEnabled {
			wantEnabled[i] = i%2 == 0
		}

		rows := segmentOrderToRows(order, wantEnabled)
		gotOrder, gotEnabled := rowsToSegmentOrder(rows)

		if len(order) == 0 {
			if len(gotOrder) != 0 || len(gotEnabled) != 0 {
				t.Fatalf("expected empty order to remain empty, got order=%#v enabled=%#v", gotOrder, gotEnabled)
			}
			continue
		}

		if len(gotOrder) != len(order) {
			t.Fatalf("expected order len %d, got %d (got=%#v, want=%#v)", len(order), len(gotOrder), gotOrder, order)
		}
		for i := range order {
			if gotOrder[i] != order[i] {
				t.Fatalf("mismatch at %d: got=%q want=%q (got=%#v, want=%#v)", i, gotOrder[i], order[i], gotOrder, order)
			}
		}

		if len(gotEnabled) != len(wantEnabled) {
			t.Fatalf("expected enabled len %d, got %d (got=%#v, want=%#v)", len(wantEnabled), len(gotEnabled), gotEnabled, wantEnabled)
		}
		for i := range wantEnabled {
			if gotEnabled[i] != wantEnabled[i] {
				t.Fatalf("enabled mismatch at %d: got=%v want=%v (got=%#v, want=%#v)", i, gotEnabled[i], wantEnabled[i], gotEnabled, wantEnabled)
			}
		}
	}
}

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
