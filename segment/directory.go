package segment

import (
	"path/filepath"

	"github.com/WAY29/cchline/config"
)

type DirectorySegment struct{}

func (s *DirectorySegment) Collect(input *config.InputData) SegmentData {
	dir := input.Workspace.CurrentDir
	// 只取最后一级目录名
	name := filepath.Base(dir)
	return SegmentData{Primary: name}
}
