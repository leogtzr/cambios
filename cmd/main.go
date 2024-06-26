package main

import (
	"cambios/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getGitStatus(repoDirectory string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoDirectory, "status", "--short")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

func writeCommandStatusToFile(command, filePath, repositoryPath string) error {
	file, err := os.Create("/tmp/cmd.out")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s|%s|%s", command, repositoryPath, filePath))
	if err != nil {
		return err
	}

	return nil
}

func handleFileCommand(command, fileNameLine, repositoryPath string) error {
	_, fileName, err := utils.GetGitStatusToken(utils.GitStatusRegex, fileNameLine)
	if err != nil {
		return err
	}

	return writeCommandStatusToFile(command, fileName, repositoryPath)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <repository path>\n", os.Args[0])

		os.Exit(1)
	}

	repositoryPath := os.Args[1]

	menuItems, err := getGitStatus(repositoryPath)
	if err != nil {
		os.Exit(1)
	}

	repoStatusCounts, err := utils.RepoCounts(&menuItems)
	if err != nil {
		ui.Close()
		fmt.Fprintln(os.Stderr, "error: parsing Git output")
	}

	if (repoStatusCounts.Added == 0) && (repoStatusCounts.Deleted == 0) &&
		(repoStatusCounts.Modified == 0) && (repoStatusCounts.Renamed == 0) &&
		(repoStatusCounts.Untracked == 0) {
		fmt.Println("no changes in repo")
		os.Exit(0)
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		ui.Close()
		os.Exit(0)
	}()

	statusText := utils.GetStatusTextLegend(&repoStatusCounts)

	statusWidget := widgets.NewParagraph()
	statusWidget.Text = statusText
	statusWidget.TextStyle = ui.NewStyle(ui.ColorCyan, ui.ColorClear, ui.ModifierBold)
	statusWidget.SetRect(0, 0, len(statusText)+3, 3)
	statusWidget.BorderStyle.Fg = ui.ColorBlack
	statusWidget.BorderStyle.Bg = ui.ColorBlue
	statusWidget.BorderStyle.Modifier = ui.ModifierBold

	maxLength := utils.GetMaxLengthLineInGitRepositoryOutput(menuItems) + 1

	list := widgets.NewList()
	list.Title = "Repository"
	list.Rows = menuItems
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false
	list.SetRect(0, 3, maxLength+2, len(menuItems)+5)

	ui.Render(statusWidget, list)

	currentRow := 0
	maxRow := len(menuItems) - 1

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			if currentRow < maxRow {
				currentRow++
				list.ScrollDown()
			} else {
				currentRow = 0
				list.SelectedRow = 0
				list.ScrollTop()
			}
		case "k", "<Up>":
			if currentRow > 0 {
				currentRow--
				list.ScrollUp()
			} else {
				currentRow = maxRow
				list.SelectedRow = maxRow
				list.ScrollBottom()
			}

		case "<Enter>":
			if err := handleFileCommand("clipboard", menuItems[currentRow], repositoryPath); err != nil {
				os.Exit(1)
			}

			ui.Close()
			os.Exit(0)
		case "v": // view file
			if err := handleFileCommand("v", menuItems[currentRow], repositoryPath); err != nil {
				os.Exit(1)
			}

			ui.Close()
			os.Exit(0)
		case "d", "f": // diff
			if err := handleFileCommand("diff", menuItems[currentRow], repositoryPath); err != nil {
				os.Exit(1)
			}

			ui.Close()
			os.Exit(0)
		case "<Escape>":
			return
		}

		list.SelectedRow = currentRow
		ui.Render(statusWidget, list)
	}
}
