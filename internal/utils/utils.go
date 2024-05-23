package utils

import (
	"cambios/internal/types"
	"errors"
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
