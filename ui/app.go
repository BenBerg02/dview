// Package ui adapts the connection screen, model and view to the Bubble Tea
// program interface.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"dview/db"
	"dview/ui/connect"
	"dview/ui/model"
	"dview/ui/view"
)

type App struct {
	w, h    int
	connect *connect.Model // non-nil until a connection is opened
	m       *model.Model
}

// New starts directly on an already opened database.
func New(d db.Database) (*App, error) {
	m, err := model.New(d)
	if err != nil {
		return nil, err
	}
	return &App{m: m}, nil
}

// NewConnect starts on the connection screen.
func NewConnect() *App {
	return &App{connect: connect.New()}
}

// NewConnectAt starts on the connection screen of the driver selected by a
// command line flag.
func NewConnectAt(flag string) *App {
	return &App{connect: connect.NewAt(flag)}
}

// Close releases the database, if one was opened.
func (a *App) Close() error {
	if a.m != nil {
		return a.m.DB.Close()
	}
	return nil
}

func (a *App) Init() tea.Cmd {
	if a.connect != nil {
		return a.connect.Init()
	}
	return a.m.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		a.w, a.h = ws.Width, ws.Height
	}
	if a.connect == nil {
		return a, a.m.Update(msg)
	}
	cmd := a.connect.Update(msg)
	if conn := a.connect.Conn; conn != nil {
		m, err := model.New(conn)
		if err != nil {
			conn.Close()
			a.connect.Conn, a.connect.Err = nil, err.Error()
			return a, cmd
		}
		a.m, a.connect = m, nil
		m.Update(tea.WindowSizeMsg{Width: a.w, Height: a.h})
		return a, tea.Batch(cmd, tea.DisableMouse)
	}
	return a, cmd
}

func (a *App) View() string {
	if a.connect != nil {
		return a.connect.View()
	}
	return view.Render(a.m)
}
