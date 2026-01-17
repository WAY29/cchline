package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
	"github.com/WAY29/cchline/render"
	"github.com/WAY29/cchline/segment"
	"github.com/WAY29/cchline/tui"
)

var (
	// Version 通过 ldflags 注入
	Version = "dev"

	configMode  = flag.Bool("c", false, "Open interactive configuration")
	versionMode = flag.Bool("v", false, "Show version")
	// CCH flags
	cchApiKey = flag.String("k", "", "CCH API key(Optional)")
	cchURL    = flag.String("u", "", "CCH server URL(Optional)")
)

func main() {
	flag.Parse()

	// Show version
	if *versionMode {
		fmt.Printf("cchline %s\n", Version)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Interactive configuration mode
	if *configMode {
		if err := tui.Run(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// No stdin data, show help or exit
		fmt.Printf("CCHLine %s - Claude Code Status Line\n", Version)
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  claude-code | cchline    Run status line")
		fmt.Println("  cchline -c               Open configuration")
		fmt.Println("  cchline -v               Show version")
		fmt.Println("  cchline -k KEY -u URL    Use CCH API")
		return
	}

	// Initialize CCH client if configured
	var cchClient *cch.Client
	apiKey := *cchApiKey
	url := *cchURL
	// Fallback to config file values
	if apiKey == "" {
		apiKey = cfg.CCHApiKey
	}
	if url == "" {
		url = cfg.CCHURL
	}
	if apiKey != "" && url != "" {
		cchClient = cch.NewClient(url, apiKey)
	}

	// Read stdin JSON
	reader := bufio.NewReader(os.Stdin)
	var input config.InputData
	if err := json.NewDecoder(reader).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Collect all segment data
	segments := collectAllSegments(cfg, &input, cchClient)

	// Generate status line
	generator := render.NewStatusLineGenerator(cfg)
	output := generator.Generate(segments)

	// Output to stdout
	fmt.Print(output)
}

func collectAllSegments(cfg *config.SimpleConfig, input *config.InputData, cchClient *cch.Client) []segment.SegmentResult {
	cfg.SegmentEnabled = config.NormalizeSegmentEnabled(cfg.SegmentOrder, cfg.SegmentEnabled)

	// Create a map of segment ID to collector
	collectorMap := map[string]struct {
		id      config.SegmentID
		collect func(*config.InputData) segment.SegmentData
	}{
		"model": {
			id:      config.SegmentModel,
			collect: (&segment.ModelSegment{}).Collect,
		},
		"directory": {
			id:      config.SegmentDirectory,
			collect: (&segment.DirectorySegment{}).Collect,
		},
		"git": {
			id:      config.SegmentGit,
			collect: (&segment.GitSegment{}).Collect,
		},
		"context_window": {
			id:      config.SegmentContextWindow,
			collect: (&segment.ContextWindowSegment{}).Collect,
		},
		"usage": {
			id:      config.SegmentUsage,
			collect: (&segment.UsageSegment{}).Collect,
		},
		"cost": {
			id:      config.SegmentCost,
			collect: (&segment.CostSegment{}).Collect,
		},
		"session": {
			id:      config.SegmentSession,
			collect: (&segment.SessionSegment{}).Collect,
		},
		"output_style": {
			id:      config.SegmentOutputStyle,
			collect: (&segment.OutputStyleSegment{}).Collect,
		},
		"update": {
			id:      config.SegmentUpdate,
			collect: (&segment.UpdateSegment{}).Collect,
		},
		// CCH Segments
		"cch_model": {
			id:      config.SegmentCCHModel,
			collect: (&segment.CCHModelSegment{Client: cchClient}).Collect,
		},
		"cch_provider": {
			id:      config.SegmentCCHProvider,
			collect: (&segment.CCHProviderSegment{Client: cchClient}).Collect,
		},
		"cch_cost": {
			id:      config.SegmentCCHCost,
			collect: (&segment.CCHCostSegment{Client: cchClient}).Collect,
		},
		"cch_requests": {
			id:      config.SegmentCCHRequests,
			collect: (&segment.CCHRequestsSegment{Client: cchClient}).Collect,
		},
		"cch_limits": {
			id:      config.SegmentCCHLimits,
			collect: (&segment.CCHLimitsSegment{Client: cchClient}).Collect,
		},
	}

	var results []segment.SegmentResult

	// Collect segments in the order specified by cfg.SegmentOrder
	enabledIdx := 0
	for _, name := range cfg.SegmentOrder {
		// Handle line break marker
		if name == config.LineBreakMarker {
			results = append(results, segment.SegmentResult{
				ID:   config.SegmentLineBreak,
				Data: segment.SegmentData{},
			})
			continue
		}

		instanceEnabled := true
		if enabledIdx < len(cfg.SegmentEnabled) {
			instanceEnabled = cfg.SegmentEnabled[enabledIdx]
		}
		enabledIdx++
		if !instanceEnabled {
			continue
		}

		collector, exists := collectorMap[name]
		if !exists {
			continue
		}

		data := collector.collect(input)
		if data.Primary != "" {
			results = append(results, segment.SegmentResult{
				ID:   collector.id,
				Data: data,
			})
		}
	}

	return results
}
