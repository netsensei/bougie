package tui

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/tui/constants"
)

func Start(filepath string) error {
	f, err := tea.LogToFile("debug.log", "help")
	if err != nil {
		return fmt.Errorf("couldn't open a file for logging: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal(cerr)
		}
	}()

	m, _ := initBrowser(filepath)

	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	_, err = constants.P.Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	// If the final model contains an error in its status, return it so the
	// caller (main) can exit with a non-zero status and print a message.
	// if bm, ok := finalModel.(Browser); ok {
	// 	if bm.status.err != nil {
	// 		return bm.status.err
	// 	}
	// }

	return nil
}
