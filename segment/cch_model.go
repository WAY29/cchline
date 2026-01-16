package segment

import (
	"github.com/WAY29/cchline/cch"
	"github.com/WAY29/cchline/config"
)

type CCHModelSegment struct {
	Client *cch.Client
}

func (s *CCHModelSegment) Collect(input *config.InputData) SegmentData {
	if s.Client == nil {
		return SegmentData{}
	}

	stats, err := s.Client.GetStats()
	if err != nil || stats.LastUsedModel == "" {
		return SegmentData{}
	}

	return SegmentData{
		Primary: stats.LastUsedModel,
	}
}
