package segment

import (
	"fmt"

	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
)

type CCHRequestsSegment struct {
	Client *cch.Client
}

func (s *CCHRequestsSegment) Collect(input *config.InputData) SegmentData {
	if s.Client == nil {
		return SegmentData{}
	}

	stats, err := s.Client.GetStats()
	if err != nil {
		return SegmentData{}
	}

	if stats.TodayRequests == 0 {
		return SegmentData{}
	}

	return SegmentData{
		Primary: fmt.Sprintf("%d reqs", stats.TodayRequests),
	}
}
