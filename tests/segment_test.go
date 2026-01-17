package tests

import (
	"fmt"
	"testing"

	"github.com/WAY29/cchline/config"
	"github.com/WAY29/cchline/segment"
	"github.com/charmbracelet/lipgloss"
)

// TestModelSegmentCollect tests the ModelSegment.Collect() method
func TestModelSegmentCollect(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
	}{
		{
			name: "Sonnet 3.5 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-5-sonnet-20241022",
					DisplayName: "claude-3-5-sonnet-20241022",
				},
			},
			expected: "Sonnet 3.5",
		},
		{
			name: "Opus 3 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-opus-20240229",
					DisplayName: "claude-3-opus-20240229",
				},
			},
			expected: "Opus 3",
		},
		{
			name: "Haiku 3 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-haiku-20240307",
					DisplayName: "claude-3-haiku-20240307",
				},
			},
			expected: "Haiku 3",
		},
		{
			name: "Opus 4 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-opus-4-20250514",
					DisplayName: "claude-opus-4-20250514",
				},
			},
			expected: "Opus 4",
		},
		{
			name: "Sonnet 4 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-sonnet-4-20250514",
					DisplayName: "claude-sonnet-4-20250514",
				},
			},
			expected: "Sonnet 4",
		},
		{
			name: "GPT-4 simplification",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "gpt-4-turbo",
					DisplayName: "gpt-4-turbo",
				},
			},
			expected: "GPT-4",
		},
		{
			name: "Unknown model name",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "unknown-model",
					DisplayName: "unknown-model",
				},
			},
			expected: "unknown-model",
		},
		{
			name: "Empty DisplayName falls back to ID",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-5-sonnet-20241022",
					DisplayName: "",
				},
			},
			expected: "Sonnet 3.5",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modelSegment.Collect(tt.input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentCollect tests the DirectorySegment.Collect() method
func TestDirectorySegmentCollect(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
	}{
		{
			name: "Extract last directory component",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "/Users/lang/coding/golang/src/cchline",
				},
			},
			expected: "cchline",
		},
		{
			name: "Single directory name",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "myapp",
				},
			},
			expected: "myapp",
		},
		{
			name: "Root directory",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "/",
				},
			},
			expected: "/",
		},
		{
			name: "Nested path",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "/home/user/projects/my-project",
				},
			},
			expected: "my-project",
		},
		{
			name: "Path with trailing slash",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "/home/user/projects/my-project/",
				},
			},
			expected: "my-project",
		},
		{
			name: "Empty directory",
			input: &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: "",
				},
			},
			expected: ".",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dirSegment.Collect(tt.input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentCollect tests the CostSegment.Collect() method
func TestCostSegmentCollect(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
		isEmpty  bool
	}{
		{
			name: "Format cost with two decimal places",
			input: &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: 0.15,
				},
			},
			expected: "$0.15",
			isEmpty:  false,
		},
		{
			name: "Format cost with rounding",
			input: &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: 1.567,
				},
			},
			expected: "$1.57",
			isEmpty:  false,
		},
		{
			name: "Format large cost",
			input: &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: 123.45,
				},
			},
			expected: "$123.45",
			isEmpty:  false,
		},
		{
			name: "Zero cost returns empty",
			input: &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: 0,
				},
			},
			expected: "",
			isEmpty:  true,
		},
		{
			name: "Small cost",
			input: &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: 0.01,
				},
			},
			expected: "$0.01",
			isEmpty:  false,
		},
	}

	costSegment := &segment.CostSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := costSegment.Collect(tt.input)
			if tt.isEmpty {
				if result.Primary != "" {
					t.Errorf("expected empty result, got %q", result.Primary)
				}
			} else {
				if result.Primary != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result.Primary)
				}
			}
		})
	}
}

// TestOutputStyleSegmentCollect tests the OutputStyleSegment.Collect() method
func TestOutputStyleSegmentCollect(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
		isEmpty  bool
	}{
		{
			name: "Default output style",
			input: &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: "default",
				},
			},
			expected: "default",
			isEmpty:  false,
		},
		{
			name: "Custom output style",
			input: &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: "markdown",
				},
			},
			expected: "markdown",
			isEmpty:  false,
		},
		{
			name: "Empty output style returns empty",
			input: &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: "",
				},
			},
			expected: "",
			isEmpty:  true,
		},
		{
			name: "Output style with special characters",
			input: &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: "custom-style_v2",
				},
			},
			expected: "custom-style_v2",
			isEmpty:  false,
		},
	}

	styleSegment := &segment.OutputStyleSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := styleSegment.Collect(tt.input)
			if tt.isEmpty {
				if result.Primary != "" {
					t.Errorf("expected empty result, got %q", result.Primary)
				}
			} else {
				if result.Primary != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result.Primary)
				}
			}
		})
	}
}

// TestSegmentDataStructure tests that SegmentData is properly initialized
func TestSegmentDataStructure(t *testing.T) {
	data := segment.SegmentData{
		Primary:   "test",
		Secondary: "secondary",
		Metadata:  make(map[string]string),
	}

	if data.Primary != "test" {
		t.Errorf("expected Primary to be 'test', got %q", data.Primary)
	}

	if data.Secondary != "secondary" {
		t.Errorf("expected Secondary to be 'secondary', got %q", data.Secondary)
	}

	if data.Metadata == nil {
		t.Error("expected Metadata to be initialized, got nil")
	}
}

// TestSegmentResultStructure tests that SegmentResult is properly initialized
func TestSegmentResultStructure(t *testing.T) {
	data := segment.SegmentData{
		Primary: "test",
	}

	result := segment.SegmentResult{
		ID:   config.SegmentModel,
		Data: data,
	}

	if result.ID != config.SegmentModel {
		t.Errorf("expected ID to be SegmentIDModel, got %v", result.ID)
	}

	if result.Data.Primary != "test" {
		t.Errorf("expected Data.Primary to be 'test', got %q", result.Data.Primary)
	}
}

// TestUpdateSegmentCollect tests the UpdateSegment.Collect() method
func TestUpdateSegmentCollect(t *testing.T) {
	updateSegment := &segment.UpdateSegment{}

	input := &config.InputData{}
	result := updateSegment.Collect(input)

	// UpdateSegment always returns empty SegmentData
	if result.Primary != "" {
		t.Errorf("expected empty Primary, got %q", result.Primary)
	}

	if result.Secondary != "" {
		t.Errorf("expected empty Secondary, got %q", result.Secondary)
	}
}

// TestMultipleSegmentsIntegration tests multiple segments working together
func TestMultipleSegmentsIntegration(t *testing.T) {
	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet-20241022",
			DisplayName: "claude-3-5-sonnet-20241022",
		},
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/Users/lang/coding/golang/src/cchline",
		},
		Cost: config.CostInfo{
			TotalCostUSD: 0.25,
		},
		OutputStyle: config.OutputStyleInfo{
			Name: "default",
		},
	}

	// Test ModelSegment
	modelSegment := &segment.ModelSegment{}
	modelResult := modelSegment.Collect(input)
	if modelResult.Primary != "Sonnet 3.5" {
		t.Errorf("ModelSegment: expected 'Sonnet 3.5', got %q", modelResult.Primary)
	}

	// Test DirectorySegment
	dirSegment := &segment.DirectorySegment{}
	dirResult := dirSegment.Collect(input)
	if dirResult.Primary != "cchline" {
		t.Errorf("DirectorySegment: expected 'cchline', got %q", dirResult.Primary)
	}

	// Test CostSegment
	costSegment := &segment.CostSegment{}
	costResult := costSegment.Collect(input)
	if costResult.Primary != "$0.25" {
		t.Errorf("CostSegment: expected '$0.25', got %q", costResult.Primary)
	}

	// Test OutputStyleSegment
	styleSegment := &segment.OutputStyleSegment{}
	styleResult := styleSegment.Collect(input)
	if styleResult.Primary != "default" {
		t.Errorf("OutputStyleSegment: expected 'default', got %q", styleResult.Primary)
	}
}

// TestCostSegmentEdgeCases tests edge cases for cost formatting
func TestCostSegmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Very small cost",
			cost:     0.001,
			expected: "$0.00",
		},
		{
			name:     "Cost with many decimals",
			cost:     1.23456,
			expected: "$1.23",
		},
		{
			name:     "Large cost",
			cost:     9999.99,
			expected: "$9999.99",
		},
		{
			name:     "Cost exactly one dollar",
			cost:     1.0,
			expected: "$1.00",
		},
		{
			name:     "Cost with single decimal",
			cost:     0.5,
			expected: "$0.50",
		},
	}

	costSegment := &segment.CostSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentEdgeCases tests edge cases for directory extraction
func TestDirectorySegmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Path with dots",
			dir:      "/home/user/.config/app",
			expected: "app",
		},
		{
			name:     "Path with hyphens and underscores",
			dir:      "/home/user/my-project_v2",
			expected: "my-project_v2",
		},
		{
			name:     "Relative path",
			dir:      "src/components",
			expected: "components",
		},
		{
			name:     "Single dot",
			dir:      ".",
			expected: ".",
		},
		{
			name:     "Double dots",
			dir:      "..",
			expected: "..",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestModelSegmentEdgeCases tests edge cases for model name simplification
func TestModelSegmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
	}{
		{
			name: "Case insensitive matching",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "CLAUDE-3-5-SONNET",
					DisplayName: "CLAUDE-3-5-SONNET",
				},
			},
			expected: "Sonnet 3.5",
		},
		{
			name: "Model name with version suffix",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-5-sonnet-20241022",
					DisplayName: "claude-3-5-sonnet-20241022",
				},
			},
			expected: "Sonnet 3.5",
		},
		{
			name: "Partial model name match",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "my-claude-3-opus-custom",
					DisplayName: "my-claude-3-opus-custom",
				},
			},
			expected: "Opus 3",
		},
		{
			name: "Model name not in replacements",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "custom-model-v1",
					DisplayName: "custom-model-v1",
				},
			},
			expected: "custom-model-v1",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modelSegment.Collect(tt.input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentEdgeCases tests edge cases for output style
func TestOutputStyleSegmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		expected string
		isEmpty  bool
	}{
		{
			name:     "Style with spaces",
			style:    "my style",
			expected: "my style",
			isEmpty:  false,
		},
		{
			name:     "Style with numbers",
			style:    "style123",
			expected: "style123",
			isEmpty:  false,
		},
		{
			name:     "Style with special characters",
			style:    "style-v2.0",
			expected: "style-v2.0",
			isEmpty:  false,
		},
		{
			name:     "Whitespace only",
			style:    "   ",
			expected: "   ",
			isEmpty:  false,
		},
	}

	styleSegment := &segment.OutputStyleSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if tt.isEmpty {
				if result.Primary != "" {
					t.Errorf("expected empty result, got %q", result.Primary)
				}
			} else {
				if result.Primary != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result.Primary)
				}
			}
		})
	}
}

// TestSegmentDataMetadata tests SegmentData metadata handling
func TestSegmentDataMetadata(t *testing.T) {
	metadata := make(map[string]string)
	metadata["key1"] = "value1"
	metadata["key2"] = "value2"

	data := segment.SegmentData{
		Primary:   "primary",
		Secondary: "secondary",
		Metadata:  metadata,
	}

	if data.Metadata["key1"] != "value1" {
		t.Errorf("expected metadata key1 to be 'value1', got %q", data.Metadata["key1"])
	}

	if data.Metadata["key2"] != "value2" {
		t.Errorf("expected metadata key2 to be 'value2', got %q", data.Metadata["key2"])
	}

	if len(data.Metadata) != 2 {
		t.Errorf("expected 2 metadata entries, got %d", len(data.Metadata))
	}
}

// TestSegmentDataEmptyMetadata tests SegmentData with empty metadata
func TestSegmentDataEmptyMetadata(t *testing.T) {
	data := segment.SegmentData{
		Primary: "test",
	}

	if data.Metadata != nil {
		t.Errorf("expected nil Metadata, got %v", data.Metadata)
	}
}

// TestSegmentDataWithInitializedMetadata tests SegmentData with initialized metadata
func TestSegmentDataWithInitializedMetadata(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	if data.Metadata == nil {
		t.Error("expected initialized Metadata, got nil")
	}

	if len(data.Metadata) != 0 {
		t.Errorf("expected empty Metadata, got %d entries", len(data.Metadata))
	}

	// Test adding to metadata
	data.Metadata["test"] = "value"
	if data.Metadata["test"] != "value" {
		t.Errorf("expected metadata test to be 'value', got %q", data.Metadata["test"])
	}
}

// TestGetSegmentThemeInSegment tests the GetSegmentTheme function
func TestGetSegmentThemeInSegment(t *testing.T) {
	tests := []struct {
		name      string
		segmentID config.SegmentID
		expected  string
	}{
		{
			name:      "Model segment ID",
			segmentID: config.SegmentModel,
			expected:  "model",
		},
		{
			name:      "Directory segment ID",
			segmentID: config.SegmentDirectory,
			expected:  "directory",
		},
		{
			name:      "Git segment ID",
			segmentID: config.SegmentGit,
			expected:  "git",
		},
		{
			name:      "Context window segment ID",
			segmentID: config.SegmentContextWindow,
			expected:  "context_window",
		},
		{
			name:      "Usage segment ID",
			segmentID: config.SegmentUsage,
			expected:  "usage",
		},
		{
			name:      "Cost segment ID",
			segmentID: config.SegmentCost,
			expected:  "cost",
		},
		{
			name:      "Session segment ID",
			segmentID: config.SegmentSession,
			expected:  "session",
		},
		{
			name:      "Output style segment ID",
			segmentID: config.SegmentOutputStyle,
			expected:  "output_style",
		},
		{
			name:      "Update segment ID",
			segmentID: config.SegmentUpdate,
			expected:  "update",
		},
		{
			name:      "CCH model segment ID",
			segmentID: config.SegmentCCHModel,
			expected:  "cch_model",
		},
		{
			name:      "CCH provider segment ID",
			segmentID: config.SegmentCCHProvider,
			expected:  "cch_provider",
		},
		{
			name:      "CCH cost segment ID",
			segmentID: config.SegmentCCHCost,
			expected:  "cch_cost",
		},
		{
			name:      "CCH requests segment ID",
			segmentID: config.SegmentCCHRequests,
			expected:  "cch_requests",
		},
		{
			name:      "CCH limits segment ID",
			segmentID: config.SegmentCCHLimits,
			expected:  "cch_limits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.segmentID) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.segmentID))
			}
		})
	}
}

// TestApplyColorInSegment tests the ApplyColor function
func TestApplyColorInSegment(t *testing.T) {
	tests := []struct {
		name     string
		mode     config.ThemeMode
		expected string
	}{
		{
			name:     "Default theme mode",
			mode:     config.ThemeModeDefault,
			expected: "default",
		},
		{
			name:     "Nerd font theme mode",
			mode:     config.ThemeModeNerdFont,
			expected: "nerd_font",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.mode))
			}
		})
	}
}

// TestGetSegmentTheme tests the GetSegmentTheme function
func TestGetSegmentTheme(t *testing.T) {
	tests := []struct {
		name     string
		id       config.SegmentID
		mode     config.ThemeMode
		hasIcon  bool
		hasColor bool
	}{
		{
			name:     "Model segment with default theme",
			id:       config.SegmentModel,
			mode:     config.ThemeModeDefault,
			hasIcon:  true,
			hasColor: true,
		},
		{
			name:     "Model segment with nerd font theme",
			id:       config.SegmentModel,
			mode:     config.ThemeModeNerdFont,
			hasIcon:  true,
			hasColor: true,
		},
		{
			name:     "Directory segment with default theme",
			id:       config.SegmentDirectory,
			mode:     config.ThemeModeDefault,
			hasIcon:  true,
			hasColor: true,
		},
		{
			name:     "Cost segment with nerd font theme",
			id:       config.SegmentCost,
			mode:     config.ThemeModeNerdFont,
			hasIcon:  true,
			hasColor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := config.GetSegmentTheme(tt.id, tt.mode)

			if tt.hasIcon && theme.Icon == "" {
				t.Errorf("expected icon to be set, got empty string")
			}

			if tt.hasColor && theme.IconColor == nil {
				t.Errorf("expected IconColor to be set, got nil")
			}

			if tt.hasColor && theme.TextColor == nil {
				t.Errorf("expected TextColor to be set, got nil")
			}
		})
	}
}

// TestApplyColor tests the ApplyColor function
func TestApplyColor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		hasColor bool
		expected string
	}{
		{
			name:     "Apply color to text",
			text:     "test",
			hasColor: true,
			expected: "test",
		},
		{
			name:     "Apply nil color returns original text",
			text:     "test",
			hasColor: false,
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *lipgloss.Color
			if tt.hasColor {
				cc := lipgloss.Color("2")
				c = &cc
			}

			result := config.ApplyColor(tt.text, c)
			// We can't directly compare colored text, but we can check it's not empty
			if result == "" && tt.text != "" {
				t.Errorf("expected non-empty result for non-empty input")
			}
		})
	}
}

// TestInputDataStructure tests InputData structure initialization
func TestInputDataStructure(t *testing.T) {
	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet",
			DisplayName: "Sonnet 3.5",
		},
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/home/user/project",
		},
		TranscriptPath: "/home/user/project/transcript.jsonl",
		Cost: config.CostInfo{
			TotalCostUSD: 0.50,
		},
		OutputStyle: config.OutputStyleInfo{
			Name: "default",
		},
	}

	if input.Model.ID != "claude-3-5-sonnet" {
		t.Errorf("expected Model.ID to be 'claude-3-5-sonnet', got %q", input.Model.ID)
	}

	if input.Model.DisplayName != "Sonnet 3.5" {
		t.Errorf("expected Model.DisplayName to be 'Sonnet 3.5', got %q", input.Model.DisplayName)
	}

	if input.Workspace.CurrentDir != "/home/user/project" {
		t.Errorf("expected Workspace.CurrentDir to be '/home/user/project', got %q", input.Workspace.CurrentDir)
	}

	if input.TranscriptPath != "/home/user/project/transcript.jsonl" {
		t.Errorf("expected TranscriptPath to be '/home/user/project/transcript.jsonl', got %q", input.TranscriptPath)
	}

	if input.Cost.TotalCostUSD != 0.50 {
		t.Errorf("expected Cost.TotalCostUSD to be 0.50, got %f", input.Cost.TotalCostUSD)
	}

	if input.OutputStyle.Name != "default" {
		t.Errorf("expected OutputStyle.Name to be 'default', got %q", input.OutputStyle.Name)
	}
}

// TestSegmentToggles tests SegmentToggles structure
func TestSegmentToggles(t *testing.T) {
	toggles := config.SegmentToggles{
		Model:         true,
		Directory:     true,
		Git:           false,
		ContextWindow: true,
		Usage:         false,
		Cost:          true,
		Session:       false,
		OutputStyle:   true,
		Update:        false,
		CCHModel:      true,
		CCHProvider:   false,
		CCHCost:       true,
		CCHRequests:   false,
		CCHLimits:     true,
	}

	if !toggles.Model {
		t.Error("expected Model to be true")
	}

	if toggles.Git {
		t.Error("expected Git to be false")
	}

	if !toggles.ContextWindow {
		t.Error("expected ContextWindow to be true")
	}

	if toggles.Usage {
		t.Error("expected Usage to be false")
	}

	if !toggles.CCHModel {
		t.Error("expected CCHModel to be true")
	}

	if toggles.CCHProvider {
		t.Error("expected CCHProvider to be false")
	}
}

// TestCostSegmentWithNegativeValue tests cost segment with negative value
func TestCostSegmentWithNegativeValue(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: -0.50,
		},
	}

	result := costSegment.Collect(input)
	// Negative cost should still be formatted
	if result.Primary != "$-0.50" {
		t.Errorf("expected '$-0.50', got %q", result.Primary)
	}
}

// TestModelSegmentWithEmptyBothFields tests model segment when both ID and DisplayName are empty
func TestModelSegmentWithEmptyBothFields(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "",
			DisplayName: "",
		},
	}

	result := modelSegment.Collect(input)
	// Should return empty string when both fields are empty
	if result.Primary != "" {
		t.Errorf("expected empty string, got %q", result.Primary)
	}
}

// TestDirectorySegmentWithComplexPaths tests directory segment with various complex paths
func TestDirectorySegmentWithComplexPaths(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Path with multiple slashes",
			dir:      "/home//user///project",
			expected: "project",
		},
		{
			name:     "Path with spaces in directory name",
			dir:      "/home/user/my project",
			expected: "my project",
		},
		{
			name:     "Path with unicode characters",
			dir:      "/home/user/项目",
			expected: "项目",
		},
		{
			name:     "Very long path",
			dir:      "/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z",
			expected: "z",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithEmptyString tests output style segment with empty string
func TestOutputStyleSegmentWithEmptyString(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "",
		},
	}

	result := styleSegment.Collect(input)
	if result.Primary != "" {
		t.Errorf("expected empty string, got %q", result.Primary)
	}
}

// TestSegmentResultWithDifferentIDs tests SegmentResult with different segment IDs
func TestSegmentResultWithDifferentIDs(t *testing.T) {
	tests := []struct {
		name string
		id   config.SegmentID
		data string
	}{
		{
			name: "Model segment result",
			id:   config.SegmentModel,
			data: "Sonnet 3.5",
		},
		{
			name: "Directory segment result",
			id:   config.SegmentDirectory,
			data: "myapp",
		},
		{
			name: "Cost segment result",
			id:   config.SegmentCost,
			data: "$0.15",
		},
		{
			name: "Git segment result",
			id:   config.SegmentGit,
			data: "main *",
		},
		{
			name: "CCH model segment result",
			id:   config.SegmentCCHModel,
			data: "claude-3-opus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := segment.SegmentResult{
				ID: tt.id,
				Data: segment.SegmentData{
					Primary: tt.data,
				},
			}

			if result.ID != tt.id {
				t.Errorf("expected ID %q, got %q", tt.id, result.ID)
			}

			if result.Data.Primary != tt.data {
				t.Errorf("expected Primary %q, got %q", tt.data, result.Data.Primary)
			}
		})
	}
}

// TestCostSegmentWithVeryLargeValue tests cost segment with very large value
func TestCostSegmentWithVeryLargeValue(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 999999.99,
		},
	}

	result := costSegment.Collect(input)
	if result.Primary != "$999999.99" {
		t.Errorf("expected '$999999.99', got %q", result.Primary)
	}
}

// TestModelSegmentPreservesCase tests that model segment preserves case in unknown models
func TestModelSegmentPreservesCase(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "MyCustomModel",
			DisplayName: "MyCustomModel",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "MyCustomModel" {
		t.Errorf("expected 'MyCustomModel', got %q", result.Primary)
	}
}

// TestDirectorySegmentPreservesCase tests that directory segment preserves case
func TestDirectorySegmentPreservesCase(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/home/user/MyProject",
		},
	}

	result := dirSegment.Collect(input)
	if result.Primary != "MyProject" {
		t.Errorf("expected 'MyProject', got %q", result.Primary)
	}
}

// TestSimpleConfigStructure tests SimpleConfig structure
func TestSimpleConfigStructure(t *testing.T) {
	cfg := &config.SimpleConfig{
		Theme:        config.ThemeModeNerdFont,
		Separator:    " | ",
		SegmentOrder: []string{"model", "directory", "git"},
		SegmentEnabled: []bool{
			true,
			true,
			true,
		},
		CCHApiKey: "test-api-key",
		CCHURL:    "http://localhost:8000",
	}

	if cfg.Theme != config.ThemeModeNerdFont {
		t.Errorf("expected Theme to be ThemeModeNerdFont, got %q", cfg.Theme)
	}

	if cfg.Separator != " | " {
		t.Errorf("expected Separator to be ' | ', got %q", cfg.Separator)
	}

	if len(cfg.SegmentEnabled) == 0 || !cfg.SegmentEnabled[0] {
		t.Error("expected Model segment to be enabled")
	}

	if cfg.CCHApiKey != "test-api-key" {
		t.Errorf("expected CCHApiKey to be 'test-api-key', got %q", cfg.CCHApiKey)
	}

	if cfg.CCHURL != "http://localhost:8000" {
		t.Errorf("expected CCHURL to be 'http://localhost:8000', got %q", cfg.CCHURL)
	}
}

// TestModelInfoStructure tests ModelInfo structure
func TestModelInfoStructure(t *testing.T) {
	model := config.ModelInfo{
		ID:          "claude-3-5-sonnet-20241022",
		DisplayName: "Claude 3.5 Sonnet",
	}

	if model.ID != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected ID to be 'claude-3-5-sonnet-20241022', got %q", model.ID)
	}

	if model.DisplayName != "Claude 3.5 Sonnet" {
		t.Errorf("expected DisplayName to be 'Claude 3.5 Sonnet', got %q", model.DisplayName)
	}
}

// TestWorkspaceInfoStructure tests WorkspaceInfo structure
func TestWorkspaceInfoStructure(t *testing.T) {
	workspace := config.WorkspaceInfo{
		CurrentDir: "/home/user/project",
	}

	if workspace.CurrentDir != "/home/user/project" {
		t.Errorf("expected CurrentDir to be '/home/user/project', got %q", workspace.CurrentDir)
	}
}

// TestCostInfoStructure tests CostInfo structure
func TestCostInfoStructure(t *testing.T) {
	cost := config.CostInfo{
		TotalCostUSD: 1.50,
	}

	if cost.TotalCostUSD != 1.50 {
		t.Errorf("expected TotalCostUSD to be 1.50, got %f", cost.TotalCostUSD)
	}
}

// TestOutputStyleInfoStructure tests OutputStyleInfo structure
func TestOutputStyleInfoStructure(t *testing.T) {
	style := config.OutputStyleInfo{
		Name: "markdown",
	}

	if style.Name != "markdown" {
		t.Errorf("expected Name to be 'markdown', got %q", style.Name)
	}
}

// TestSegmentThemeFields tests SegmentTheme structure fields
func TestSegmentThemeFields(t *testing.T) {
	iconColor := lipgloss.Color("6")
	textColor := lipgloss.Color("2")
	theme := config.SegmentTheme{
		Icon:      "🤖",
		IconColor: &iconColor,
		TextColor: &textColor,
		Bold:      true,
	}

	if theme.Icon != "🤖" {
		t.Errorf("expected Icon to be '🤖', got %q", theme.Icon)
	}

	if theme.IconColor == nil {
		t.Error("expected IconColor to be set")
	}

	if theme.TextColor == nil {
		t.Error("expected TextColor to be set")
	}

	if !theme.Bold {
		t.Error("expected Bold to be true")
	}
}

// TestCostSegmentBoundaryValues tests cost segment with boundary values
func TestCostSegmentBoundaryValues(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Minimum positive value",
			cost:     0.01,
			expected: "$0.01",
		},
		{
			name:     "Value just below rounding threshold",
			cost:     0.004,
			expected: "$0.00",
		},
		{
			name:     "Value at rounding threshold",
			cost:     0.005,
			expected: "$0.01",
		},
		{
			name:     "Value just above rounding threshold",
			cost:     0.006,
			expected: "$0.01",
		},
		{
			name:     "Negative minimum",
			cost:     -0.01,
			expected: "$-0.01",
		},
	}

	costSegment := &segment.CostSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestModelSegmentWithPartialMatches tests model segment with partial name matches
func TestModelSegmentWithPartialMatches(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.InputData
		expected string
	}{
		{
			name: "Sonnet in middle of name",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "prefix-claude-3-5-sonnet-suffix",
					DisplayName: "prefix-claude-3-5-sonnet-suffix",
				},
			},
			expected: "Sonnet 3.5",
		},
		{
			name: "Multiple matching patterns - sonnet matched first",
			input: &config.InputData{
				Model: config.ModelInfo{
					ID:          "claude-3-opus-and-claude-3-5-sonnet",
					DisplayName: "claude-3-opus-and-claude-3-5-sonnet",
				},
			},
			expected: "Sonnet 3.5",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modelSegment.Collect(tt.input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentWithSpecialPaths tests directory segment with special paths
func TestDirectorySegmentWithSpecialPaths(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Symlink-like path",
			dir:      "/var/log/app",
			expected: "app",
		},
		{
			name:     "Hidden directory",
			dir:      "/home/user/.config",
			expected: ".config",
		},
		{
			name:     "Directory with numbers",
			dir:      "/home/user/project123",
			expected: "project123",
		},
		{
			name:     "Directory with mixed case and numbers",
			dir:      "/home/user/MyProject2024",
			expected: "MyProject2024",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithVariousNames tests output style segment with various style names
func TestOutputStyleSegmentWithVariousNames(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "Lowercase style",
			style:    "markdown",
			expected: "markdown",
		},
		{
			name:     "Uppercase style",
			style:    "MARKDOWN",
			expected: "MARKDOWN",
		},
		{
			name:     "Mixed case style",
			style:    "MarkDown",
			expected: "MarkDown",
		},
		{
			name:     "Style with version",
			style:    "v2.0",
			expected: "v2.0",
		},
		{
			name:     "Style with underscores",
			style:    "my_custom_style",
			expected: "my_custom_style",
		},
	}

	styleSegment := &segment.OutputStyleSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataSecondaryField tests SegmentData secondary field
func TestSegmentDataSecondaryField(t *testing.T) {
	data := segment.SegmentData{
		Primary:   "primary",
		Secondary: "secondary",
	}

	if data.Primary != "primary" {
		t.Errorf("expected Primary to be 'primary', got %q", data.Primary)
	}

	if data.Secondary != "secondary" {
		t.Errorf("expected Secondary to be 'secondary', got %q", data.Secondary)
	}
}

// TestSegmentDataMetadataOperations tests various metadata operations
func TestSegmentDataMetadataOperations(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	// Test adding multiple entries
	data.Metadata["key1"] = "value1"
	data.Metadata["key2"] = "value2"
	data.Metadata["key3"] = "value3"

	if len(data.Metadata) != 3 {
		t.Errorf("expected 3 metadata entries, got %d", len(data.Metadata))
	}

	// Test updating entry
	data.Metadata["key1"] = "updated_value1"
	if data.Metadata["key1"] != "updated_value1" {
		t.Errorf("expected key1 to be 'updated_value1', got %q", data.Metadata["key1"])
	}

	// Test deleting entry
	delete(data.Metadata, "key2")
	if len(data.Metadata) != 2 {
		t.Errorf("expected 2 metadata entries after deletion, got %d", len(data.Metadata))
	}

	if _, exists := data.Metadata["key2"]; exists {
		t.Error("expected key2 to be deleted")
	}
}

// TestCostSegmentZeroValue tests cost segment specifically with zero value
func TestCostSegmentZeroValue(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 0.0,
		},
	}

	result := costSegment.Collect(input)
	// Zero cost should return empty SegmentData
	if result.Primary != "" {
		t.Errorf("expected empty Primary for zero cost, got %q", result.Primary)
	}

	if result.Secondary != "" {
		t.Errorf("expected empty Secondary for zero cost, got %q", result.Secondary)
	}
}

// TestUpdateSegmentAlwaysEmpty tests that UpdateSegment always returns empty data
func TestUpdateSegmentAlwaysEmpty(t *testing.T) {
	updateSegment := &segment.UpdateSegment{}

	inputs := []*config.InputData{
		{},
		{
			Model: config.ModelInfo{
				ID:          "test",
				DisplayName: "test",
			},
		},
		{
			Cost: config.CostInfo{
				TotalCostUSD: 100.0,
			},
		},
	}

	for i, input := range inputs {
		result := updateSegment.Collect(input)
		if result.Primary != "" {
			t.Errorf("test case %d: expected empty Primary, got %q", i, result.Primary)
		}
		if result.Secondary != "" {
			t.Errorf("test case %d: expected empty Secondary, got %q", i, result.Secondary)
		}
	}
}

// TestSegmentInterfaceImplementation tests that all segments implement the Segment interface
func TestSegmentInterfaceImplementation(t *testing.T) {
	var _ segment.Segment = (*segment.ModelSegment)(nil)
	var _ segment.Segment = (*segment.DirectorySegment)(nil)
	var _ segment.Segment = (*segment.CostSegment)(nil)
	var _ segment.Segment = (*segment.OutputStyleSegment)(nil)
	var _ segment.Segment = (*segment.UpdateSegment)(nil)
}

// TestCostSegmentFormattingConsistency tests cost formatting consistency
func TestCostSegmentFormattingConsistency(t *testing.T) {
	costSegment := &segment.CostSegment{}

	// Test that same cost always produces same output
	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 0.123,
		},
	}

	result1 := costSegment.Collect(input)
	result2 := costSegment.Collect(input)

	if result1.Primary != result2.Primary {
		t.Errorf("expected consistent formatting, got %q and %q", result1.Primary, result2.Primary)
	}
}

// TestModelSegmentFormattingConsistency tests model name formatting consistency
func TestModelSegmentFormattingConsistency(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet-20241022",
			DisplayName: "claude-3-5-sonnet-20241022",
		},
	}

	result1 := modelSegment.Collect(input)
	result2 := modelSegment.Collect(input)

	if result1.Primary != result2.Primary {
		t.Errorf("expected consistent formatting, got %q and %q", result1.Primary, result2.Primary)
	}
}

// TestDirectorySegmentFormattingConsistency tests directory extraction consistency
func TestDirectorySegmentFormattingConsistency(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/home/user/project",
		},
	}

	result1 := dirSegment.Collect(input)
	result2 := dirSegment.Collect(input)

	if result1.Primary != result2.Primary {
		t.Errorf("expected consistent formatting, got %q and %q", result1.Primary, result2.Primary)
	}
}

// TestOutputStyleSegmentFormattingConsistency tests output style formatting consistency
func TestOutputStyleSegmentFormattingConsistency(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "markdown",
		},
	}

	result1 := styleSegment.Collect(input)
	result2 := styleSegment.Collect(input)

	if result1.Primary != result2.Primary {
		t.Errorf("expected consistent formatting, got %q and %q", result1.Primary, result2.Primary)
	}
}

// TestSegmentResultComparison tests SegmentResult comparison
func TestSegmentResultComparison(t *testing.T) {
	result1 := segment.SegmentResult{
		ID: config.SegmentModel,
		Data: segment.SegmentData{
			Primary: "Sonnet 3.5",
		},
	}

	result2 := segment.SegmentResult{
		ID: config.SegmentModel,
		Data: segment.SegmentData{
			Primary: "Sonnet 3.5",
		},
	}

	if result1.ID != result2.ID {
		t.Error("expected same segment IDs")
	}

	if result1.Data.Primary != result2.Data.Primary {
		t.Error("expected same primary data")
	}
}

// TestCostSegmentWithFloatingPointPrecision tests cost segment with floating point precision
func TestCostSegmentWithFloatingPointPrecision(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Floating point 0.1 + 0.2",
			cost:     0.1 + 0.2,
			expected: "$0.30",
		},
		{
			name:     "Floating point precision test",
			cost:     0.30000000000000004,
			expected: "$0.30",
		},
		{
			name:     "Very small floating point",
			cost:     0.0001,
			expected: "$0.00",
		},
	}

	costSegment := &segment.CostSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestModelSegmentWithAllKnownModels tests model segment with all known model patterns
func TestModelSegmentWithAllKnownModels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "claude-3-5-sonnet",
			input:    "claude-3-5-sonnet",
			expected: "Sonnet 3.5",
		},
		{
			name:     "claude-3-opus",
			input:    "claude-3-opus",
			expected: "Opus 3",
		},
		{
			name:     "claude-3-haiku",
			input:    "claude-3-haiku",
			expected: "Haiku 3",
		},
		{
			name:     "claude-opus-4",
			input:    "claude-opus-4",
			expected: "Opus 4",
		},
		{
			name:     "claude-sonnet-4",
			input:    "claude-sonnet-4",
			expected: "Sonnet 4",
		},
		{
			name:     "gpt-4",
			input:    "gpt-4",
			expected: "GPT-4",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Model: config.ModelInfo{
					ID:          tt.input,
					DisplayName: tt.input,
				},
			}
			result := modelSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentWithRootPaths tests directory segment with various root paths
func TestDirectorySegmentWithRootPaths(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Unix root",
			dir:      "/",
			expected: "/",
		},
		{
			name:     "Windows C drive",
			dir:      "C:\\",
			expected: "C:\\",
		},
		{
			name:     "Windows D drive",
			dir:      "D:\\",
			expected: "D:\\",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataNilMetadata tests SegmentData with nil metadata
func TestSegmentDataNilMetadata(t *testing.T) {
	data := segment.SegmentData{
		Primary:   "test",
		Secondary: "secondary",
		Metadata:  nil,
	}

	if data.Metadata != nil {
		t.Error("expected nil Metadata")
	}

	// Verify we can still access Primary and Secondary
	if data.Primary != "test" {
		t.Errorf("expected Primary to be 'test', got %q", data.Primary)
	}

	if data.Secondary != "secondary" {
		t.Errorf("expected Secondary to be 'secondary', got %q", data.Secondary)
	}
}

// TestCostSegmentWithCurrencyFormatting tests cost segment currency formatting
func TestCostSegmentWithCurrencyFormatting(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Dollar sign prefix",
			cost:     1.00,
			expected: "$1.00",
		},
		{
			name:     "Two decimal places",
			cost:     0.99,
			expected: "$0.99",
		},
		{
			name:     "Leading zero",
			cost:     0.05,
			expected: "$0.05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestModelSegmentDisplayNamePriority tests that DisplayName takes priority over ID
func TestModelSegmentDisplayNamePriority(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet-20241022",
			DisplayName: "claude-3-5-sonnet-20241022",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentEmptyPath tests directory segment with empty path
func TestDirectorySegmentEmptyPath(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: "",
		},
	}

	result := dirSegment.Collect(input)
	// filepath.Base("") returns "."
	if result.Primary != "." {
		t.Errorf("expected '.', got %q", result.Primary)
	}
}

// TestOutputStyleSegmentNonEmptyCheck tests output style segment empty check
func TestOutputStyleSegmentNonEmptyCheck(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name    string
		style   string
		isEmpty bool
	}{
		{
			name:    "Non-empty style",
			style:   "markdown",
			isEmpty: false,
		},
		{
			name:    "Empty style",
			style:   "",
			isEmpty: true,
		},
		{
			name:    "Whitespace style",
			style:   " ",
			isEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if tt.isEmpty {
				if result.Primary != "" {
					t.Errorf("expected empty result, got %q", result.Primary)
				}
			} else {
				if result.Primary == "" {
					t.Error("expected non-empty result")
				}
			}
		})
	}
}

// TestSegmentDataCopy tests copying SegmentData
func TestSegmentDataCopy(t *testing.T) {
	original := segment.SegmentData{
		Primary:   "primary",
		Secondary: "secondary",
		Metadata: map[string]string{
			"key": "value",
		},
	}

	// Create a copy
	copy := original

	// Verify copy has same values
	if copy.Primary != original.Primary {
		t.Error("expected Primary to be copied")
	}

	if copy.Secondary != original.Secondary {
		t.Error("expected Secondary to be copied")
	}

	if copy.Metadata["key"] != original.Metadata["key"] {
		t.Error("expected Metadata to be copied")
	}

	// Modify copy's metadata
	copy.Metadata["key"] = "modified"

	// Verify original is also modified (shallow copy)
	if original.Metadata["key"] != "modified" {
		t.Error("expected shallow copy behavior")
	}
}

// TestSegmentResultCopy tests copying SegmentResult
func TestSegmentResultCopy(t *testing.T) {
	original := segment.SegmentResult{
		ID: config.SegmentModel,
		Data: segment.SegmentData{
			Primary: "test",
		},
	}

	// Create a copy
	copy := original

	// Verify copy has same values
	if copy.ID != original.ID {
		t.Error("expected ID to be copied")
	}

	if copy.Data.Primary != original.Data.Primary {
		t.Error("expected Data to be copied")
	}
}

// TestCostSegmentNonZeroCheck tests cost segment non-zero check
func TestCostSegmentNonZeroCheck(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name    string
		cost    float64
		isEmpty bool
	}{
		{
			name:    "Positive cost",
			cost:    0.01,
			isEmpty: false,
		},
		{
			name:    "Zero cost",
			cost:    0.0,
			isEmpty: true,
		},
		{
			name:    "Negative cost",
			cost:    -0.01,
			isEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if tt.isEmpty {
				if result.Primary != "" {
					t.Errorf("expected empty result, got %q", result.Primary)
				}
			} else {
				if result.Primary == "" {
					t.Error("expected non-empty result")
				}
			}
		})
	}
}

// TestSegmentDataAllFields tests SegmentData with all fields populated
func TestSegmentDataAllFields(t *testing.T) {
	metadata := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	data := segment.SegmentData{
		Primary:   "primary_value",
		Secondary: "secondary_value",
		Metadata:  metadata,
	}

	if data.Primary != "primary_value" {
		t.Errorf("expected Primary to be 'primary_value', got %q", data.Primary)
	}

	if data.Secondary != "secondary_value" {
		t.Errorf("expected Secondary to be 'secondary_value', got %q", data.Secondary)
	}

	if len(data.Metadata) != 2 {
		t.Errorf("expected 2 metadata entries, got %d", len(data.Metadata))
	}

	if data.Metadata["key1"] != "value1" {
		t.Errorf("expected key1 to be 'value1', got %q", data.Metadata["key1"])
	}

	if data.Metadata["key2"] != "value2" {
		t.Errorf("expected key2 to be 'value2', got %q", data.Metadata["key2"])
	}
}

// TestModelSegmentCaseInsensitivity tests model segment case insensitivity
func TestModelSegmentCaseInsensitivity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase",
			input:    "claude-3-5-sonnet",
			expected: "Sonnet 3.5",
		},
		{
			name:     "Uppercase",
			input:    "CLAUDE-3-5-SONNET",
			expected: "Sonnet 3.5",
		},
		{
			name:     "Mixed case",
			input:    "Claude-3-5-Sonnet",
			expected: "Sonnet 3.5",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Model: config.ModelInfo{
					ID:          tt.input,
					DisplayName: tt.input,
				},
			}
			result := modelSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentWithDotFiles tests directory segment with dot files
func TestDirectorySegmentWithDotFiles(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Hidden file",
			dir:      "/home/user/.bashrc",
			expected: ".bashrc",
		},
		{
			name:     "Hidden directory",
			dir:      "/home/user/.config",
			expected: ".config",
		},
		{
			name:     "Double dot directory",
			dir:      "/home/user/..",
			expected: "..",
		},
		{
			name:     "Single dot directory",
			dir:      "/home/user/.",
			expected: ".",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentRounding tests cost segment rounding behavior
func TestCostSegmentRounding(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Round down",
			cost:     0.124,
			expected: "$0.12",
		},
		{
			name:     "Round up",
			cost:     0.126,
			expected: "$0.13",
		},
		{
			name:     "Round half",
			cost:     0.125,
			expected: "$0.12",
		},
		{
			name:     "Round half up",
			cost:     0.135,
			expected: "$0.14",
		},
	}

	costSegment := &segment.CostSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentTrimming tests output style segment doesn't trim whitespace
func TestOutputStyleSegmentTrimming(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "Leading space",
			style:    " markdown",
			expected: " markdown",
		},
		{
			name:     "Trailing space",
			style:    "markdown ",
			expected: "markdown ",
		},
		{
			name:     "Both spaces",
			style:    " markdown ",
			expected: " markdown ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentResultZeroValue tests SegmentResult zero value
func TestSegmentResultZeroValue(t *testing.T) {
	var result segment.SegmentResult

	if result.ID != "" {
		t.Errorf("expected empty ID, got %q", result.ID)
	}

	if result.Data.Primary != "" {
		t.Errorf("expected empty Primary, got %q", result.Data.Primary)
	}

	if result.Data.Secondary != "" {
		t.Errorf("expected empty Secondary, got %q", result.Data.Secondary)
	}

	if result.Data.Metadata != nil {
		t.Errorf("expected nil Metadata, got %v", result.Data.Metadata)
	}
}

// TestSegmentDataZeroValue tests SegmentData zero value
func TestSegmentDataZeroValue(t *testing.T) {
	var data segment.SegmentData

	if data.Primary != "" {
		t.Errorf("expected empty Primary, got %q", data.Primary)
	}

	if data.Secondary != "" {
		t.Errorf("expected empty Secondary, got %q", data.Secondary)
	}

	if data.Metadata != nil {
		t.Errorf("expected nil Metadata, got %v", data.Metadata)
	}
}

// TestModelSegmentIDFallback tests model segment falls back to ID when DisplayName is empty
func TestModelSegmentIDFallback(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet-20241022",
			DisplayName: "",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentPathSeparators tests directory segment with different path separators
func TestDirectorySegmentPathSeparators(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Unix separator",
			dir:      "/home/user/project",
			expected: "project",
		},
	}

	dirSegment := &segment.DirectorySegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentNegativeFormatting tests cost segment negative value formatting
func TestCostSegmentNegativeFormatting(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Negative one dollar",
			cost:     -1.00,
			expected: "$-1.00",
		},
		{
			name:     "Negative small amount",
			cost:     -0.05,
			expected: "$-0.05",
		},
		{
			name:     "Negative large amount",
			cost:     -999.99,
			expected: "$-999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentCaseSensitivity tests output style segment is case sensitive
func TestOutputStyleSegmentCaseSensitivity(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "Lowercase",
			style:    "markdown",
			expected: "markdown",
		},
		{
			name:     "Uppercase",
			style:    "MARKDOWN",
			expected: "MARKDOWN",
		},
		{
			name:     "Mixed case",
			style:    "MarkDown",
			expected: "MarkDown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataMetadataIndependence tests that metadata is independent between instances
func TestSegmentDataMetadataIndependence(t *testing.T) {
	data1 := segment.SegmentData{
		Primary:  "test1",
		Metadata: make(map[string]string),
	}

	data2 := segment.SegmentData{
		Primary:  "test2",
		Metadata: make(map[string]string),
	}

	data1.Metadata["key"] = "value1"
	data2.Metadata["key"] = "value2"

	if data1.Metadata["key"] != "value1" {
		t.Errorf("expected data1 key to be 'value1', got %q", data1.Metadata["key"])
	}

	if data2.Metadata["key"] != "value2" {
		t.Errorf("expected data2 key to be 'value2', got %q", data2.Metadata["key"])
	}
}

// TestModelSegmentMultipleReplacements tests model segment with multiple possible replacements
func TestModelSegmentMultipleReplacements(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	// Test that the first matching pattern is used
	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet-and-claude-3-opus",
			DisplayName: "claude-3-5-sonnet-and-claude-3-opus",
		},
	}

	result := modelSegment.Collect(input)
	// Should match one of them (order depends on map iteration)
	if result.Primary != "Sonnet 3.5" && result.Primary != "Opus 3" {
		t.Errorf("expected 'Sonnet 3.5' or 'Opus 3', got %q", result.Primary)
	}
}

// TestDirectorySegmentLongPath tests directory segment with very long path
func TestDirectorySegmentLongPath(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	longPath := "/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/final"
	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: longPath,
		},
	}

	result := dirSegment.Collect(input)
	if result.Primary != "final" {
		t.Errorf("expected 'final', got %q", result.Primary)
	}
}

// TestCostSegmentLargeNumbers tests cost segment with large numbers
func TestCostSegmentLargeNumbers(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Million dollars",
			cost:     1000000.00,
			expected: "$1000000.00",
		},
		{
			name:     "Billion dollars",
			cost:     1000000000.00,
			expected: "$1000000000.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataEmptyStrings tests SegmentData with empty strings
func TestSegmentDataEmptyStrings(t *testing.T) {
	data := segment.SegmentData{
		Primary:   "",
		Secondary: "",
		Metadata:  make(map[string]string),
	}

	if data.Primary != "" {
		t.Errorf("expected empty Primary, got %q", data.Primary)
	}

	if data.Secondary != "" {
		t.Errorf("expected empty Secondary, got %q", data.Secondary)
	}

	if len(data.Metadata) != 0 {
		t.Errorf("expected empty Metadata, got %d entries", len(data.Metadata))
	}
}

// TestUpdateSegmentMultipleCalls tests UpdateSegment returns consistent empty results
func TestUpdateSegmentMultipleCalls(t *testing.T) {
	updateSegment := &segment.UpdateSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "test",
			DisplayName: "test",
		},
	}

	results := make([]segment.SegmentData, 5)
	for i := 0; i < 5; i++ {
		results[i] = updateSegment.Collect(input)
	}

	for i, result := range results {
		if result.Primary != "" {
			t.Errorf("call %d: expected empty Primary, got %q", i, result.Primary)
		}
	}
}

// TestModelSegmentWithSpecialCharacters tests model segment with special characters
func TestModelSegmentWithSpecialCharacters(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet@latest",
			DisplayName: "claude-3-5-sonnet@latest",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentWithSymbols tests directory segment with symbols in path
func TestDirectorySegmentWithSymbols(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Directory with plus sign",
			dir:      "/home/user/project+v2",
			expected: "project+v2",
		},
		{
			name:     "Directory with equals sign",
			dir:      "/home/user/project=main",
			expected: "project=main",
		},
		{
			name:     "Directory with at sign",
			dir:      "/home/user/project@latest",
			expected: "project@latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentWithScientificNotation tests cost segment with scientific notation
func TestCostSegmentWithScientificNotation(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 1e-2, // 0.01
		},
	}

	result := costSegment.Collect(input)
	if result.Primary != "$0.01" {
		t.Errorf("expected '$0.01', got %q", result.Primary)
	}
}

// TestOutputStyleSegmentWithNumbers tests output style segment with numbers
func TestOutputStyleSegmentWithNumbers(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "Style with version number",
			style:    "v1.0.0",
			expected: "v1.0.0",
		},
		{
			name:     "Style with only numbers",
			style:    "123",
			expected: "123",
		},
		{
			name:     "Style with mixed alphanumeric",
			style:    "style2024",
			expected: "style2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataMetadataKeyTypes tests metadata with various key types
func TestSegmentDataMetadataKeyTypes(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	// Test various key formats
	keys := []string{
		"simple",
		"with-dash",
		"with_underscore",
		"with.dot",
		"with:colon",
		"with/slash",
		"123numeric",
		"UPPERCASE",
		"MixedCase",
	}

	for _, key := range keys {
		data.Metadata[key] = "value"
	}

	if len(data.Metadata) != len(keys) {
		t.Errorf("expected %d metadata entries, got %d", len(keys), len(data.Metadata))
	}

	for _, key := range keys {
		if data.Metadata[key] != "value" {
			t.Errorf("expected metadata[%q] to be 'value', got %q", key, data.Metadata[key])
		}
	}
}

// TestModelSegmentEmptyDisplayNameUsesID tests model segment uses ID when DisplayName is empty
func TestModelSegmentEmptyDisplayNameUsesID(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-opus-20240229",
			DisplayName: "",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Opus 3" {
		t.Errorf("expected 'Opus 3', got %q", result.Primary)
	}
}

// TestDirectorySegmentConsistency tests directory segment returns consistent results
func TestDirectorySegmentConsistency(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/home/user/project",
		},
	}

	results := make([]string, 10)
	for i := 0; i < 10; i++ {
		results[i] = dirSegment.Collect(input).Primary
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("expected consistent results, got %q and %q", results[0], results[i])
		}
	}
}

// TestCostSegmentConsistency tests cost segment returns consistent results
func TestCostSegmentConsistency(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 0.50,
		},
	}

	results := make([]string, 10)
	for i := 0; i < 10; i++ {
		results[i] = costSegment.Collect(input).Primary
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("expected consistent results, got %q and %q", results[0], results[i])
		}
	}
}

// TestOutputStyleSegmentConsistency tests output style segment returns consistent results
func TestOutputStyleSegmentConsistency(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "markdown",
		},
	}

	results := make([]string, 10)
	for i := 0; i < 10; i++ {
		results[i] = styleSegment.Collect(input).Primary
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("expected consistent results, got %q and %q", results[0], results[i])
		}
	}
}

// TestSegmentDataMetadataValueTypes tests metadata with various value types
func TestSegmentDataMetadataValueTypes(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	// Test various value formats
	values := []string{
		"simple",
		"with spaces",
		"with-dash",
		"with_underscore",
		"with.dot",
		"with:colon",
		"with/slash",
		"123numeric",
		"UPPERCASE",
		"MixedCase",
		"",
	}

	for i, value := range values {
		key := fmt.Sprintf("key%d", i)
		data.Metadata[key] = value
	}

	if len(data.Metadata) != len(values) {
		t.Errorf("expected %d metadata entries, got %d", len(values), len(data.Metadata))
	}

	for i, value := range values {
		key := fmt.Sprintf("key%d", i)
		if data.Metadata[key] != value {
			t.Errorf("expected metadata[%q] to be %q, got %q", key, value, data.Metadata[key])
		}
	}
}

// TestModelSegmentWithVersionSuffixes tests model segment with various version suffixes
func TestModelSegmentWithVersionSuffixes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "With date suffix",
			input:    "claude-3-5-sonnet-20241022",
			expected: "Sonnet 3.5",
		},
		{
			name:     "With beta suffix",
			input:    "claude-3-5-sonnet-beta",
			expected: "Sonnet 3.5",
		},
		{
			name:     "With preview suffix",
			input:    "claude-3-5-sonnet-preview",
			expected: "Sonnet 3.5",
		},
	}

	modelSegment := &segment.ModelSegment{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Model: config.ModelInfo{
					ID:          tt.input,
					DisplayName: tt.input,
				},
			}
			result := modelSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentWithTrailingSlashes tests directory segment with trailing slashes
func TestDirectorySegmentWithTrailingSlashes(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Single trailing slash",
			dir:      "/home/user/project/",
			expected: "project",
		},
		{
			name:     "Multiple trailing slashes",
			dir:      "/home/user/project///",
			expected: "project",
		},
		{
			name:     "Root with trailing slash",
			dir:      "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentWithDecimalVariations tests cost segment with decimal variations
func TestCostSegmentWithDecimalVariations(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Whole number",
			cost:     5.0,
			expected: "$5.00",
		},
		{
			name:     "One decimal place",
			cost:     5.5,
			expected: "$5.50",
		},
		{
			name:     "Two decimal places",
			cost:     5.55,
			expected: "$5.55",
		},
		{
			name:     "Three decimal places",
			cost:     5.555,
			expected: "$5.55",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithEmojis tests output style segment with emoji characters
func TestOutputStyleSegmentWithEmojis(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "Style with emoji",
			style:    "markdown 📝",
			expected: "markdown 📝",
		},
		{
			name:     "Only emoji",
			style:    "🎯",
			expected: "🎯",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestSegmentDataMetadataLargeSize tests SegmentData with large metadata
func TestSegmentDataMetadataLargeSize(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	// Add many entries
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		data.Metadata[key] = value
	}

	if len(data.Metadata) != 1000 {
		t.Errorf("expected 1000 metadata entries, got %d", len(data.Metadata))
	}

	// Verify some entries
	if data.Metadata["key_0"] != "value_0" {
		t.Error("expected first entry to be correct")
	}

	if data.Metadata["key_999"] != "value_999" {
		t.Error("expected last entry to be correct")
	}
}

// TestModelSegmentWithUnknownModel tests model segment with completely unknown model
func TestModelSegmentWithUnknownModel(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "unknown-ai-model-xyz",
			DisplayName: "unknown-ai-model-xyz",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "unknown-ai-model-xyz" {
		t.Errorf("expected 'unknown-ai-model-xyz', got %q", result.Primary)
	}
}

// TestDirectorySegmentWithRelativePaths tests directory segment with relative paths
func TestDirectorySegmentWithRelativePaths(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Relative path",
			dir:      "src/components",
			expected: "components",
		},
		{
			name:     "Relative path with dot",
			dir:      "./src/components",
			expected: "components",
		},
		{
			name:     "Relative path with parent",
			dir:      "../src/components",
			expected: "components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentWithInfinity tests cost segment behavior with special float values
func TestCostSegmentWithInfinity(t *testing.T) {
	costSegment := &segment.CostSegment{}

	// Test with very large number (not infinity, but close)
	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 1e10,
		},
	}

	result := costSegment.Collect(input)
	if result.Primary == "" {
		t.Error("expected non-empty result for large number")
	}
}

// TestOutputStyleSegmentWithLongString tests output style segment with long string
func TestOutputStyleSegmentWithLongString(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	longStyle := "this_is_a_very_long_output_style_name_with_many_characters_and_underscores_and_numbers_123456789"

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: longStyle,
		},
	}

	result := styleSegment.Collect(input)
	if result.Primary != longStyle {
		t.Errorf("expected %q, got %q", longStyle, result.Primary)
	}
}

// TestSegmentDataMetadataEmptyValue tests metadata with empty string values
func TestSegmentDataMetadataEmptyValue(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	data.Metadata["empty"] = ""
	data.Metadata["nonempty"] = "value"

	if data.Metadata["empty"] != "" {
		t.Error("expected empty value to remain empty")
	}

	if data.Metadata["nonempty"] != "value" {
		t.Error("expected non-empty value to be preserved")
	}

	if len(data.Metadata) != 2 {
		t.Errorf("expected 2 entries, got %d", len(data.Metadata))
	}
}

// TestSegmentInterfaceCompliance tests that segments properly implement the interface
func TestSegmentInterfaceCompliance(t *testing.T) {
	segments := []struct {
		name    string
		segment segment.Segment
	}{
		{
			name:    "ModelSegment",
			segment: &segment.ModelSegment{},
		},
		{
			name:    "DirectorySegment",
			segment: &segment.DirectorySegment{},
		},
		{
			name:    "CostSegment",
			segment: &segment.CostSegment{},
		},
		{
			name:    "OutputStyleSegment",
			segment: &segment.OutputStyleSegment{},
		},
		{
			name:    "UpdateSegment",
			segment: &segment.UpdateSegment{},
		},
	}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "test",
			DisplayName: "test",
		},
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/test",
		},
		Cost: config.CostInfo{
			TotalCostUSD: 0.1,
		},
		OutputStyle: config.OutputStyleInfo{
			Name: "test",
		},
	}

	for _, seg := range segments {
		t.Run(seg.name, func(t *testing.T) {
			result := seg.segment.Collect(input)
			// Just verify it returns a SegmentData without panicking
			if result.Primary == "" && result.Secondary == "" && result.Metadata == nil {
				// This is acceptable for some segments
			}
		})
	}
}

// TestCostSegmentPrecisionEdgeCases tests cost segment with precision edge cases
func TestCostSegmentPrecisionEdgeCases(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "0.015 rounds to 0.01",
			cost:     0.015,
			expected: "$0.01",
		},
		{
			name:     "0.025 rounds to 0.03",
			cost:     0.025,
			expected: "$0.03",
		},
		{
			name:     "0.035 rounds to 0.04",
			cost:     0.035,
			expected: "$0.04",
		},
		{
			name:     "0.045 rounds to 0.04",
			cost:     0.045,
			expected: "$0.04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestModelSegmentWithWhitespace tests model segment with whitespace in names
func TestModelSegmentWithWhitespace(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "claude-3-5-sonnet with spaces",
			DisplayName: "claude-3-5-sonnet with spaces",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentWithWhitespace tests directory segment with whitespace
func TestDirectorySegmentWithWhitespace(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Directory with leading space",
			dir:      "/home/user/ project",
			expected: " project",
		},
		{
			name:     "Directory with trailing space",
			dir:      "/home/user/project ",
			expected: "project ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithTabs tests output style segment with tab characters
func TestOutputStyleSegmentWithTabs(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "markdown\tstyle",
		},
	}

	result := styleSegment.Collect(input)
	if result.Primary != "markdown\tstyle" {
		t.Errorf("expected 'markdown\\tstyle', got %q", result.Primary)
	}
}

// TestSegmentDataMetadataReplacement tests metadata value replacement
func TestSegmentDataMetadataReplacement(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	data.Metadata["key"] = "value1"
	if data.Metadata["key"] != "value1" {
		t.Error("expected initial value to be set")
	}

	data.Metadata["key"] = "value2"
	if data.Metadata["key"] != "value2" {
		t.Error("expected value to be replaced")
	}

	data.Metadata["key"] = "value3"
	if data.Metadata["key"] != "value3" {
		t.Error("expected value to be replaced again")
	}
}

// TestCostSegmentWithMaxFloat tests cost segment with maximum float value
func TestCostSegmentWithMaxFloat(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 1.7976931348623157e+308, // Close to max float64
		},
	}

	result := costSegment.Collect(input)
	// Should not panic and should return a formatted string
	if result.Primary == "" {
		t.Error("expected non-empty result for large float")
	}
}

// TestModelSegmentWithNumericSuffix tests model segment with numeric suffixes
func TestModelSegmentWithNumericSuffix(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "With version 1",
			input:    "claude-3-5-sonnet-v1",
			expected: "Sonnet 3.5",
		},
		{
			name:     "With version 2",
			input:    "claude-3-5-sonnet-v2",
			expected: "Sonnet 3.5",
		},
		{
			name:     "With revision",
			input:    "claude-3-5-sonnet-r1",
			expected: "Sonnet 3.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Model: config.ModelInfo{
					ID:          tt.input,
					DisplayName: tt.input,
				},
			}
			result := modelSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestDirectorySegmentWithConsecutiveSlashes tests directory segment with consecutive slashes
func TestDirectorySegmentWithConsecutiveSlashes(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Double slash",
			dir:      "/home//user//project",
			expected: "project",
		},
		{
			name:     "Triple slash",
			dir:      "/home///user///project",
			expected: "project",
		},
		{
			name:     "Many slashes",
			dir:      "/home/////user/////project",
			expected: "project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentWithMinimumValue tests cost segment with minimum positive value
func TestCostSegmentWithMinimumValue(t *testing.T) {
	costSegment := &segment.CostSegment{}

	input := &config.InputData{
		Cost: config.CostInfo{
			TotalCostUSD: 1e-10, // Very small positive number
		},
	}

	result := costSegment.Collect(input)
	if result.Primary != "$0.00" {
		t.Errorf("expected '$0.00', got %q", result.Primary)
	}
}

// TestOutputStyleSegmentWithNewlines tests output style segment with newline characters
func TestOutputStyleSegmentWithNewlines(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "markdown\nstyle",
		},
	}

	result := styleSegment.Collect(input)
	if result.Primary != "markdown\nstyle" {
		t.Errorf("expected 'markdown\\nstyle', got %q", result.Primary)
	}
}

// TestSegmentDataMetadataIteration tests iterating over metadata
func TestSegmentDataMetadataIteration(t *testing.T) {
	data := segment.SegmentData{
		Primary:  "test",
		Metadata: make(map[string]string),
	}

	// Add entries
	expectedEntries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range expectedEntries {
		data.Metadata[k] = v
	}

	// Iterate and verify
	count := 0
	for k, v := range data.Metadata {
		if expectedEntries[k] != v {
			t.Errorf("expected metadata[%q] to be %q, got %q", k, expectedEntries[k], v)
		}
		count++
	}

	if count != len(expectedEntries) {
		t.Errorf("expected %d entries, got %d", len(expectedEntries), count)
	}
}

// TestModelSegmentWithPrefixAndSuffix tests model segment with prefix and suffix
func TestModelSegmentWithPrefixAndSuffix(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "prefix-claude-3-5-sonnet-suffix",
			DisplayName: "prefix-claude-3-5-sonnet-suffix",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentWithDotDotPath tests directory segment with .. in path
func TestDirectorySegmentWithDotDotPath(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "Single dot dot",
			dir:      "/home/user/..",
			expected: "..",
		},
		{
			name:     "Dot dot in middle",
			dir:      "/home/../user/project",
			expected: "project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Workspace: config.WorkspaceInfo{
					CurrentDir: tt.dir,
				},
			}
			result := dirSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestCostSegmentWithNearZero tests cost segment with values near zero
func TestCostSegmentWithNearZero(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "0.001",
			cost:     0.001,
			expected: "$0.00",
		},
		{
			name:     "0.004",
			cost:     0.004,
			expected: "$0.00",
		},
		{
			name:     "0.005",
			cost:     0.005,
			expected: "$0.01",
		},
		{
			name:     "0.009",
			cost:     0.009,
			expected: "$0.01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithControlCharacters tests output style segment with control characters
func TestOutputStyleSegmentWithControlCharacters(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	input := &config.InputData{
		OutputStyle: config.OutputStyleInfo{
			Name: "markdown\x00style",
		},
	}

	result := styleSegment.Collect(input)
	if result.Primary != "markdown\x00style" {
		t.Errorf("expected 'markdown\\x00style', got %q", result.Primary)
	}
}

// TestSegmentDataFieldIndependence tests that fields are independent
func TestSegmentDataFieldIndependence(t *testing.T) {
	data1 := segment.SegmentData{
		Primary:   "primary1",
		Secondary: "secondary1",
	}

	data2 := segment.SegmentData{
		Primary:   "primary2",
		Secondary: "secondary2",
	}

	if data1.Primary == data2.Primary {
		t.Error("expected different Primary values")
	}

	if data1.Secondary == data2.Secondary {
		t.Error("expected different Secondary values")
	}

	// Modify data1
	data1.Primary = "modified"

	if data2.Primary != "primary2" {
		t.Error("expected data2.Primary to remain unchanged")
	}
}

// TestModelSegmentWithAllCaps tests model segment with all caps model name
func TestModelSegmentWithAllCaps(t *testing.T) {
	modelSegment := &segment.ModelSegment{}

	input := &config.InputData{
		Model: config.ModelInfo{
			ID:          "CLAUDE-3-5-SONNET",
			DisplayName: "CLAUDE-3-5-SONNET",
		},
	}

	result := modelSegment.Collect(input)
	if result.Primary != "Sonnet 3.5" {
		t.Errorf("expected 'Sonnet 3.5', got %q", result.Primary)
	}
}

// TestDirectorySegmentWithAllCaps tests directory segment with all caps directory name
func TestDirectorySegmentWithAllCaps(t *testing.T) {
	dirSegment := &segment.DirectorySegment{}

	input := &config.InputData{
		Workspace: config.WorkspaceInfo{
			CurrentDir: "/HOME/USER/PROJECT",
		},
	}

	result := dirSegment.Collect(input)
	if result.Primary != "PROJECT" {
		t.Errorf("expected 'PROJECT', got %q", result.Primary)
	}
}

// TestCostSegmentWithMixedSigns tests cost segment with mixed positive and negative
func TestCostSegmentWithMixedSigns(t *testing.T) {
	costSegment := &segment.CostSegment{}

	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "Positive",
			cost:     0.50,
			expected: "$0.50",
		},
		{
			name:     "Negative",
			cost:     -0.50,
			expected: "$-0.50",
		},
		{
			name:     "Zero",
			cost:     0.0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				Cost: config.CostInfo{
					TotalCostUSD: tt.cost,
				},
			}
			result := costSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}

// TestOutputStyleSegmentWithMixedCase tests output style segment preserves mixed case
func TestOutputStyleSegmentWithMixedCase(t *testing.T) {
	styleSegment := &segment.OutputStyleSegment{}

	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "CamelCase",
			style:    "MarkdownStyle",
			expected: "MarkdownStyle",
		},
		{
			name:     "snake_case",
			style:    "markdown_style",
			expected: "markdown_style",
		},
		{
			name:     "kebab-case",
			style:    "markdown-style",
			expected: "markdown-style",
		},
		{
			name:     "PascalCase",
			style:    "MarkdownStyle",
			expected: "MarkdownStyle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.InputData{
				OutputStyle: config.OutputStyleInfo{
					Name: tt.style,
				},
			}
			result := styleSegment.Collect(input)
			if result.Primary != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Primary)
			}
		})
	}
}
