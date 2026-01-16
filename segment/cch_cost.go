package segment

import (
	"fmt"

	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
)

type CCHCostSegment struct {
	Client *cch.Client
}

func (s *CCHCostSegment) Collect(input *config.InputData) SegmentData {
	if s.Client == nil {
		return SegmentData{}
	}

	stats, err := s.Client.GetStats()
	if err != nil {
		return SegmentData{}
	}

	// Format: $1.50/$10
	var primary string
	if stats.DailyQuota > 0 {
		primary = fmt.Sprintf("$%.2f/$%.0f", stats.TodayCost, stats.DailyQuota)
	} else {
		primary = fmt.Sprintf("$%.2f", stats.TodayCost)
	}

	return SegmentData{
		Primary: primary,
	}
}
