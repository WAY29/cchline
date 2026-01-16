package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

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
		return
	}

	// Read stdin JSON
	reader := bufio.NewReader(os.Stdin)
	var input config.InputData
	if err := json.NewDecoder(reader).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Collect all segment data
	segments := collectAllSegments(cfg, &input)

	// Generate status line
	generator := render.NewStatusLineGenerator(cfg)
	output := generator.Generate(segments)

	// Output to stdout
	fmt.Print(output)
}

func collectAllSegments(cfg *config.SimpleConfig, input *config.InputData) []segment.SegmentResult {
	// Create a map of segment ID to collector
	collectorMap := map[string]struct {
		id      config.SegmentID
		enabled bool
		collect func(*config.InputData) segment.SegmentData
	}{
		"model": {
			id:      config.SegmentModel,
			enabled: cfg.Segments.Model,
			collect: (&segment.ModelSegment{}).Collect,
		},
		"directory": {
			id:      config.SegmentDirectory,
			enabled: cfg.Segments.Directory,
			collect: (&segment.DirectorySegment{}).Collect,
		},
		"git": {
			id:      config.SegmentGit,
			enabled: cfg.Segments.Git,
			collect: (&segment.GitSegment{}).Collect,
		},
		"context_window": {
			id:      config.SegmentContextWindow,
			enabled: cfg.Segments.ContextWindow,
			collect: (&segment.ContextWindowSegment{}).Collect,
		},
		"usage": {
			id:      config.SegmentUsage,
			enabled: cfg.Segments.Usage,
			collect: (&segment.UsageSegment{}).Collect,
		},
		"cost": {
			id:      config.SegmentCost,
			enabled: cfg.Segments.Cost,
			collect: (&segment.CostSegment{}).Collect,
		},
		"session": {
			id:      config.SegmentSession,
			enabled: cfg.Segments.Session,
			collect: (&segment.SessionSegment{}).Collect,
		},
		"output_style": {
			id:      config.SegmentOutputStyle,
			enabled: cfg.Segments.OutputStyle,
			collect: (&segment.OutputStyleSegment{}).Collect,
		},
		"update": {
			id:      config.SegmentUpdate,
			enabled: cfg.Segments.Update,
			collect: (&segment.UpdateSegment{}).Collect,
		},
	}

	var results []segment.SegmentResult

	// Collect segments in the order specified by cfg.SegmentOrder
	for _, name := range cfg.SegmentOrder {
		collector, exists := collectorMap[name]
		if !exists || !collector.enabled {
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
