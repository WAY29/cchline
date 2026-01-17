package tests

import (
	"reflect"
	"strings"
	"testing"

	"github.com/WAY29/cchline/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// TestThemeModeConstants verifies ThemeMode constant values
func TestThemeModeConstants(t *testing.T) {
	tests := []struct {
		name     string
		mode     config.ThemeMode
		expected string
	}{
		{
			name:     "ThemeModeDefault constant",
			mode:     config.ThemeModeDefault,
			expected: "default",
		},
		{
			name:     "ThemeModeNerdFont constant",
			mode:     config.ThemeModeNerdFont,
			expected: "nerd_font",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("got %q, want %q", string(tt.mode), tt.expected)
			}
		})
	}
}

// TestSegmentIDConstants verifies all SegmentID constants are defined
func TestSegmentIDConstants(t *testing.T) {
	expectedSegments := []config.SegmentID{
		config.SegmentModel,
		config.SegmentDirectory,
		config.SegmentGit,
		config.SegmentContextWindow,
		config.SegmentUsage,
		config.SegmentCost,
		config.SegmentSession,
		config.SegmentOutputStyle,
		config.SegmentUpdate,
		config.SegmentCCHModel,
		config.SegmentCCHProvider,
		config.SegmentCCHCost,
		config.SegmentCCHRequests,
		config.SegmentCCHLimits,
	}

	expectedValues := map[config.SegmentID]string{
		config.SegmentModel:         "model",
		config.SegmentDirectory:     "directory",
		config.SegmentGit:           "git",
		config.SegmentContextWindow: "context_window",
		config.SegmentUsage:         "usage",
		config.SegmentCost:          "cost",
		config.SegmentSession:       "session",
		config.SegmentOutputStyle:   "output_style",
		config.SegmentUpdate:        "update",
		config.SegmentCCHModel:      "cch_model",
		config.SegmentCCHProvider:   "cch_provider",
		config.SegmentCCHCost:       "cch_cost",
		config.SegmentCCHRequests:   "cch_requests",
		config.SegmentCCHLimits:     "cch_limits",
	}

	for _, segment := range expectedSegments {
		t.Run(string(segment), func(t *testing.T) {
			expected, ok := expectedValues[segment]
			if !ok {
				t.Errorf("segment %q not found in expected values", segment)
				return
			}
			if string(segment) != expected {
				t.Errorf("got %q, want %q", string(segment), expected)
			}
		})
	}
}

// TestGetSegmentThemeDefaultMode verifies GetSegmentTheme returns correct theme for default mode
func TestGetSegmentThemeDefaultMode(t *testing.T) {
	tests := []struct {
		name      string
		segmentID config.SegmentID
		hasIcon   bool
		hasColor  bool
	}{
		{
			name:      "Model segment default theme",
			segmentID: config.SegmentModel,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Directory segment default theme",
			segmentID: config.SegmentDirectory,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Git segment default theme",
			segmentID: config.SegmentGit,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "ContextWindow segment default theme",
			segmentID: config.SegmentContextWindow,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Usage segment default theme",
			segmentID: config.SegmentUsage,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Cost segment default theme",
			segmentID: config.SegmentCost,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Session segment default theme",
			segmentID: config.SegmentSession,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "OutputStyle segment default theme",
			segmentID: config.SegmentOutputStyle,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Update segment default theme",
			segmentID: config.SegmentUpdate,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Model segment default theme",
			segmentID: config.SegmentCCHModel,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Provider segment default theme",
			segmentID: config.SegmentCCHProvider,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Cost segment default theme",
			segmentID: config.SegmentCCHCost,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Requests segment default theme",
			segmentID: config.SegmentCCHRequests,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Limits segment default theme",
			segmentID: config.SegmentCCHLimits,
			hasIcon:   true,
			hasColor:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := config.GetSegmentTheme(tt.segmentID, config.ThemeModeDefault)

			if tt.hasIcon && theme.Icon == "" {
				t.Errorf("expected icon for segment %q, got empty string", tt.segmentID)
			}

			if tt.hasColor {
				if theme.IconColor == nil {
					t.Errorf("expected IconColor for segment %q, got nil", tt.segmentID)
				}
				if theme.TextColor == nil {
					t.Errorf("expected TextColor for segment %q, got nil", tt.segmentID)
				}
			}
		})
	}
}

// TestGetSegmentThemeNerdFontMode verifies GetSegmentTheme returns correct theme for nerd_font mode
func TestGetSegmentThemeNerdFontMode(t *testing.T) {
	tests := []struct {
		name      string
		segmentID config.SegmentID
		hasIcon   bool
		hasColor  bool
	}{
		{
			name:      "Model segment nerd font theme",
			segmentID: config.SegmentModel,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Directory segment nerd font theme",
			segmentID: config.SegmentDirectory,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Git segment nerd font theme",
			segmentID: config.SegmentGit,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "ContextWindow segment nerd font theme",
			segmentID: config.SegmentContextWindow,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Usage segment nerd font theme",
			segmentID: config.SegmentUsage,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Cost segment nerd font theme",
			segmentID: config.SegmentCost,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Session segment nerd font theme",
			segmentID: config.SegmentSession,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "OutputStyle segment nerd font theme",
			segmentID: config.SegmentOutputStyle,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "Update segment nerd font theme",
			segmentID: config.SegmentUpdate,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Model segment nerd font theme",
			segmentID: config.SegmentCCHModel,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Provider segment nerd font theme",
			segmentID: config.SegmentCCHProvider,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Cost segment nerd font theme",
			segmentID: config.SegmentCCHCost,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Requests segment nerd font theme",
			segmentID: config.SegmentCCHRequests,
			hasIcon:   true,
			hasColor:  true,
		},
		{
			name:      "CCH Limits segment nerd font theme",
			segmentID: config.SegmentCCHLimits,
			hasIcon:   true,
			hasColor:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := config.GetSegmentTheme(tt.segmentID, config.ThemeModeNerdFont)

			if tt.hasIcon && theme.Icon == "" {
				t.Errorf("expected icon for segment %q, got empty string", tt.segmentID)
			}

			if tt.hasColor {
				if theme.IconColor == nil {
					t.Errorf("expected IconColor for segment %q, got nil", tt.segmentID)
				}
				if theme.TextColor == nil {
					t.Errorf("expected TextColor for segment %q, got nil", tt.segmentID)
				}
			}
		})
	}
}

// TestGetSegmentThemeDifferentIcons verifies default and nerd_font modes return different icons
func TestGetSegmentThemeDifferentIcons(t *testing.T) {
	segments := []config.SegmentID{
		config.SegmentModel,
		config.SegmentDirectory,
		config.SegmentGit,
		config.SegmentContextWindow,
		config.SegmentUsage,
		config.SegmentCost,
		config.SegmentSession,
		config.SegmentOutputStyle,
		config.SegmentUpdate,
		config.SegmentCCHModel,
		config.SegmentCCHProvider,
		config.SegmentCCHCost,
		config.SegmentCCHRequests,
		config.SegmentCCHLimits,
	}

	for _, segment := range segments {
		t.Run(string(segment), func(t *testing.T) {
			defaultTheme := config.GetSegmentTheme(segment, config.ThemeModeDefault)
			nerdFontTheme := config.GetSegmentTheme(segment, config.ThemeModeNerdFont)

			if defaultTheme.Icon == nerdFontTheme.Icon {
				t.Errorf("expected different icons for segment %q, got same icon %q", segment, defaultTheme.Icon)
			}
		})
	}
}

// TestGetSegmentThemeSameColors verifies both modes return same colors
func TestGetSegmentThemeSameColors(t *testing.T) {
	segments := []config.SegmentID{
		config.SegmentModel,
		config.SegmentDirectory,
		config.SegmentGit,
		config.SegmentContextWindow,
		config.SegmentUsage,
		config.SegmentCost,
		config.SegmentSession,
		config.SegmentOutputStyle,
		config.SegmentUpdate,
		config.SegmentCCHModel,
		config.SegmentCCHProvider,
		config.SegmentCCHCost,
		config.SegmentCCHRequests,
		config.SegmentCCHLimits,
	}

	for _, segment := range segments {
		t.Run(string(segment), func(t *testing.T) {
			defaultTheme := config.GetSegmentTheme(segment, config.ThemeModeDefault)
			nerdFontTheme := config.GetSegmentTheme(segment, config.ThemeModeNerdFont)

			// Colors should be the same regardless of theme mode
			if defaultTheme.IconColor != nerdFontTheme.IconColor {
				t.Errorf("expected same IconColor for segment %q in both modes", segment)
			}
			if defaultTheme.TextColor != nerdFontTheme.TextColor {
				t.Errorf("expected same TextColor for segment %q in both modes", segment)
			}
			if defaultTheme.Bold != nerdFontTheme.Bold {
				t.Errorf("expected same Bold for segment %q in both modes", segment)
			}
		})
	}
}

// TestApplyColorWithValidColor verifies ApplyColor applies color to text
func TestApplyColorWithValidColor(t *testing.T) {
	testText := "test"
	testColor := lipgloss.Color("1")

	result := config.ApplyColor(testText, &testColor)

	// Result should contain ANSI escape codes
	if result == testText {
		t.Errorf("expected colored text with ANSI codes, got plain text")
	}

	// Result should contain the original text
	if !strings.Contains(result, testText) {
		t.Errorf("expected result to contain original text %q, got %q", testText, result)
	}
}

// TestApplyColorWithNilColor verifies ApplyColor returns plain text when color is nil
func TestApplyColorWithNilColor(t *testing.T) {
	testText := "test"

	result := config.ApplyColor(testText, nil)

	if result != testText {
		t.Errorf("expected plain text %q, got %q", testText, result)
	}
}

// TestApplyColorWithEmptyString verifies ApplyColor handles empty string
func TestApplyColorWithEmptyString(t *testing.T) {
	testColor := lipgloss.Color("2")

	result := config.ApplyColor("", &testColor)

	if ansi.Strip(result) != "" {
		t.Errorf("expected empty string after stripping ANSI, got %q", ansi.Strip(result))
	}
}

// TestDefaultSegmentOrder verifies DefaultSegmentOrder matches the default layout
func TestDefaultSegmentOrder(t *testing.T) {
	want := []string{"model", "directory", "output_style", config.LineBreakMarker, "context_window"}
	if !reflect.DeepEqual(config.DefaultSegmentOrder, want) {
		t.Errorf("DefaultSegmentOrder mismatch: got=%#v want=%#v", config.DefaultSegmentOrder, want)
	}
}

// TestDefaultSegmentOrderNoDuplicates verifies DefaultSegmentOrder has no duplicates
func TestDefaultSegmentOrderNoDuplicates(t *testing.T) {
	seen := make(map[string]bool)
	for _, segment := range config.DefaultSegmentOrder {
		if seen[segment] {
			t.Errorf("duplicate segment %q found in DefaultSegmentOrder", segment)
		}
		seen[segment] = true
	}
}

// TestGetSegmentThemeInvalidSegment verifies GetSegmentTheme returns empty theme for invalid segment
func TestGetSegmentThemeInvalidSegment(t *testing.T) {
	invalidSegment := config.SegmentID("invalid_segment")

	theme := config.GetSegmentTheme(invalidSegment, config.ThemeModeDefault)

	if theme.Icon != "" {
		t.Errorf("expected empty icon for invalid segment, got %q", theme.Icon)
	}
	if theme.IconColor != nil {
		t.Errorf("expected nil IconColor for invalid segment, got %v", theme.IconColor)
	}
	if theme.TextColor != nil {
		t.Errorf("expected nil TextColor for invalid segment, got %v", theme.TextColor)
	}
}

// TestSegmentThemeStructure verifies SegmentTheme has all required fields
func TestSegmentThemeStructure(t *testing.T) {
	theme := config.GetSegmentTheme(config.SegmentModel, config.ThemeModeDefault)

	// Verify all fields are accessible
	_ = theme.Icon
	_ = theme.IconColor
	_ = theme.TextColor
	_ = theme.BgColor
	_ = theme.Bold

	// Verify Icon is not empty for valid segment
	if theme.Icon == "" {
		t.Errorf("expected non-empty Icon for valid segment")
	}
}
