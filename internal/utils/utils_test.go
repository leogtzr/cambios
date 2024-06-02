package utils

import (
	"cambios/internal/types"
	"testing"
)

func TestGetMaxLengthLineInGitRepositoryOutput(t *testing.T) {
	t.Run("TestGetMaxLengthLineInGitRepositoryOutput", func(t *testing.T) {
		type testCase struct {
			gitStatusOutputLines []string
			expectedMaxLength    int
		}

		tests := []testCase{
			{[]string{
				"master",
				"leonardo",
				"leo",
			}, 8},
			{[]string{
				"abc",
				"xx",
				"",
			}, 3},
			{[]string{
				"",
			}, 0},
		}

		for _, test := range tests {
			if got := GetMaxLengthLineInGitRepositoryOutput(test.gitStatusOutputLines); got != test.expectedMaxLength {
				t.Errorf("GetMaxLengthLineInGitRepositoryOutput() = %d, want %d", got, test.expectedMaxLength)
			}
		}
	})
}

func TestGetGitStatusToken(t *testing.T) {
	t.Run("TestGetGitStatusToken", func(t *testing.T) {
		type expected struct {
			statusToken string
			file        string
			hasError    bool
		}

		type testCase struct {
			fileGitStatusOutput string
			expected            expected
		}

		tests := []testCase{
			{
				fileGitStatusOutput: " M .idea/runConfigurations/Tomcat_App.xml",
				expected: expected{
					statusToken: " M",
					file:        ".idea/runConfigurations/Tomcat_App.xml",
				},
			},
			{
				fileGitStatusOutput: "?? test-automation/src/test/resources/leogtzr.properties",
				expected: expected{
					statusToken: "??",
					file:        "test-automation/src/test/resources/leogtzr.properties",
				},
			},
			{
				fileGitStatusOutput: " D build.xml",
				expected: expected{
					statusToken: " D",
					file:        "build.xml",
				},
			},
		}

		for _, test := range tests {
			statusToken, fileName, err := GetGitStatusToken(GitStatusRegex, test.fileGitStatusOutput)
			if err != nil && !test.expected.hasError {
				t.Errorf("GetGitStatusToken() = (%v), want no error", err)
			}
			if statusToken != test.expected.statusToken {
				t.Errorf("got=(%s), want=(%s)", statusToken, test.expected.statusToken)
			}
			if fileName != test.expected.file {
				t.Errorf("got=(%s), want=(%s)", fileName, test.expected.file)
			}
		}

	})
}

func TestGetStatusRepositoryCounts(t *testing.T) {
	t.Run("TestGetStatusRepositoryCounts", func(t *testing.T) {
		type testCase struct {
			gitStatusCountLines []string
			hasError            bool
			expected            types.RepositoryStatusCount
		}

		tests := []testCase{
			{
				gitStatusCountLines: []string{
					"?? test-automation/src/test/resources/leogtzr.properties",
					" M .idea/runConfigurations/Tomcat_App.xml",
					" D build.xml",
					" D file.xml",
					" A dir1/file_a.txt",
				},
				hasError: false,
				expected: types.RepositoryStatusCount{
					Untracked: 1,
					Modified:  1,
					Deleted:   2,
					Added:     1,
				},
			},
		}

		for _, test := range tests {
			repoStatusCounts, err := RepoCounts(&test.gitStatusCountLines)
			if err != nil && !test.hasError {
				t.Errorf("GetStatusRepositoryCounts() = (%v), want no error", err)
			}

			if repoStatusCounts.Modified != test.expected.Modified {
				t.Errorf("got=(%v), want=(%v)", repoStatusCounts.Modified, test.expected.Modified)
			}

			if repoStatusCounts.Deleted != test.expected.Deleted {
				t.Errorf("got=(%v), want=(%v)", repoStatusCounts.Deleted, test.expected.Deleted)
			}

			if repoStatusCounts.Added != test.expected.Added {
				t.Errorf("got=(%v), want=(%v)", repoStatusCounts.Added, test.expected.Added)
			}

			if repoStatusCounts.Untracked != test.expected.Untracked {
				t.Errorf("got=%v, want=%v", repoStatusCounts.Untracked, test.expected.Untracked)
			}

			if repoStatusCounts.Untracked != test.expected.Untracked {
				t.Errorf("got=%v, want=%v", repoStatusCounts.Untracked, test.expected.Untracked)
			}
		}
	})
}

func TestGetStatusTextLegend(t *testing.T) {
	t.Run("TestGetStatusTextLegend", func(t *testing.T) {
		type testCase struct {
			repoStatusCount types.RepositoryStatusCount
			expected        string
		}

		tests := []testCase{
			{
				repoStatusCount: types.RepositoryStatusCount{
					Added:     3,
					Untracked: 2,
				},
				expected: "3 added, 2 untracked",
			},
			{
				repoStatusCount: types.RepositoryStatusCount{},
				expected:        "no changes",
			},
			{
				repoStatusCount: types.RepositoryStatusCount{
					Deleted: 5,
				},
				expected: "5 deleted",
			},
		}

		for _, test := range tests {
			if statusText := GetStatusTextLegend(&test.repoStatusCount); statusText != test.expected {
				t.Errorf("got=(%s), want=(%s)", statusText, test.expected)
			}
		}
	})
}
