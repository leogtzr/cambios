package utils

import (
	"cambios/internal/types"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var GitStatusRegex = regexp.MustCompile(`^\s*([ MADRCU?]{2})\s+(.+)$`)

func GetGitStatusToken(re *regexp.Regexp, fileNameLine string) (string, string, error) {
	match := re.FindStringSubmatch(fileNameLine)

	if len(match) >= 2 {
		statusCount := match[1]
		fileName := match[2]

		return statusCount, fileName, nil
	}

	return "", "", errors.New("could not determine file name")
}

func GetMaxLengthLineInGitRepositoryOutput(lines []string) int {
	maxLength := 0

	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	return maxLength
}

func RepoCounts(menuItems *[]string) (types.RepositoryStatusCount, error) {
	repoStatusCounts := types.RepositoryStatusCount{}

	for _, item := range *menuItems {
		status, _, err := GetGitStatusToken(GitStatusRegex, item)
		if err != nil {
			return types.RepositoryStatusCount{}, err
		}

		status = strings.TrimSpace(status)
		switch {
		case status == "M":
			repoStatusCounts.Modified++
		case status == "??":
			repoStatusCounts.Untracked++
		case status == "A":
			repoStatusCounts.Added++
		case status == "D":
			repoStatusCounts.Deleted++
		case status == "R":
			repoStatusCounts.Renamed++
		}

	}

	return repoStatusCounts, nil
}

func GetStatusTextLegend(repoStatusCount *types.RepositoryStatusCount) string {
	var builder strings.Builder

	if repoStatusCount.Added > 0 {
		builder.WriteString(fmt.Sprintf("%d added", repoStatusCount.Added))
	}

	if (repoStatusCount.Deleted > 0) && (builder.Len() > 0) {
		builder.WriteString(fmt.Sprintf(", %d deleted", repoStatusCount.Deleted))
	} else if (repoStatusCount.Deleted > 0) && (builder.Len() == 0) {
		builder.WriteString(fmt.Sprintf("%d deleted", repoStatusCount.Deleted))
	}

	if (repoStatusCount.Modified > 0) && (builder.Len() > 0) {
		builder.WriteString(fmt.Sprintf(", %d modified", repoStatusCount.Modified))
	} else if (repoStatusCount.Modified > 0) && (builder.Len() == 0) {
		builder.WriteString(fmt.Sprintf("%d modified", repoStatusCount.Modified))
	}

	if (repoStatusCount.Untracked > 0) && (builder.Len() > 0) {
		builder.WriteString(fmt.Sprintf(", %d untracked", repoStatusCount.Untracked))
	} else if (repoStatusCount.Untracked > 0) && (builder.Len() == 0) {
		builder.WriteString(fmt.Sprintf("%d untracked", repoStatusCount.Untracked))
	}

	if builder.Len() == 0 {
		builder.WriteString("no changes")
	}

	return builder.String()
}
