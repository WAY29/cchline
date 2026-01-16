package tests

import (
	"strings"
	"testing"

	"github.com/WAY29/cchline/config"
	"github.com/WAY29/cchline/render"
	"github.com/WAY29/cchline/segment"
)

// Helper function to create a test config
func createTestConfig(theme config.ThemeMode, separator string) *config.SimpleConfig {
	return &config.SimpleConfig{
		Theme:     theme,
		Separator: separator,
		Segments: config.SegmentToggles{
			Model:         true,
			Directory:     true,
			Git:           true,
			ContextWindow: true,
		},
	}
}

// Helper function to create a test segment result
func createTestSegment(id config.SegmentID, primary string) segment.SegmentResult {
	return segment.SegmentResult{
		ID: id,
		Data: segment.SegmentData{
			Primary:   primary,
			Secondary: "",
			Metadata:  make(map[string]string),
		},
	}
}

// TestNewStatusLineGenerator tests the creation of StatusLineGenerator
func TestNewStatusLineGenerator(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	if generator == nil {
		t.Error("NewStatusLineGenerator returned nil")
	}
}

// TestGenerateEmptySegments tests Generate with empty segments list
func TestGenerateEmptySegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	result := generator.Generate([]segment.SegmentResult{})

	if result != "" {
		t.Errorf("Expected empty string for empty segments, got: %q", result)
	}
}

// TestGenerateSingleSegment tests Generate with a single segment
func TestGenerateSingleSegment(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for single segment")
	}

	// Check that the result contains the segment data
	if !strings.Contains(result, "Sonnet 3.5") {
		t.Errorf("Expected result to contain 'Sonnet 3.5', got: %q", result)
	}
}

// TestGenerateMultipleSegments tests Generate with multiple segments
func TestGenerateMultipleSegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentDirectory, "myapp"),
		createTestSegment(config.SegmentGit, "main *"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for multiple segments")
	}

	// Check that all segment data is present
	if !strings.Contains(result, "Sonnet 3.5") {
		t.Errorf("Expected result to contain 'Sonnet 3.5', got: %q", result)
	}
	if !strings.Contains(result, "myapp") {
		t.Errorf("Expected result to contain 'myapp', got: %q", result)
	}
	if !strings.Contains(result, "main *") {
		t.Errorf("Expected result to contain 'main *', got: %q", result)
	}
}

// TestGenerateSeparatorJoining tests that segments are properly joined with separator
func TestGenerateSeparatorJoining(t *testing.T) {
	separator := " | "
	cfg := createTestConfig(config.ThemeModeDefault, separator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentDirectory, "myapp"),
	}

	result := generator.Generate(segments)

	// Check that separator is present
	if !strings.Contains(result, separator) {
		t.Errorf("Expected result to contain separator %q, got: %q", separator, result)
	}
}

// TestGenerateCustomSeparator tests Generate with custom separator
func TestGenerateCustomSeparator(t *testing.T) {
	customSeparator := " :: "
	cfg := createTestConfig(config.ThemeModeDefault, customSeparator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentDirectory, "myapp"),
	}

	result := generator.Generate(segments)

	// Check that custom separator is present
	if !strings.Contains(result, customSeparator) {
		t.Errorf("Expected result to contain custom separator %q, got: %q", customSeparator, result)
	}
}

// TestGenerateDefaultTheme tests Generate with default (emoji) theme
func TestGenerateDefaultTheme(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with default theme")
	}

	// Default theme should contain emoji icon
	if !strings.Contains(result, "🤖") {
		t.Errorf("Expected default theme to contain model emoji, got: %q", result)
	}
}

// TestGenerateNerdFontTheme tests Generate with nerd font theme
func TestGenerateNerdFontTheme(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeNerdFont, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with nerd font theme")
	}

	// Nerd font theme should contain nerd font icon (not emoji)
	if strings.Contains(result, "🤖") {
		t.Errorf("Expected nerd font theme to not contain emoji, got: %q", result)
	}
}

// TestGenerateMultipleSegmentsWithDifferentTypes tests multiple different segment types
func TestGenerateMultipleSegmentsWithDifferentTypes(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentDirectory, "myapp"),
		createTestSegment(config.SegmentGit, "main *"),
		createTestSegment(config.SegmentContextWindow, "15.6% · 31.1k tokens"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for multiple different segment types")
	}

	// Verify all segments are present
	expectedStrings := []string{
		"Sonnet 3.5",
		"myapp",
		"main *",
		"15.6% · 31.1k tokens",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected result to contain %q, got: %q", expected, result)
		}
	}
}

// TestGenerateWithEmptyPrimaryData tests segment with empty primary data
func TestGenerateWithEmptyPrimaryData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentUpdate, ""), // Empty primary data
	}

	result := generator.Generate(segments)

	// Should still contain the model segment
	if !strings.Contains(result, "Sonnet 3.5") {
		t.Errorf("Expected result to contain 'Sonnet 3.5', got: %q", result)
	}
}

// TestGenerateCCHSegments tests CCH-related segments
func TestGenerateCCHSegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCCHModel, "claude-3-opus"),
		createTestSegment(config.SegmentCCHProvider, "anthropic"),
		createTestSegment(config.SegmentCCHCost, "$1.50/$10"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for CCH segments")
	}

	// Verify CCH segments are present
	if !strings.Contains(result, "claude-3-opus") {
		t.Errorf("Expected result to contain 'claude-3-opus', got: %q", result)
	}
	if !strings.Contains(result, "anthropic") {
		t.Errorf("Expected result to contain 'anthropic', got: %q", result)
	}
	if !strings.Contains(result, "$1.50/$10") {
		t.Errorf("Expected result to contain '$1.50/$10', got: %q", result)
	}
}

// TestGenerateSegmentOrderPreservation tests that segment order is preserved in output
func TestGenerateSegmentOrderPreservation(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
		createTestSegment(config.SegmentGit, "C"),
	}

	result := generator.Generate(segments)

	// Find positions of each segment in the result
	posA := strings.Index(result, "A")
	posB := strings.Index(result, "B")
	posC := strings.Index(result, "C")

	if posA == -1 || posB == -1 || posC == -1 {
		t.Errorf("Expected all segments to be present in result: %q", result)
	}

	// Verify order is preserved
	if posA > posB || posB > posC {
		t.Errorf("Expected segment order to be preserved (A before B before C), got: %q", result)
	}
}

// TestGenerateWithSpecialCharacters tests segments with special characters
func TestGenerateWithSpecialCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model-v1.0"),
		createTestSegment(config.SegmentDirectory, "/path/to/dir"),
		createTestSegment(config.SegmentGit, "feature/test-branch"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with special characters")
	}

	// Verify special characters are preserved
	if !strings.Contains(result, "Model-v1.0") {
		t.Errorf("Expected result to contain 'Model-v1.0', got: %q", result)
	}
	if !strings.Contains(result, "/path/to/dir") {
		t.Errorf("Expected result to contain '/path/to/dir', got: %q", result)
	}
	if !strings.Contains(result, "feature/test-branch") {
		t.Errorf("Expected result to contain 'feature/test-branch', got: %q", result)
	}
}

// TestGenerateWithUnicodeCharacters tests segments with unicode characters
func TestGenerateWithUnicodeCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "模型-3.5"),
		createTestSegment(config.SegmentDirectory, "目录"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with unicode characters")
	}

	// Verify unicode characters are preserved
	if !strings.Contains(result, "模型-3.5") {
		t.Errorf("Expected result to contain '模型-3.5', got: %q", result)
	}
	if !strings.Contains(result, "目录") {
		t.Errorf("Expected result to contain '目录', got: %q", result)
	}
}

// TestGenerateWithLongSegmentData tests segments with long data
func TestGenerateWithLongSegmentData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	longData := "This is a very long segment data that contains a lot of information and should still be rendered correctly"
	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, longData),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with long segment data")
	}

	if !strings.Contains(result, longData) {
		t.Errorf("Expected result to contain long data, got: %q", result)
	}
}

// TestGenerateNoSeparatorForSingleSegment tests that no separator appears for single segment
func TestGenerateNoSeparatorForSingleSegment(t *testing.T) {
	separator := " | "
	cfg := createTestConfig(config.ThemeModeDefault, separator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
	}

	result := generator.Generate(segments)

	// Count occurrences of separator - should be 0 for single segment
	count := strings.Count(result, separator)
	if count > 0 {
		t.Errorf("Expected no separator for single segment, but found %d occurrences in: %q", count, result)
	}
}

// TestGenerateMultipleSeparators tests correct number of separators for multiple segments
func TestGenerateMultipleSeparators(t *testing.T) {
	separator := " | "
	cfg := createTestConfig(config.ThemeModeDefault, separator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
		createTestSegment(config.SegmentGit, "C"),
		createTestSegment(config.SegmentContextWindow, "D"),
	}

	result := generator.Generate(segments)

	// For 4 segments, there should be 3 separators
	count := strings.Count(result, separator)
	if count != 3 {
		t.Errorf("Expected 3 separators for 4 segments, but found %d in: %q", count, result)
	}
}

// TestGenerateWithMetadata tests segments with metadata
func TestGenerateWithMetadata(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segmentWithMetadata := segment.SegmentResult{
		ID: config.SegmentModel,
		Data: segment.SegmentData{
			Primary: "Sonnet 3.5",
			Secondary: "Claude",
			Metadata: map[string]string{
				"version": "3.5",
				"provider": "Anthropic",
			},
		},
	}

	result := generator.Generate([]segment.SegmentResult{segmentWithMetadata})

	if result == "" {
		t.Error("Expected non-empty result for segment with metadata")
	}

	if !strings.Contains(result, "Sonnet 3.5") {
		t.Errorf("Expected result to contain primary data, got: %q", result)
	}
}

// TestGenerateAllSegmentTypes tests all available segment types
func TestGenerateAllSegmentTypes(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	allSegmentTypes := []config.SegmentID{
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

	segments := make([]segment.SegmentResult, len(allSegmentTypes))
	for i, segType := range allSegmentTypes {
		segments[i] = createTestSegment(segType, "test-data")
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for all segment types")
	}

	// Verify all test data is present
	if !strings.Contains(result, "test-data") {
		t.Errorf("Expected result to contain test data, got: %q", result)
	}
}

// TestGenerateWithCostData tests segments with cost/currency data
func TestGenerateWithCostData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "$0.15"),
		createTestSegment(config.SegmentCCHCost, "$1.50/$10"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for cost data")
	}

	if !strings.Contains(result, "$0.15") {
		t.Errorf("Expected result to contain '$0.15', got: %q", result)
	}
	if !strings.Contains(result, "$1.50/$10") {
		t.Errorf("Expected result to contain '$1.50/$10', got: %q", result)
	}
}

// TestGenerateWithPercentageData tests segments with percentage data
func TestGenerateWithPercentageData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentContextWindow, "15.6% · 31.1k tokens"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for percentage data")
	}

	if !strings.Contains(result, "15.6%") {
		t.Errorf("Expected result to contain '15.6%%', got: %q", result)
	}
}

// TestGenerateWithTimeData tests segments with time data
func TestGenerateWithTimeData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentSession, "1h23m"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for time data")
	}

	if !strings.Contains(result, "1h23m") {
		t.Errorf("Expected result to contain '1h23m', got: %q", result)
	}
}

// TestGenerateWithTokenCountData tests segments with token count data
func TestGenerateWithTokenCountData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentUsage, "↓31.1K ↑5.2K"),
		createTestSegment(config.SegmentCCHRequests, "123 reqs"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for token count data")
	}

	if !strings.Contains(result, "31.1K") {
		t.Errorf("Expected result to contain '31.1K', got: %q", result)
	}
	if !strings.Contains(result, "123 reqs") {
		t.Errorf("Expected result to contain '123 reqs', got: %q", result)
	}
}

// TestGenerateWithEmptySeparator tests Generate with empty separator
func TestGenerateWithEmptySeparator(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, "")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with empty separator")
	}

	// Should contain both segments without separator between them
	if !strings.Contains(result, "A") || !strings.Contains(result, "B") {
		t.Errorf("Expected result to contain both segments, got: %q", result)
	}
}

// TestGenerateWithMultiCharacterSeparator tests Generate with multi-character separator
func TestGenerateWithMultiCharacterSeparator(t *testing.T) {
	multiCharSeparator := " :: "
	cfg := createTestConfig(config.ThemeModeDefault, multiCharSeparator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
		createTestSegment(config.SegmentGit, "C"),
	}

	result := generator.Generate(segments)

	// Count occurrences of multi-character separator
	count := strings.Count(result, multiCharSeparator)
	if count != 2 {
		t.Errorf("Expected 2 occurrences of multi-character separator, but found %d in: %q", count, result)
	}
}

// TestGenerateLargeNumberOfSegments tests Generate with many segments
func TestGenerateLargeNumberOfSegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	// Create 20 segments
	segments := make([]segment.SegmentResult, 20)
	for i := 0; i < 20; i++ {
		segments[i] = createTestSegment(config.SegmentModel, "segment")
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for large number of segments")
	}

	// Should have 19 separators for 20 segments
	count := strings.Count(result, " | ")
	if count != 19 {
		t.Errorf("Expected 19 separators for 20 segments, but found %d", count)
	}
}

// TestGenerateWithWhitespaceInData tests segments with leading/trailing whitespace
func TestGenerateWithWhitespaceInData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "  Sonnet 3.5  "),
		createTestSegment(config.SegmentDirectory, "\tmyapp\t"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with whitespace in data")
	}

	// Whitespace should be preserved
	if !strings.Contains(result, "  Sonnet 3.5  ") {
		t.Errorf("Expected result to preserve whitespace, got: %q", result)
	}
}

// TestGenerateWithNumericData tests segments with numeric data
func TestGenerateWithNumericData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCCHRequests, "12345"),
		createTestSegment(config.SegmentCost, "999.99"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with numeric data")
	}

	if !strings.Contains(result, "12345") {
		t.Errorf("Expected result to contain '12345', got: %q", result)
	}
	if !strings.Contains(result, "999.99") {
		t.Errorf("Expected result to contain '999.99', got: %q", result)
	}
}

// TestGenerateWithComplexSegmentData tests segments with complex combined data
func TestGenerateWithComplexSegmentData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude 3.5 Sonnet"),
		createTestSegment(config.SegmentDirectory, "/Users/lang/coding/golang/src/cchline"),
		createTestSegment(config.SegmentGit, "main * +3 ~2 -1"),
		createTestSegment(config.SegmentContextWindow, "45.2% · 92.3k/200k tokens"),
		createTestSegment(config.SegmentCost, "$12.34"),
		createTestSegment(config.SegmentSession, "2h45m30s"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with complex segment data")
	}

	// Verify all complex data is present
	complexData := []string{
		"Claude 3.5 Sonnet",
		"/Users/lang/coding/golang/src/cchline",
		"main * +3 ~2 -1",
		"45.2%",
		"92.3k/200k tokens",
		"$12.34",
		"2h45m30s",
	}

	for _, data := range complexData {
		if !strings.Contains(result, data) {
			t.Errorf("Expected result to contain %q, got: %q", data, result)
		}
	}
}

// TestGenerateThemeConsistency tests that theme is consistently applied
func TestGenerateThemeConsistency(t *testing.T) {
	defaultCfg := createTestConfig(config.ThemeModeDefault, " | ")
	nerdFontCfg := createTestConfig(config.ThemeModeNerdFont, " | ")

	defaultGen := render.NewStatusLineGenerator(defaultCfg)
	nerdFontGen := render.NewStatusLineGenerator(nerdFontCfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
	}

	defaultResult := defaultGen.Generate(segments)
	nerdFontResult := nerdFontGen.Generate(segments)

	// Both should contain the segment data
	if !strings.Contains(defaultResult, "Sonnet 3.5") {
		t.Errorf("Expected default theme result to contain 'Sonnet 3.5', got: %q", defaultResult)
	}
	if !strings.Contains(nerdFontResult, "Sonnet 3.5") {
		t.Errorf("Expected nerd font theme result to contain 'Sonnet 3.5', got: %q", nerdFontResult)
	}

	// Results should be different due to different icons
	if defaultResult == nerdFontResult {
		t.Errorf("Expected different results for different themes, but got same: %q", defaultResult)
	}
}

// TestGenerateWithMixedSegmentTypes tests mixing different segment types
func TestGenerateWithMixedSegmentTypes(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentCCHModel, "claude-3-opus"),
		createTestSegment(config.SegmentDirectory, "myapp"),
		createTestSegment(config.SegmentCCHProvider, "anthropic"),
		createTestSegment(config.SegmentGit, "main *"),
		createTestSegment(config.SegmentCost, "$0.15"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for mixed segment types")
	}

	// Verify all segments are present
	expectedData := []string{
		"Sonnet 3.5",
		"claude-3-opus",
		"myapp",
		"anthropic",
		"main *",
		"$0.15",
	}

	for _, data := range expectedData {
		if !strings.Contains(result, data) {
			t.Errorf("Expected result to contain %q, got: %q", data, result)
		}
	}
}

// TestGenerateConsistency tests that Generate produces consistent results
func TestGenerateConsistency(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Sonnet 3.5"),
		createTestSegment(config.SegmentDirectory, "myapp"),
	}

	result1 := generator.Generate(segments)
	result2 := generator.Generate(segments)

	if result1 != result2 {
		t.Errorf("Expected consistent results, but got different outputs:\nFirst: %q\nSecond: %q", result1, result2)
	}
}

// TestGenerateWithDifferentConfigs tests Generate with different configurations
func TestGenerateWithDifferentConfigs(t *testing.T) {
	testCases := []struct {
		name      string
		theme     config.ThemeMode
		separator string
	}{
		{"Default theme with pipe separator", config.ThemeModeDefault, " | "},
		{"Nerd font theme with pipe separator", config.ThemeModeNerdFont, " | "},
		{"Default theme with arrow separator", config.ThemeModeDefault, " → "},
		{"Nerd font theme with dash separator", config.ThemeModeNerdFont, " - "},
		{"Default theme with no separator", config.ThemeModeDefault, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := createTestConfig(tc.theme, tc.separator)
			generator := render.NewStatusLineGenerator(cfg)

			segments := []segment.SegmentResult{
				createTestSegment(config.SegmentModel, "Sonnet 3.5"),
				createTestSegment(config.SegmentDirectory, "myapp"),
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			if !strings.Contains(result, "Sonnet 3.5") {
				t.Errorf("Expected result to contain 'Sonnet 3.5', got: %q", result)
			}
		})
	}
}

// TestGenerateWithSpecialSeparators tests Generate with special character separators
func TestGenerateWithSpecialSeparators(t *testing.T) {
	specialSeparators := []string{
		" → ",
		" • ",
		" ◆ ",
		" ▸ ",
		" ❯ ",
	}

	for _, sep := range specialSeparators {
		t.Run("separator: "+sep, func(t *testing.T) {
			cfg := createTestConfig(config.ThemeModeDefault, sep)
			generator := render.NewStatusLineGenerator(cfg)

			segments := []segment.SegmentResult{
				createTestSegment(config.SegmentModel, "A"),
				createTestSegment(config.SegmentDirectory, "B"),
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			if !strings.Contains(result, sep) {
				t.Errorf("Expected result to contain separator %q, got: %q", sep, result)
			}
		})
	}
}

// TestGenerateWithOnlyWhitespaceData tests segments with only whitespace
func TestGenerateWithOnlyWhitespaceData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "   "),
		createTestSegment(config.SegmentDirectory, "\t\t"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for whitespace-only data")
	}
}

// TestGenerateWithNumbersOnly tests segments with only numeric data
func TestGenerateWithNumbersOnly(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "12345"),
		createTestSegment(config.SegmentCCHRequests, "999"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for numeric-only data")
	}

	if !strings.Contains(result, "12345") || !strings.Contains(result, "999") {
		t.Errorf("Expected result to contain numeric data, got: %q", result)
	}
}

// TestGenerateWithSymbolsOnly tests segments with only symbols
func TestGenerateWithSymbolsOnly(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "* + ~ -"),
		createTestSegment(config.SegmentUsage, "↓ ↑"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for symbol-only data")
	}

	if !strings.Contains(result, "*") || !strings.Contains(result, "↓") {
		t.Errorf("Expected result to contain symbols, got: %q", result)
	}
}

// TestGenerateWithMixedCase tests segments with mixed case data
func TestGenerateWithMixedCase(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ClAuDe SoNnEt"),
		createTestSegment(config.SegmentDirectory, "MyApP/SrC"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for mixed case data")
	}

	if !strings.Contains(result, "ClAuDe SoNnEt") {
		t.Errorf("Expected result to preserve mixed case, got: %q", result)
	}
}

// TestGenerateWithRepeatedCharacters tests segments with repeated characters
func TestGenerateWithRepeatedCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "aaaaaaa"),
		createTestSegment(config.SegmentDirectory, "bbbbbbb"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for repeated characters")
	}

	if !strings.Contains(result, "aaaaaaa") {
		t.Errorf("Expected result to contain repeated characters, got: %q", result)
	}
}

// TestGenerateWithVeryLongString tests segments with very long strings
func TestGenerateWithVeryLongString(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	longString := strings.Repeat("a", 1000)
	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, longString),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for very long string")
	}

	if !strings.Contains(result, longString) {
		t.Errorf("Expected result to contain very long string")
	}
}

// TestGenerateWithVeryShortString tests segments with single character
func TestGenerateWithVeryShortString(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for single character data")
	}

	if !strings.Contains(result, "A") || !strings.Contains(result, "B") {
		t.Errorf("Expected result to contain single characters, got: %q", result)
	}
}

// TestGenerateWithPunctuation tests segments with punctuation
func TestGenerateWithPunctuation(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model: v3.5!"),
		createTestSegment(config.SegmentDirectory, "path/to/dir?"),
		createTestSegment(config.SegmentGit, "branch@main#1"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for punctuation data")
	}

	if !strings.Contains(result, "Model: v3.5!") {
		t.Errorf("Expected result to contain punctuation, got: %q", result)
	}
}

// TestGenerateWithMultipleSpaces tests segments with multiple consecutive spaces
func TestGenerateWithMultipleSpaces(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model     3.5"),
		createTestSegment(config.SegmentDirectory, "my     app"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for multiple spaces")
	}

	if !strings.Contains(result, "Model     3.5") {
		t.Errorf("Expected result to preserve multiple spaces, got: %q", result)
	}
}

// TestGenerateWithPathData tests segments with file path data
func TestGenerateWithPathData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentDirectory, "/home/user/projects/golang/src"),
		createTestSegment(config.SegmentGit, "feature/add-tests"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for path data")
	}

	if !strings.Contains(result, "/home/user/projects/golang/src") {
		t.Errorf("Expected result to contain path data, got: %q", result)
	}
}

// TestGenerateWithVersionData tests segments with version data
func TestGenerateWithVersionData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude 3.5.1"),
		createTestSegment(config.SegmentCCHModel, "claude-3-opus-20240229"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for version data")
	}

	if !strings.Contains(result, "Claude 3.5.1") {
		t.Errorf("Expected result to contain version data, got: %q", result)
	}
}

// TestGenerateWithRatioData tests segments with ratio/fraction data
func TestGenerateWithRatioData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentContextWindow, "50/100"),
		createTestSegment(config.SegmentCCHCost, "$5/$100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for ratio data")
	}

	if !strings.Contains(result, "50/100") {
		t.Errorf("Expected result to contain ratio data, got: %q", result)
	}
}

// TestGenerateWithDotNotation tests segments with dot notation
func TestGenerateWithDotNotation(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "claude.3.5.sonnet"),
		createTestSegment(config.SegmentCCHModel, "anthropic.claude.opus"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for dot notation")
	}

	if !strings.Contains(result, "claude.3.5.sonnet") {
		t.Errorf("Expected result to contain dot notation, got: %q", result)
	}
}

// TestGenerateWithHyphenatedData tests segments with hyphenated data
func TestGenerateWithHyphenatedData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "claude-3-5-sonnet"),
		createTestSegment(config.SegmentGit, "feature-branch-name"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for hyphenated data")
	}

	if !strings.Contains(result, "claude-3-5-sonnet") {
		t.Errorf("Expected result to contain hyphenated data, got: %q", result)
	}
}

// TestGenerateWithUnderscoreData tests segments with underscore data
func TestGenerateWithUnderscoreData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "claude_3_5_sonnet"),
		createTestSegment(config.SegmentDirectory, "my_project_dir"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for underscore data")
	}

	if !strings.Contains(result, "claude_3_5_sonnet") {
		t.Errorf("Expected result to contain underscore data, got: %q", result)
	}
}

// TestGenerateWithBracketData tests segments with bracket data
func TestGenerateWithBracketData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "[claude-3.5]"),
		createTestSegment(config.SegmentGit, "(main)"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for bracket data")
	}

	if !strings.Contains(result, "[claude-3.5]") {
		t.Errorf("Expected result to contain bracket data, got: %q", result)
	}
}

// TestGenerateWithArrowData tests segments with arrow symbols
func TestGenerateWithArrowData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentUsage, "↓100 ↑50"),
		createTestSegment(config.SegmentGit, "→ main"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for arrow data")
	}

	if !strings.Contains(result, "↓100") {
		t.Errorf("Expected result to contain arrow data, got: %q", result)
	}
}

// TestGenerateWithDotSeparatorInData tests segments with dot separator in data
func TestGenerateWithDotSeparatorInData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentContextWindow, "45.2% · 92.3k tokens"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for dot separator in data")
	}

	if !strings.Contains(result, "45.2% · 92.3k tokens") {
		t.Errorf("Expected result to contain dot separator, got: %q", result)
	}
}

// TestGenerateWithEqualsData tests segments with equals sign
func TestGenerateWithEqualsData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "model=claude-3.5"),
		createTestSegment(config.SegmentCost, "cost=$0.15"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for equals data")
	}

	if !strings.Contains(result, "model=claude-3.5") {
		t.Errorf("Expected result to contain equals data, got: %q", result)
	}
}

// TestGenerateWithColonData tests segments with colon data
func TestGenerateWithColonData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentSession, "1h:23m:45s"),
		createTestSegment(config.SegmentModel, "claude:3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for colon data")
	}

	if !strings.Contains(result, "1h:23m:45s") {
		t.Errorf("Expected result to contain colon data, got: %q", result)
	}
}

// TestGenerateWithCommaData tests segments with comma data
func TestGenerateWithCommaData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "$1,234.56"),
		createTestSegment(config.SegmentCCHRequests, "1,000 reqs"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for comma data")
	}

	if !strings.Contains(result, "$1,234.56") {
		t.Errorf("Expected result to contain comma data, got: %q", result)
	}
}

// TestGenerateWithParenthesesData tests segments with parentheses
func TestGenerateWithParenthesesData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude (3.5)"),
		createTestSegment(config.SegmentGit, "main (active)"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for parentheses data")
	}

	if !strings.Contains(result, "Claude (3.5)") {
		t.Errorf("Expected result to contain parentheses data, got: %q", result)
	}
}

// TestGenerateWithAsteriskData tests segments with asterisk data
func TestGenerateWithAsteriskData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "main *"),
		createTestSegment(config.SegmentModel, "Claude *3.5*"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for asterisk data")
	}

	if !strings.Contains(result, "main *") {
		t.Errorf("Expected result to contain asterisk data, got: %q", result)
	}
}

// TestGenerateWithPlusMinusData tests segments with plus and minus signs
func TestGenerateWithPlusMinusData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "+3 -2 ~1"),
		createTestSegment(config.SegmentCost, "+$0.10"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for plus/minus data")
	}

	if !strings.Contains(result, "+3 -2 ~1") {
		t.Errorf("Expected result to contain plus/minus data, got: %q", result)
	}
}

// TestGenerateWithSlashData tests segments with slash data
func TestGenerateWithSlashData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentDirectory, "/home/user/projects"),
		createTestSegment(config.SegmentContextWindow, "50/100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for slash data")
	}

	if !strings.Contains(result, "/home/user/projects") {
		t.Errorf("Expected result to contain slash data, got: %q", result)
	}
}

// TestGenerateWithBackslashData tests segments with backslash data
func TestGenerateWithBackslashData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentDirectory, "C:\\Users\\project"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for backslash data")
	}

	if !strings.Contains(result, "C:\\Users\\project") {
		t.Errorf("Expected result to contain backslash data, got: %q", result)
	}
}

// TestGenerateWithAtSignData tests segments with at sign
func TestGenerateWithAtSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "main@origin"),
		createTestSegment(config.SegmentModel, "claude@anthropic"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for at sign data")
	}

	if !strings.Contains(result, "main@origin") {
		t.Errorf("Expected result to contain at sign data, got: %q", result)
	}
}

// TestGenerateWithHashData tests segments with hash/pound sign
func TestGenerateWithHashData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "#123"),
		createTestSegment(config.SegmentModel, "claude#3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for hash data")
	}

	if !strings.Contains(result, "#123") {
		t.Errorf("Expected result to contain hash data, got: %q", result)
	}
}

// TestGenerateWithPercentSignData tests segments with percent sign
func TestGenerateWithPercentSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentContextWindow, "75%"),
		createTestSegment(config.SegmentUsage, "50% used"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for percent sign data")
	}

	if !strings.Contains(result, "75%") {
		t.Errorf("Expected result to contain percent sign data, got: %q", result)
	}
}

// TestGenerateWithAmpersandData tests segments with ampersand
func TestGenerateWithAmpersandData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude & Sonnet"),
		createTestSegment(config.SegmentGit, "main & develop"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for ampersand data")
	}

	if !strings.Contains(result, "Claude & Sonnet") {
		t.Errorf("Expected result to contain ampersand data, got: %q", result)
	}
}

// TestGenerateWithPipeInData tests segments with pipe character in data
func TestGenerateWithPipeInData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " :: ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude | Sonnet"),
		createTestSegment(config.SegmentGit, "main | develop"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for pipe in data")
	}

	if !strings.Contains(result, "Claude | Sonnet") {
		t.Errorf("Expected result to contain pipe in data, got: %q", result)
	}
}

// TestGenerateWithCaretData tests segments with caret
func TestGenerateWithCaretData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude^3.5"),
		createTestSegment(config.SegmentGit, "main^HEAD"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for caret data")
	}

	if !strings.Contains(result, "Claude^3.5") {
		t.Errorf("Expected result to contain caret data, got: %q", result)
	}
}

// TestGenerateWithTildeData tests segments with tilde
func TestGenerateWithTildeData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentDirectory, "~/projects"),
		createTestSegment(config.SegmentGit, "~1"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for tilde data")
	}

	if !strings.Contains(result, "~/projects") {
		t.Errorf("Expected result to contain tilde data, got: %q", result)
	}
}

// TestGenerateWithBacktickData tests segments with backtick
func TestGenerateWithBacktickData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "`claude-3.5`"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for backtick data")
	}

	if !strings.Contains(result, "`claude-3.5`") {
		t.Errorf("Expected result to contain backtick data, got: %q", result)
	}
}

// TestGenerateWithQuoteData tests segments with quotes
func TestGenerateWithQuoteData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "\"Claude 3.5\""),
		createTestSegment(config.SegmentDirectory, "'my-project'"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for quote data")
	}

	if !strings.Contains(result, "\"Claude 3.5\"") {
		t.Errorf("Expected result to contain quote data, got: %q", result)
	}
}

// TestGenerateWithExclamationData tests segments with exclamation mark
func TestGenerateWithExclamationData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude!"),
		createTestSegment(config.SegmentGit, "main!"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for exclamation data")
	}

	if !strings.Contains(result, "Claude!") {
		t.Errorf("Expected result to contain exclamation data, got: %q", result)
	}
}

// TestGenerateWithQuestionMarkData tests segments with question mark
func TestGenerateWithQuestionMarkData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude?"),
		createTestSegment(config.SegmentDirectory, "dir?"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for question mark data")
	}

	if !strings.Contains(result, "Claude?") {
		t.Errorf("Expected result to contain question mark data, got: %q", result)
	}
}

// TestGenerateWithDollarSignData tests segments with dollar sign
func TestGenerateWithDollarSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "$100"),
		createTestSegment(config.SegmentCCHCost, "$50/$200"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for dollar sign data")
	}

	if !strings.Contains(result, "$100") {
		t.Errorf("Expected result to contain dollar sign data, got: %q", result)
	}
}

// TestGenerateWithEuroSignData tests segments with euro sign
func TestGenerateWithEuroSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "€100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for euro sign data")
	}

	if !strings.Contains(result, "€100") {
		t.Errorf("Expected result to contain euro sign data, got: %q", result)
	}
}

// TestGenerateWithPoundSignData tests segments with pound sign
func TestGenerateWithPoundSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "£100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for pound sign data")
	}

	if !strings.Contains(result, "£100") {
		t.Errorf("Expected result to contain pound sign data, got: %q", result)
	}
}

// TestGenerateWithYenSignData tests segments with yen sign
func TestGenerateWithYenSignData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentCost, "¥100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for yen sign data")
	}

	if !strings.Contains(result, "¥100") {
		t.Errorf("Expected result to contain yen sign data, got: %q", result)
	}
}

// TestGenerateWithCopyrightSymbol tests segments with copyright symbol
func TestGenerateWithCopyrightSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude © 2024"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for copyright symbol")
	}

	if !strings.Contains(result, "©") {
		t.Errorf("Expected result to contain copyright symbol, got: %q", result)
	}
}

// TestGenerateWithRegisteredSymbol tests segments with registered symbol
func TestGenerateWithRegisteredSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude® 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for registered symbol")
	}

	if !strings.Contains(result, "®") {
		t.Errorf("Expected result to contain registered symbol, got: %q", result)
	}
}

// TestGenerateWithTrademarkSymbol tests segments with trademark symbol
func TestGenerateWithTrademarkSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude™"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for trademark symbol")
	}

	if !strings.Contains(result, "™") {
		t.Errorf("Expected result to contain trademark symbol, got: %q", result)
	}
}

// TestGenerateWithDegreeSymbol tests segments with degree symbol
func TestGenerateWithDegreeSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "90°"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for degree symbol")
	}

	if !strings.Contains(result, "°") {
		t.Errorf("Expected result to contain degree symbol, got: %q", result)
	}
}

// TestGenerateWithMicroSymbol tests segments with micro symbol
func TestGenerateWithMicroSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "100µs"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for micro symbol")
	}

	if !strings.Contains(result, "µ") {
		t.Errorf("Expected result to contain micro symbol, got: %q", result)
	}
}

// TestGenerateWithPilcrowSymbol tests segments with pilcrow symbol
func TestGenerateWithPilcrowSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "¶ Section"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for pilcrow symbol")
	}

	if !strings.Contains(result, "¶") {
		t.Errorf("Expected result to contain pilcrow symbol, got: %q", result)
	}
}

// TestGenerateWithSectionSymbol tests segments with section symbol
func TestGenerateWithSectionSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "§ 1.2.3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for section symbol")
	}

	if !strings.Contains(result, "§") {
		t.Errorf("Expected result to contain section symbol, got: %q", result)
	}
}

// TestGenerateWithDaggerSymbol tests segments with dagger symbol
func TestGenerateWithDaggerSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "†"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for dagger symbol")
	}

	if !strings.Contains(result, "†") {
		t.Errorf("Expected result to contain dagger symbol, got: %q", result)
	}
}

// TestGenerateWithDoubleDaggerSymbol tests segments with double dagger symbol
func TestGenerateWithDoubleDaggerSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "‡"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for double dagger symbol")
	}

	if !strings.Contains(result, "‡") {
		t.Errorf("Expected result to contain double dagger symbol, got: %q", result)
	}
}

// TestGenerateWithBulletSymbol tests segments with bullet symbol
func TestGenerateWithBulletSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "• Item"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for bullet symbol")
	}

	if !strings.Contains(result, "•") {
		t.Errorf("Expected result to contain bullet symbol, got: %q", result)
	}
}

// TestGenerateWithEllipsisSymbol tests segments with ellipsis symbol
func TestGenerateWithEllipsisSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Loading…"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for ellipsis symbol")
	}

	if !strings.Contains(result, "…") {
		t.Errorf("Expected result to contain ellipsis symbol, got: %q", result)
	}
}

// TestGenerateWithPrimeSymbol tests segments with prime symbol
func TestGenerateWithPrimeSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5′"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for prime symbol")
	}

	if !strings.Contains(result, "′") {
		t.Errorf("Expected result to contain prime symbol, got: %q", result)
	}
}

// TestGenerateWithDoublePrimeSymbol tests segments with double prime symbol
func TestGenerateWithDoublePrimeSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5″"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for double prime symbol")
	}

	if !strings.Contains(result, "″") {
		t.Errorf("Expected result to contain double prime symbol, got: %q", result)
	}
}

// TestGenerateWithLeftRightArrows tests segments with left-right arrows
func TestGenerateWithLeftRightArrows(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentUsage, "↔ bidirectional"),
		createTestSegment(config.SegmentGit, "← ↑ → ↓"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for left-right arrows")
	}

	if !strings.Contains(result, "↔") {
		t.Errorf("Expected result to contain left-right arrow, got: %q", result)
	}
}

// TestGenerateWithDoubleArrows tests segments with double arrows
func TestGenerateWithDoubleArrows(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "⇐ ⇑ ⇒ ⇓"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for double arrows")
	}

	if !strings.Contains(result, "⇒") {
		t.Errorf("Expected result to contain double arrow, got: %q", result)
	}
}

// TestGenerateWithMathSymbols tests segments with mathematical symbols
func TestGenerateWithMathSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∑ ∏ ∫"),
		createTestSegment(config.SegmentCost, "≈ $100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for math symbols")
	}

	if !strings.Contains(result, "∑") {
		t.Errorf("Expected result to contain math symbols, got: %q", result)
	}
}

// TestGenerateWithLogicalSymbols tests segments with logical symbols
func TestGenerateWithLogicalSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∧ ∨ ¬"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for logical symbols")
	}

	if !strings.Contains(result, "∧") {
		t.Errorf("Expected result to contain logical symbols, got: %q", result)
	}
}

// TestGenerateWithSetSymbols tests segments with set symbols
func TestGenerateWithSetSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∈ ∉ ⊂ ⊃"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for set symbols")
	}

	if !strings.Contains(result, "∈") {
		t.Errorf("Expected result to contain set symbols, got: %q", result)
	}
}

// TestGenerateWithGeometricShapes tests segments with geometric shapes
func TestGenerateWithGeometricShapes(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "■ ● ▲ ◆"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for geometric shapes")
	}

	if !strings.Contains(result, "■") {
		t.Errorf("Expected result to contain geometric shapes, got: %q", result)
	}
}

// TestGenerateWithCheckmarkAndX tests segments with checkmark and X
func TestGenerateWithCheckmarkAndX(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "✓ ✗"),
		createTestSegment(config.SegmentGit, "✔ ✘"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for checkmark and X")
	}

	if !strings.Contains(result, "✓") {
		t.Errorf("Expected result to contain checkmark, got: %q", result)
	}
}

// TestGenerateWithStarSymbols tests segments with star symbols
func TestGenerateWithStarSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "★ ☆ ⭐"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for star symbols")
	}

	if !strings.Contains(result, "★") {
		t.Errorf("Expected result to contain star symbols, got: %q", result)
	}
}

// TestGenerateWithHeartSymbols tests segments with heart symbols
func TestGenerateWithHeartSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "♥ ♡"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for heart symbols")
	}

	if !strings.Contains(result, "♥") {
		t.Errorf("Expected result to contain heart symbols, got: %q", result)
	}
}

// TestGenerateWithSuitSymbols tests segments with card suit symbols
func TestGenerateWithSuitSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "♠ ♣ ♥ ♦"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for suit symbols")
	}

	if !strings.Contains(result, "♠") {
		t.Errorf("Expected result to contain suit symbols, got: %q", result)
	}
}

// TestGenerateWithWeatherSymbols tests segments with weather symbols
func TestGenerateWithWeatherSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "☀ ☁ ☂ ☃"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for weather symbols")
	}

	if !strings.Contains(result, "☀") {
		t.Errorf("Expected result to contain weather symbols, got: %q", result)
	}
}

// TestGenerateWithMusicSymbols tests segments with music symbols
func TestGenerateWithMusicSymbols(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "♪ ♫ ♬"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for music symbols")
	}

	if !strings.Contains(result, "♪") {
		t.Errorf("Expected result to contain music symbols, got: %q", result)
	}
}

// TestGenerateWithFractionData tests segments with fraction data
func TestGenerateWithFractionData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentContextWindow, "½ ⅓ ¼"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for fraction data")
	}

	if !strings.Contains(result, "½") {
		t.Errorf("Expected result to contain fraction data, got: %q", result)
	}
}

// TestGenerateWithRomanNumerals tests segments with roman numerals
func TestGenerateWithRomanNumerals(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Ⅰ Ⅱ Ⅲ Ⅳ Ⅴ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for roman numerals")
	}

	if !strings.Contains(result, "Ⅰ") {
		t.Errorf("Expected result to contain roman numerals, got: %q", result)
	}
}

// TestGenerateWithCircledNumbers tests segments with circled numbers
func TestGenerateWithCircledNumbers(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "① ② ③ ④ ⑤"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for circled numbers")
	}

	if !strings.Contains(result, "①") {
		t.Errorf("Expected result to contain circled numbers, got: %q", result)
	}
}

// TestGenerateWithParenthesizedNumbers tests segments with parenthesized numbers
func TestGenerateWithParenthesizedNumbers(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⑴ ⑵ ⑶"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for parenthesized numbers")
	}

	if !strings.Contains(result, "⑴") {
		t.Errorf("Expected result to contain parenthesized numbers, got: %q", result)
	}
}

// TestGenerateWithFullwidthNumbers tests segments with fullwidth numbers
func TestGenerateWithFullwidthNumbers(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "０１２３４５"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for fullwidth numbers")
	}

	if !strings.Contains(result, "０") {
		t.Errorf("Expected result to contain fullwidth numbers, got: %q", result)
	}
}

// TestGenerateWithCombinedComplexData tests segments with all types of complex data combined
func TestGenerateWithCombinedComplexData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude 3.5 (Sonnet) ©2024"),
		createTestSegment(config.SegmentDirectory, "/home/user/my-project_v1.0"),
		createTestSegment(config.SegmentGit, "feature/add-tests * +5 -2 ~1"),
		createTestSegment(config.SegmentContextWindow, "45.2% · 92.3k/200k tokens"),
		createTestSegment(config.SegmentCost, "$12.34 ≈ €11.50"),
		createTestSegment(config.SegmentSession, "2h:45m:30s"),
		createTestSegment(config.SegmentCCHModel, "claude-3-opus-20240229"),
		createTestSegment(config.SegmentCCHCost, "$5.00/$100.00"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for combined complex data")
	}

	// Verify key elements are present
	keyElements := []string{
		"Claude 3.5",
		"/home/user/my-project_v1.0",
		"feature/add-tests",
		"45.2%",
		"92.3k/200k tokens",
		"$12.34",
		"2h:45m:30s",
		"claude-3-opus-20240229",
	}

	for _, elem := range keyElements {
		if !strings.Contains(result, elem) {
			t.Errorf("Expected result to contain %q, got: %q", elem, result)
		}
	}
}

// TestGenerateWithArabicNumerals tests segments with Arabic numerals
func TestGenerateWithArabicNumerals(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "٠١٢٣٤٥"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Arabic numerals")
	}

	if !strings.Contains(result, "٠") {
		t.Errorf("Expected result to contain Arabic numerals, got: %q", result)
	}
}

// TestGenerateWithDevanagariNumerals tests segments with Devanagari numerals
func TestGenerateWithDevanagariNumerals(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "०१२३४५"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Devanagari numerals")
	}

	if !strings.Contains(result, "०") {
		t.Errorf("Expected result to contain Devanagari numerals, got: %q", result)
	}
}

// TestGenerateWithChineseCharacters tests segments with Chinese characters
func TestGenerateWithChineseCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "模型 3.5"),
		createTestSegment(config.SegmentDirectory, "项目目录"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Chinese characters")
	}

	if !strings.Contains(result, "模型") {
		t.Errorf("Expected result to contain Chinese characters, got: %q", result)
	}
}

// TestGenerateWithJapaneseCharacters tests segments with Japanese characters
func TestGenerateWithJapaneseCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "モデル 3.5"),
		createTestSegment(config.SegmentDirectory, "プロジェクト"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Japanese characters")
	}

	if !strings.Contains(result, "モデル") {
		t.Errorf("Expected result to contain Japanese characters, got: %q", result)
	}
}

// TestGenerateWithKoreanCharacters tests segments with Korean characters
func TestGenerateWithKoreanCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "모델 3.5"),
		createTestSegment(config.SegmentDirectory, "프로젝트"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Korean characters")
	}

	if !strings.Contains(result, "모델") {
		t.Errorf("Expected result to contain Korean characters, got: %q", result)
	}
}

// TestGenerateWithThaiCharacters tests segments with Thai characters
func TestGenerateWithThaiCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "โมเดล 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Thai characters")
	}

	if !strings.Contains(result, "โมเดล") {
		t.Errorf("Expected result to contain Thai characters, got: %q", result)
	}
}

// TestGenerateWithHebrewCharacters tests segments with Hebrew characters
func TestGenerateWithHebrewCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "דגם 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Hebrew characters")
	}

	if !strings.Contains(result, "דגם") {
		t.Errorf("Expected result to contain Hebrew characters, got: %q", result)
	}
}

// TestGenerateWithArabicCharacters tests segments with Arabic characters
func TestGenerateWithArabicCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "نموذج 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Arabic characters")
	}

	if !strings.Contains(result, "نموذج") {
		t.Errorf("Expected result to contain Arabic characters, got: %q", result)
	}
}

// TestGenerateWithGreekCharacters tests segments with Greek characters
func TestGenerateWithGreekCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Μοντέλο 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Greek characters")
	}

	if !strings.Contains(result, "Μοντέλο") {
		t.Errorf("Expected result to contain Greek characters, got: %q", result)
	}
}

// TestGenerateWithCyrillicCharacters tests segments with Cyrillic characters
func TestGenerateWithCyrillicCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Модель 3.5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Cyrillic characters")
	}

	if !strings.Contains(result, "Модель") {
		t.Errorf("Expected result to contain Cyrillic characters, got: %q", result)
	}
}

// TestGenerateWithLatinExtended tests segments with Latin extended characters
func TestGenerateWithLatinExtended(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Café Naïve Résumé"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Latin extended characters")
	}

	if !strings.Contains(result, "Café") {
		t.Errorf("Expected result to contain Latin extended characters, got: %q", result)
	}
}

// TestGenerateWithCombiningDiacritics tests segments with combining diacritics
func TestGenerateWithCombiningDiacritics(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "e\u0301"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for combining diacritics")
	}
}

// TestGenerateWithZeroWidthCharacters tests segments with zero-width characters
func TestGenerateWithZeroWidthCharacters(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude\u200bSonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for zero-width characters")
	}
}

// TestGenerateWithRightToLeftMark tests segments with right-to-left mark
func TestGenerateWithRightToLeftMark(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude\u200fSonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for right-to-left mark")
	}
}

// TestGenerateWithLeftToRightMark tests segments with left-to-right mark
func TestGenerateWithLeftToRightMark(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude\u200eSonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for left-to-right mark")
	}
}

// TestGenerateWithNonBreakingSpace tests segments with non-breaking space
func TestGenerateWithNonBreakingSpace(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude\u00a0Sonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for non-breaking space")
	}

	if !strings.Contains(result, "Claude") {
		t.Errorf("Expected result to contain data with non-breaking space, got: %q", result)
	}
}

// TestGenerateWithEmDash tests segments with em dash
func TestGenerateWithEmDash(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude — Sonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for em dash")
	}

	if !strings.Contains(result, "—") {
		t.Errorf("Expected result to contain em dash, got: %q", result)
	}
}

// TestGenerateWithEnDash tests segments with en dash
func TestGenerateWithEnDash(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude – Sonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for en dash")
	}

	if !strings.Contains(result, "–") {
		t.Errorf("Expected result to contain en dash, got: %q", result)
	}
}

// TestGenerateWithHyphenMinus tests segments with hyphen-minus
func TestGenerateWithHyphenMinus(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude-Sonnet"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for hyphen-minus")
	}

	if !strings.Contains(result, "Claude-Sonnet") {
		t.Errorf("Expected result to contain hyphen-minus, got: %q", result)
	}
}

// TestGenerateWithMinusSign tests segments with minus sign
func TestGenerateWithMinusSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 − 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for minus sign")
	}

	if !strings.Contains(result, "−") {
		t.Errorf("Expected result to contain minus sign, got: %q", result)
	}
}

// TestGenerateWithPlusSign tests segments with plus sign
func TestGenerateWithPlusSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 + 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for plus sign")
	}

	if !strings.Contains(result, "+") {
		t.Errorf("Expected result to contain plus sign, got: %q", result)
	}
}

// TestGenerateWithMultiplicationSign tests segments with multiplication sign
func TestGenerateWithMultiplicationSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 × 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for multiplication sign")
	}

	if !strings.Contains(result, "×") {
		t.Errorf("Expected result to contain multiplication sign, got: %q", result)
	}
}

// TestGenerateWithDivisionSign tests segments with division sign
func TestGenerateWithDivisionSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 ÷ 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for division sign")
	}

	if !strings.Contains(result, "÷") {
		t.Errorf("Expected result to contain division sign, got: %q", result)
	}
}

// TestGenerateWithEqualsSign tests segments with equals sign
func TestGenerateWithEqualsSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 = 5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for equals sign")
	}

	if !strings.Contains(result, "=") {
		t.Errorf("Expected result to contain equals sign, got: %q", result)
	}
}

// TestGenerateWithNotEqualSign tests segments with not equal sign
func TestGenerateWithNotEqualSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 ≠ 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for not equal sign")
	}

	if !strings.Contains(result, "≠") {
		t.Errorf("Expected result to contain not equal sign, got: %q", result)
	}
}

// TestGenerateWithLessThanSign tests segments with less than sign
func TestGenerateWithLessThanSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "3 < 5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for less than sign")
	}

	if !strings.Contains(result, "<") {
		t.Errorf("Expected result to contain less than sign, got: %q", result)
	}
}

// TestGenerateWithGreaterThanSign tests segments with greater than sign
func TestGenerateWithGreaterThanSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 > 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for greater than sign")
	}

	if !strings.Contains(result, ">") {
		t.Errorf("Expected result to contain greater than sign, got: %q", result)
	}
}

// TestGenerateWithLessThanOrEqualSign tests segments with less than or equal sign
func TestGenerateWithLessThanOrEqualSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "3 ≤ 5"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for less than or equal sign")
	}

	if !strings.Contains(result, "≤") {
		t.Errorf("Expected result to contain less than or equal sign, got: %q", result)
	}
}

// TestGenerateWithGreaterThanOrEqualSign tests segments with greater than or equal sign
func TestGenerateWithGreaterThanOrEqualSign(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "5 ≥ 3"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for greater than or equal sign")
	}

	if !strings.Contains(result, "≥") {
		t.Errorf("Expected result to contain greater than or equal sign, got: %q", result)
	}
}

// TestGenerateWithInfinitySymbol tests segments with infinity symbol
func TestGenerateWithInfinitySymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∞"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for infinity symbol")
	}

	if !strings.Contains(result, "∞") {
		t.Errorf("Expected result to contain infinity symbol, got: %q", result)
	}
}

// TestGenerateWithPartialDifferentialSymbol tests segments with partial differential symbol
func TestGenerateWithPartialDifferentialSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∂"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for partial differential symbol")
	}

	if !strings.Contains(result, "∂") {
		t.Errorf("Expected result to contain partial differential symbol, got: %q", result)
	}
}

// TestGenerateWithNablaSymbol tests segments with nabla symbol
func TestGenerateWithNablaSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∇"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for nabla symbol")
	}

	if !strings.Contains(result, "∇") {
		t.Errorf("Expected result to contain nabla symbol, got: %q", result)
	}
}

// TestGenerateWithLaplacianSymbol tests segments with Laplacian symbol
func TestGenerateWithLaplacianSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∆"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Laplacian symbol")
	}

	if !strings.Contains(result, "∆") {
		t.Errorf("Expected result to contain Laplacian symbol, got: %q", result)
	}
}

// TestGenerateWithForAllSymbol tests segments with for all symbol
func TestGenerateWithForAllSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∀"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for for all symbol")
	}

	if !strings.Contains(result, "∀") {
		t.Errorf("Expected result to contain for all symbol, got: %q", result)
	}
}

// TestGenerateWithThereExistsSymbol tests segments with there exists symbol
func TestGenerateWithThereExistsSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∃"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for there exists symbol")
	}

	if !strings.Contains(result, "∃") {
		t.Errorf("Expected result to contain there exists symbol, got: %q", result)
	}
}

// TestGenerateWithEmptySetSymbol tests segments with empty set symbol
func TestGenerateWithEmptySetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∅"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for empty set symbol")
	}

	if !strings.Contains(result, "∅") {
		t.Errorf("Expected result to contain empty set symbol, got: %q", result)
	}
}

// TestGenerateWithUnionSymbol tests segments with union symbol
func TestGenerateWithUnionSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∪"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for union symbol")
	}

	if !strings.Contains(result, "∪") {
		t.Errorf("Expected result to contain union symbol, got: %q", result)
	}
}

// TestGenerateWithIntersectionSymbol tests segments with intersection symbol
func TestGenerateWithIntersectionSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "∩"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for intersection symbol")
	}

	if !strings.Contains(result, "∩") {
		t.Errorf("Expected result to contain intersection symbol, got: %q", result)
	}
}

// TestGenerateWithProperSubsetSymbol tests segments with proper subset symbol
func TestGenerateWithProperSubsetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊂"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for proper subset symbol")
	}

	if !strings.Contains(result, "⊂") {
		t.Errorf("Expected result to contain proper subset symbol, got: %q", result)
	}
}

// TestGenerateWithProperSupersetSymbol tests segments with proper superset symbol
func TestGenerateWithProperSupersetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊃"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for proper superset symbol")
	}

	if !strings.Contains(result, "⊃") {
		t.Errorf("Expected result to contain proper superset symbol, got: %q", result)
	}
}

// TestGenerateWithSubsetOrEqualSymbol tests segments with subset or equal symbol
func TestGenerateWithSubsetOrEqualSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊆"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for subset or equal symbol")
	}

	if !strings.Contains(result, "⊆") {
		t.Errorf("Expected result to contain subset or equal symbol, got: %q", result)
	}
}

// TestGenerateWithSupersetOrEqualSymbol tests segments with superset or equal symbol
func TestGenerateWithSupersetOrEqualSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊇"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for superset or equal symbol")
	}

	if !strings.Contains(result, "⊇") {
		t.Errorf("Expected result to contain superset or equal symbol, got: %q", result)
	}
}

// TestGenerateWithNotSubsetSymbol tests segments with not subset symbol
func TestGenerateWithNotSubsetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊄"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for not subset symbol")
	}

	if !strings.Contains(result, "⊄") {
		t.Errorf("Expected result to contain not subset symbol, got: %q", result)
	}
}

// TestGenerateWithNotSupersetSymbol tests segments with not superset symbol
func TestGenerateWithNotSupersetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊅"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for not superset symbol")
	}

	if !strings.Contains(result, "⊅") {
		t.Errorf("Expected result to contain not superset symbol, got: %q", result)
	}
}

// TestGenerateWithPowerSetSymbol tests segments with power set symbol
func TestGenerateWithPowerSetSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "℘"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for power set symbol")
	}

	if !strings.Contains(result, "℘") {
		t.Errorf("Expected result to contain power set symbol, got: %q", result)
	}
}

// TestGenerateWithAlephSymbol tests segments with aleph symbol
func TestGenerateWithAlephSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℵ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for aleph symbol")
	}

	if !strings.Contains(result, "ℵ") {
		t.Errorf("Expected result to contain aleph symbol, got: %q", result)
	}
}

// TestGenerateWithBethSymbol tests segments with beth symbol
func TestGenerateWithBethSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℶ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for beth symbol")
	}

	if !strings.Contains(result, "ℶ") {
		t.Errorf("Expected result to contain beth symbol, got: %q", result)
	}
}

// TestGenerateWithGimelSymbol tests segments with gimel symbol
func TestGenerateWithGimelSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℷ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for gimel symbol")
	}

	if !strings.Contains(result, "ℷ") {
		t.Errorf("Expected result to contain gimel symbol, got: %q", result)
	}
}

// TestGenerateWithDaletSymbol tests segments with dalet symbol
func TestGenerateWithDaletSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℸ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for dalet symbol")
	}

	if !strings.Contains(result, "ℸ") {
		t.Errorf("Expected result to contain dalet symbol, got: %q", result)
	}
}

// TestGenerateWithWeierstraussSymbol tests segments with Weierstrass symbol
func TestGenerateWithWeierstraussSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "℘"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for Weierstrass symbol")
	}

	if !strings.Contains(result, "℘") {
		t.Errorf("Expected result to contain Weierstrass symbol, got: %q", result)
	}
}

// TestGenerateWithImaginaryUnitSymbol tests segments with imaginary unit symbol
func TestGenerateWithImaginaryUnitSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℑ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for imaginary unit symbol")
	}

	if !strings.Contains(result, "ℑ") {
		t.Errorf("Expected result to contain imaginary unit symbol, got: %q", result)
	}
}

// TestGenerateWithRealPartSymbol tests segments with real part symbol
func TestGenerateWithRealPartSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "ℜ"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for real part symbol")
	}

	if !strings.Contains(result, "ℜ") {
		t.Errorf("Expected result to contain real part symbol, got: %q", result)
	}
}

// TestGenerateWithTradeSymbol tests segments with trade symbol
func TestGenerateWithTradeSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "℠"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for trade symbol")
	}

	if !strings.Contains(result, "℠") {
		t.Errorf("Expected result to contain trade symbol, got: %q", result)
	}
}

// TestGenerateWithServiceMarkSymbol tests segments with service mark symbol
func TestGenerateWithServiceMarkSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "℠"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for service mark symbol")
	}

	if !strings.Contains(result, "℠") {
		t.Errorf("Expected result to contain service mark symbol, got: %q", result)
	}
}

// TestGenerateWithOhmSymbol tests segments with ohm symbol
func TestGenerateWithOhmSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Ω"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for ohm symbol")
	}

	if !strings.Contains(result, "Ω") {
		t.Errorf("Expected result to contain ohm symbol, got: %q", result)
	}
}

// TestGenerateWithAngstromSymbol tests segments with angstrom symbol
func TestGenerateWithAngstromSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Å"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for angstrom symbol")
	}

	if !strings.Contains(result, "Å") {
		t.Errorf("Expected result to contain angstrom symbol, got: %q", result)
	}
}

// TestGenerateWithEstimateSymbol tests segments with estimate symbol
func TestGenerateWithEstimateSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "℮"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for estimate symbol")
	}

	if !strings.Contains(result, "℮") {
		t.Errorf("Expected result to contain estimate symbol, got: %q", result)
	}
}

// TestGenerateWithTurnstilesSymbol tests segments with turnstiles symbol
func TestGenerateWithTurnstilesSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊢"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for turnstiles symbol")
	}

	if !strings.Contains(result, "⊢") {
		t.Errorf("Expected result to contain turnstiles symbol, got: %q", result)
	}
}

// TestGenerateWithDoubleTurnstilesSymbol tests segments with double turnstiles symbol
func TestGenerateWithDoubleTurnstilesSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊨"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for double turnstiles symbol")
	}

	if !strings.Contains(result, "⊨") {
		t.Errorf("Expected result to contain double turnstiles symbol, got: %q", result)
	}
}

// TestGenerateWithBottomSymbol tests segments with bottom symbol
func TestGenerateWithBottomSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊥"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for bottom symbol")
	}

	if !strings.Contains(result, "⊥") {
		t.Errorf("Expected result to contain bottom symbol, got: %q", result)
	}
}

// TestGenerateWithTopSymbol tests segments with top symbol
func TestGenerateWithTopSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊤"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for top symbol")
	}

	if !strings.Contains(result, "⊤") {
		t.Errorf("Expected result to contain top symbol, got: %q", result)
	}
}

// TestGenerateWithUpTackSymbol tests segments with up tack symbol
func TestGenerateWithUpTackSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊢"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for up tack symbol")
	}

	if !strings.Contains(result, "⊢") {
		t.Errorf("Expected result to contain up tack symbol, got: %q", result)
	}
}

// TestGenerateWithDownTackSymbol tests segments with down tack symbol
func TestGenerateWithDownTackSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊣"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for down tack symbol")
	}

	if !strings.Contains(result, "⊣") {
		t.Errorf("Expected result to contain down tack symbol, got: %q", result)
	}
}

// TestGenerateWithLeftTackSymbol tests segments with left tack symbol
func TestGenerateWithLeftTackSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊣"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for left tack symbol")
	}

	if !strings.Contains(result, "⊣") {
		t.Errorf("Expected result to contain left tack symbol, got: %q", result)
	}
}

// TestGenerateWithRightTackSymbol tests segments with right tack symbol
func TestGenerateWithRightTackSymbol(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "⊢"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for right tack symbol")
	}

	if !strings.Contains(result, "⊢") {
		t.Errorf("Expected result to contain right tack symbol, got: %q", result)
	}
}

// TestGenerateWithFinalSummaryTest tests a final comprehensive scenario
func TestGenerateWithFinalSummaryTest(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude 3.5 Sonnet"),
		createTestSegment(config.SegmentDirectory, "/Users/lang/coding/golang/src/cchline"),
		createTestSegment(config.SegmentGit, "main * +5 -2"),
		createTestSegment(config.SegmentContextWindow, "50% · 100k tokens"),
		createTestSegment(config.SegmentCost, "$0.50"),
		createTestSegment(config.SegmentSession, "1h30m"),
		createTestSegment(config.SegmentCCHModel, "claude-3-opus"),
		createTestSegment(config.SegmentCCHCost, "$5/$100"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for final summary test")
	}

	// Verify the result contains all expected data
	expectedData := []string{
		"Claude 3.5 Sonnet",
		"/Users/lang/coding/golang/src/cchline",
		"main",
		"50%",
		"100k tokens",
		"$0.50",
		"1h30m",
		"claude-3-opus",
		"$5/$100",
	}

	for _, data := range expectedData {
		if !strings.Contains(result, data) {
			t.Errorf("Expected result to contain %q, got: %q", data, result)
		}
	}

	// Verify separators are present
	separatorCount := strings.Count(result, " | ")
	if separatorCount != 7 {
		t.Errorf("Expected 7 separators for 8 segments, but found %d", separatorCount)
	}
}

// TestGenerateWithBenchmarkScenario tests a realistic benchmark scenario
func TestGenerateWithBenchmarkScenario(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeNerdFont, " → ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Claude 3.5"),
		createTestSegment(config.SegmentDirectory, "project"),
		createTestSegment(config.SegmentGit, "develop"),
		createTestSegment(config.SegmentContextWindow, "25%"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for benchmark scenario")
	}

	// Verify all segments are present
	if !strings.Contains(result, "Claude 3.5") {
		t.Errorf("Expected result to contain model, got: %q", result)
	}
	if !strings.Contains(result, "project") {
		t.Errorf("Expected result to contain directory, got: %q", result)
	}
	if !strings.Contains(result, "develop") {
		t.Errorf("Expected result to contain git branch, got: %q", result)
	}
	if !strings.Contains(result, "25%") {
		t.Errorf("Expected result to contain context window, got: %q", result)
	}

	// Verify custom separator is used
	if !strings.Contains(result, " → ") {
		t.Errorf("Expected result to contain custom separator ' → ', got: %q", result)
	}
}

// TestGenerateWithMinimalSegments tests with minimal segment data
func TestGenerateWithMinimalSegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for minimal segments")
	}

	if !strings.Contains(result, "A") {
		t.Errorf("Expected result to contain 'A', got: %q", result)
	}
}

// TestGenerateWithMaximalSegments tests with maximum segment data
func TestGenerateWithMaximalSegments(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	// Create segments with all available types
	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model"),
		createTestSegment(config.SegmentDirectory, "Directory"),
		createTestSegment(config.SegmentGit, "Git"),
		createTestSegment(config.SegmentContextWindow, "Context"),
		createTestSegment(config.SegmentUsage, "Usage"),
		createTestSegment(config.SegmentCost, "Cost"),
		createTestSegment(config.SegmentSession, "Session"),
		createTestSegment(config.SegmentOutputStyle, "OutputStyle"),
		createTestSegment(config.SegmentUpdate, "Update"),
		createTestSegment(config.SegmentCCHModel, "CCHModel"),
		createTestSegment(config.SegmentCCHProvider, "CCHProvider"),
		createTestSegment(config.SegmentCCHCost, "CCHCost"),
		createTestSegment(config.SegmentCCHRequests, "CCHRequests"),
		createTestSegment(config.SegmentCCHLimits, "CCHLimits"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for maximal segments")
	}

	// Verify all segments are present
	expectedSegments := []string{
		"Model", "Directory", "Git", "Context", "Usage", "Cost", "Session",
		"OutputStyle", "Update", "CCHModel", "CCHProvider", "CCHCost",
		"CCHRequests", "CCHLimits",
	}

	for _, seg := range expectedSegments {
		if !strings.Contains(result, seg) {
			t.Errorf("Expected result to contain %q, got: %q", seg, result)
		}
	}

	// Verify correct number of separators (13 for 14 segments)
	separatorCount := strings.Count(result, " | ")
	if separatorCount != 13 {
		t.Errorf("Expected 13 separators for 14 segments, but found %d", separatorCount)
	}
}

// TestGenerateWithAlternatingThemes tests alternating between themes
func TestGenerateWithAlternatingThemes(t *testing.T) {
	defaultCfg := createTestConfig(config.ThemeModeDefault, " | ")
	nerdFontCfg := createTestConfig(config.ThemeModeNerdFont, " | ")

	defaultGen := render.NewStatusLineGenerator(defaultCfg)
	nerdFontGen := render.NewStatusLineGenerator(nerdFontCfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Test"),
		createTestSegment(config.SegmentDirectory, "Dir"),
	}

	result1 := defaultGen.Generate(segments)
	result2 := nerdFontGen.Generate(segments)
	result3 := defaultGen.Generate(segments)

	// First and third should be identical (same theme)
	if result1 != result3 {
		t.Errorf("Expected same results for same theme, got:\nFirst: %q\nThird: %q", result1, result3)
	}

	// First and second should be different (different themes)
	if result1 == result2 {
		t.Errorf("Expected different results for different themes, got same: %q", result1)
	}
}

// TestGenerateWithDifferentSeparatorLengths tests various separator lengths
func TestGenerateWithDifferentSeparatorLengths(t *testing.T) {
	separators := []string{
		"",
		" ",
		" | ",
		" :: ",
		" → ",
		" ◆ ",
		"     ",
	}

	for _, sep := range separators {
		t.Run("separator_length_"+string(rune(len(sep))), func(t *testing.T) {
			cfg := createTestConfig(config.ThemeModeDefault, sep)
			generator := render.NewStatusLineGenerator(cfg)

			segments := []segment.SegmentResult{
				createTestSegment(config.SegmentModel, "A"),
				createTestSegment(config.SegmentDirectory, "B"),
				createTestSegment(config.SegmentGit, "C"),
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			// Verify all data is present
			if !strings.Contains(result, "A") || !strings.Contains(result, "B") || !strings.Contains(result, "C") {
				t.Errorf("Expected result to contain all segments, got: %q", result)
			}
		})
	}
}

// TestGenerateWithRepeatedCalls tests repeated calls produce consistent results
func TestGenerateWithRepeatedCalls(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Test"),
		createTestSegment(config.SegmentDirectory, "Dir"),
	}

	results := make([]string, 5)
	for i := 0; i < 5; i++ {
		results[i] = generator.Generate(segments)
	}

	// All results should be identical
	for i := 1; i < 5; i++ {
		if results[i] != results[0] {
			t.Errorf("Expected consistent results, but result %d differs from result 0", i)
		}
	}
}

// TestGenerateWithDifferentGenerators tests different generator instances
func TestGenerateWithDifferentGenerators(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	gen1 := render.NewStatusLineGenerator(cfg)
	gen2 := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Test"),
	}

	result1 := gen1.Generate(segments)
	result2 := gen2.Generate(segments)

	if result1 != result2 {
		t.Errorf("Expected same results from different generators with same config, got:\nGen1: %q\nGen2: %q", result1, result2)
	}
}

// TestGenerateWithConfigVariations tests various config combinations
func TestGenerateWithConfigVariations(t *testing.T) {
	testCases := []struct {
		name      string
		theme     config.ThemeMode
		separator string
	}{
		{"Default theme, pipe separator", config.ThemeModeDefault, " | "},
		{"Default theme, arrow separator", config.ThemeModeDefault, " → "},
		{"Default theme, no separator", config.ThemeModeDefault, ""},
		{"Nerd font theme, pipe separator", config.ThemeModeNerdFont, " | "},
		{"Nerd font theme, arrow separator", config.ThemeModeNerdFont, " → "},
		{"Nerd font theme, no separator", config.ThemeModeNerdFont, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := createTestConfig(tc.theme, tc.separator)
			generator := render.NewStatusLineGenerator(cfg)

			segments := []segment.SegmentResult{
				createTestSegment(config.SegmentModel, "Test"),
				createTestSegment(config.SegmentDirectory, "Dir"),
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			if !strings.Contains(result, "Test") || !strings.Contains(result, "Dir") {
				t.Errorf("Expected result to contain all segments, got: %q", result)
			}
		})
	}
}

// TestGenerateWithSegmentDataVariations tests various segment data patterns
func TestGenerateWithSegmentDataVariations(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	testCases := []struct {
		name     string
		segments []segment.SegmentResult
	}{
		{
			"Simple alphanumeric",
			[]segment.SegmentResult{
				createTestSegment(config.SegmentModel, "abc123"),
			},
		},
		{
			"With special characters",
			[]segment.SegmentResult{
				createTestSegment(config.SegmentModel, "test-123_abc"),
			},
		},
		{
			"With spaces",
			[]segment.SegmentResult{
				createTestSegment(config.SegmentModel, "test data here"),
			},
		},
		{
			"With symbols",
			[]segment.SegmentResult{
				createTestSegment(config.SegmentModel, "$100 @ 50%"),
			},
		},
		{
			"With unicode",
			[]segment.SegmentResult{
				createTestSegment(config.SegmentModel, "测试 テスト"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := generator.Generate(tc.segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			// Verify the primary data is present
			if !strings.Contains(result, tc.segments[0].Data.Primary) {
				t.Errorf("Expected result to contain %q, got: %q", tc.segments[0].Data.Primary, result)
			}
		})
	}
}

// TestGenerateEdgeCaseEmptyConfig tests with minimal config
func TestGenerateEdgeCaseEmptyConfig(t *testing.T) {
	cfg := &config.SimpleConfig{
		Theme:     config.ThemeModeDefault,
		Separator: " | ",
	}
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Test"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with minimal config")
	}

	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain 'Test', got: %q", result)
	}
}

// TestGenerateWithNilSegmentData tests handling of segments with nil-like data
func TestGenerateWithNilSegmentData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		{
			ID: config.SegmentModel,
			Data: segment.SegmentData{
				Primary:   "Test",
				Secondary: "",
				Metadata:  nil,
			},
		},
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with nil metadata")
	}

	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain 'Test', got: %q", result)
	}
}

// TestGenerateWithEmptyMetadata tests segments with empty metadata
func TestGenerateWithEmptyMetadata(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		{
			ID: config.SegmentModel,
			Data: segment.SegmentData{
				Primary:   "Test",
				Secondary: "",
				Metadata:  make(map[string]string),
			},
		},
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with empty metadata")
	}

	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain 'Test', got: %q", result)
	}
}

// TestGenerateWithPopulatedMetadata tests segments with populated metadata
func TestGenerateWithPopulatedMetadata(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	metadata := make(map[string]string)
	metadata["key1"] = "value1"
	metadata["key2"] = "value2"

	segments := []segment.SegmentResult{
		{
			ID: config.SegmentModel,
			Data: segment.SegmentData{
				Primary:   "Test",
				Secondary: "Secondary",
				Metadata:  metadata,
			},
		},
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with populated metadata")
	}

	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain 'Test', got: %q", result)
	}
}

// TestGenerateWithSecondaryData tests segments with secondary data
func TestGenerateWithSecondaryData(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		{
			ID: config.SegmentModel,
			Data: segment.SegmentData{
				Primary:   "Primary",
				Secondary: "Secondary",
				Metadata:  make(map[string]string),
			},
		},
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with secondary data")
	}

	// Primary should be rendered, secondary may or may not be visible depending on implementation
	if !strings.Contains(result, "Primary") {
		t.Errorf("Expected result to contain 'Primary', got: %q", result)
	}
}

// TestGenerateStressTest tests with a large number of segments
func TestGenerateStressTest(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	// Create 100 segments
	segments := make([]segment.SegmentResult, 100)
	for i := 0; i < 100; i++ {
		segments[i] = createTestSegment(config.SegmentModel, "Segment")
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for stress test")
	}

	// Verify correct number of separators (99 for 100 segments)
	separatorCount := strings.Count(result, " | ")
	if separatorCount != 99 {
		t.Errorf("Expected 99 separators for 100 segments, but found %d", separatorCount)
	}
}

// TestGenerateWithVeryLongSeparator tests with very long separator
func TestGenerateWithVeryLongSeparator(t *testing.T) {
	longSeparator := strings.Repeat("=", 50)
	cfg := createTestConfig(config.ThemeModeDefault, longSeparator)
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "A"),
		createTestSegment(config.SegmentDirectory, "B"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result with very long separator")
	}

	if !strings.Contains(result, longSeparator) {
		t.Errorf("Expected result to contain long separator, got: %q", result)
	}
}

// TestGenerateWithSpecialCharacterSeparator tests with special character separator
func TestGenerateWithSpecialCharacterSeparator(t *testing.T) {
	specialSeparators := []string{
		"🔹",
		"▸",
		"❯",
		"➜",
		"⟶",
	}

	for _, sep := range specialSeparators {
		t.Run("separator_"+sep, func(t *testing.T) {
			cfg := createTestConfig(config.ThemeModeDefault, " "+sep+" ")
			generator := render.NewStatusLineGenerator(cfg)

			segments := []segment.SegmentResult{
				createTestSegment(config.SegmentModel, "A"),
				createTestSegment(config.SegmentDirectory, "B"),
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			if !strings.Contains(result, sep) {
				t.Errorf("Expected result to contain separator %q, got: %q", sep, result)
			}
		})
	}
}

// TestGenerateWithAllSegmentIDTypes tests all segment ID types
func TestGenerateWithAllSegmentIDTypes(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	allSegmentIDs := []config.SegmentID{
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

	for _, segID := range allSegmentIDs {
		t.Run(string(segID), func(t *testing.T) {
			segments := []segment.SegmentResult{
				{
					ID: segID,
					Data: segment.SegmentData{
						Primary:   "TestData",
						Secondary: "",
						Metadata:  make(map[string]string),
					},
				},
			}

			result := generator.Generate(segments)

			if result == "" {
				t.Errorf("Expected non-empty result for segment ID %q", segID)
			}

			if !strings.Contains(result, "TestData") {
				t.Errorf("Expected result to contain 'TestData' for segment ID %q, got: %q", segID, result)
			}
		})
	}
}

// TestGenerateWithRandomSegmentOrder tests segments in random order
func TestGenerateWithRandomSegmentOrder(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	// Create segments in different orders
	order1 := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model"),
		createTestSegment(config.SegmentDirectory, "Dir"),
		createTestSegment(config.SegmentGit, "Git"),
	}

	order2 := []segment.SegmentResult{
		createTestSegment(config.SegmentGit, "Git"),
		createTestSegment(config.SegmentModel, "Model"),
		createTestSegment(config.SegmentDirectory, "Dir"),
	}

	result1 := generator.Generate(order1)
	result2 := generator.Generate(order2)

	// Results should be different due to different order
	if result1 == result2 {
		t.Errorf("Expected different results for different segment orders, but got same: %q", result1)
	}

	// But both should contain all data
	for _, data := range []string{"Model", "Dir", "Git"} {
		if !strings.Contains(result1, data) || !strings.Contains(result2, data) {
			t.Errorf("Expected both results to contain %q", data)
		}
	}
}

// TestGenerateWithDuplicateSegmentIDs tests with duplicate segment IDs
func TestGenerateWithDuplicateSegmentIDs(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "First"),
		createTestSegment(config.SegmentModel, "Second"),
		createTestSegment(config.SegmentModel, "Third"),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for duplicate segment IDs")
	}

	// All data should be present
	if !strings.Contains(result, "First") || !strings.Contains(result, "Second") || !strings.Contains(result, "Third") {
		t.Errorf("Expected result to contain all data, got: %q", result)
	}

	// Should have 2 separators for 3 segments
	separatorCount := strings.Count(result, " | ")
	if separatorCount != 2 {
		t.Errorf("Expected 2 separators for 3 segments, but found %d", separatorCount)
	}
}

// TestGenerateWithMixedEmptyAndNonEmpty tests mix of empty and non-empty segments
func TestGenerateWithMixedEmptyAndNonEmpty(t *testing.T) {
	cfg := createTestConfig(config.ThemeModeDefault, " | ")
	generator := render.NewStatusLineGenerator(cfg)

	segments := []segment.SegmentResult{
		createTestSegment(config.SegmentModel, "Model"),
		createTestSegment(config.SegmentDirectory, ""),
		createTestSegment(config.SegmentGit, "Git"),
		createTestSegment(config.SegmentContextWindow, ""),
	}

	result := generator.Generate(segments)

	if result == "" {
		t.Error("Expected non-empty result for mixed empty and non-empty segments")
	}

	// Non-empty data should be present
	if !strings.Contains(result, "Model") || !strings.Contains(result, "Git") {
		t.Errorf("Expected result to contain non-empty data, got: %q", result)
	}
}

// TestGenerateWithThemeAndSeparatorCombinations tests all theme and separator combinations
func TestGenerateWithThemeAndSeparatorCombinations(t *testing.T) {
	themes := []config.ThemeMode{
		config.ThemeModeDefault,
		config.ThemeModeNerdFont,
	}

	separators := []string{
		" | ",
		" → ",
		" :: ",
		"",
	}

	for _, theme := range themes {
		for _, sep := range separators {
			t.Run(string(theme)+"_"+sep, func(t *testing.T) {
				cfg := createTestConfig(theme, sep)
				generator := render.NewStatusLineGenerator(cfg)

				segments := []segment.SegmentResult{
					createTestSegment(config.SegmentModel, "Test"),
				}

				result := generator.Generate(segments)

				if result == "" {
					t.Error("Expected non-empty result")
				}

				if !strings.Contains(result, "Test") {
					t.Errorf("Expected result to contain 'Test', got: %q", result)
				}
			})
		}
	}
}
