package segment

import (
	"fmt"

	"github.com/WAY29/cchline/config"
)

type CostSegment struct{}

func (s *CostSegment) Collect(input *config.InputData) SegmentData {
	cost := input.Cost.TotalCostUSD
	if cost == 0 {
		return SegmentData{}
	}
	return SegmentData{
		Primary: fmt.Sprintf("$%.2f", cost),
	}
}
