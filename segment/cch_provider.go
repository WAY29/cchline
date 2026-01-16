package segment

import (
	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
)

type CCHProviderSegment struct {
	Client *cch.Client
}

func (s *CCHProviderSegment) Collect(input *config.InputData) SegmentData {
	if s.Client == nil {
		return SegmentData{}
	}

	stats, err := s.Client.GetStats()
	if err != nil || stats.LastProviderName == "" {
		return SegmentData{}
	}

	return SegmentData{
		Primary: stats.LastProviderName,
	}
}
