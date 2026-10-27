package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"dview/db"
	"dview/ui"
)

const usage = "usage: dview [database.sqlite | --psql [connection-url]]"

func main() {
	app, err := start(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		if err == errUsage {
			os.Exit(2)
		}
		os.Exit(1)
	}

	_, err = tea.NewProgram(app, tea.WithAltScreen()).Run()
	app.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

var errUsage = fmt.Errorf("%s", usage)

// start builds the initial screen from the command line arguments.
func start(args []string) (*ui.App, error) {
	switch {
	case len(args) == 0:
		return ui.NewConnect(), nil
	case args[0] == "--psql" && len(args) == 1:
		return ui.NewConnectAt("psql"), nil
	case args[0] == "--psql" && len(args) == 2:
		p, err := db.OpenPostgres(args[1], "")
		if err != nil {
			return nil, err
		}
		return ui.New(p)
	case len(args) == 1 && args[0][0] != '-':
		database, err := db.Drivers[0].Open(args[0])
		if err != nil {
			return nil, err
		}
		return ui.New(database)
	}
	return nil, errUsage
}
