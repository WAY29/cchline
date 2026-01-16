package segment

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/WAY29/cchline/config"
)

type GitSegment struct{}

func (s *GitSegment) Collect(input *config.InputData) SegmentData {
	dir := input.Workspace.CurrentDir

	branch := execGit(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if branch == "" {
		return SegmentData{} // 非 Git 仓库
	}

	status := execGit(dir, "status", "--porcelain")
	result := branch
	if status != "" {
		result += " *"
	}

	ahead, behind := getAheadBehind(dir)
	if ahead > 0 {
		result += fmt.Sprintf(" ↑%d", ahead)
	}
	if behind > 0 {
		result += fmt.Sprintf(" ↓%d", behind)
	}

	return SegmentData{Primary: result}
}

func execGit(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getAheadBehind(dir string) (ahead, behind int) {
	output := execGit(dir, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if output == "" {
		return 0, 0
	}
	parts := strings.Fields(output)
	if len(parts) == 2 {
		behind, _ = strconv.Atoi(parts[0])
		ahead, _ = strconv.Atoi(parts[1])
	}
	return
}
