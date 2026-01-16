package segment

import (
	"github.com/WAY29/cchline/config"
)

type UpdateSegment struct{}

func (s *UpdateSegment) Collect(input *config.InputData) SegmentData {
	// 检查是否有新版本（可选实现）
	// 建议使用缓存避免频繁检查
	return SegmentData{}
}
