package segment

import (
	"github.com/WAY29/cchline/config"
)

type OutputStyleSegment struct{}

func (s *OutputStyleSegment) Collect(input *config.InputData) SegmentData {
	style := input.OutputStyle.Name
	if style == "" {
		return SegmentData{}
	}
	return SegmentData{Primary: style}
}
