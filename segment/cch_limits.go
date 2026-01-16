package segment

import (
	"fmt"
	"strings"

	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
)

type CCHLimitsSegment struct {
	Client *cch.Client
}

func (s *CCHLimitsSegment) Collect(input *config.InputData) SegmentData {
	if s.Client == nil {
		return SegmentData{}
	}

	stats, err := s.Client.GetStats()
	if err != nil {
		return SegmentData{}
	}

	// Format: 5h:$0/$5 W:$0/$50 M:$0/$100
	var parts []string

	if stats.Limit5h > 0 {
		parts = append(parts, fmt.Sprintf("5h:$%.0f", stats.Limit5h))
	}
	if stats.LimitWeekly > 0 {
		parts = append(parts, fmt.Sprintf("W:$%.0f", stats.LimitWeekly))
	}
	if stats.LimitMonthly > 0 {
		parts = append(parts, fmt.Sprintf("M:$%.0f", stats.LimitMonthly))
	}

	if len(parts) == 0 {
		return SegmentData{}
	}

	return SegmentData{
		Primary: strings.Join(parts, " "),
	}
}
