// Contains the state that bubble tea requires that will define Theia's TUI
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"theia/filesystem"
)

// Defining the configs so I can add more flags easily without crowding the Model struct
type Config struct {
	ShowDetails bool
	ShowHidden  bool
	CDMode      bool
	FileMode    bool
	DirMode     bool
}
type Model struct {
	Path        string
	SystemFiles []filesystem.SystemFile
	Cursor      int
	Selected    string
	Settings    Config
	TopRow      int
	Height      int
	Searching   bool
	SearchTerm  string
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Declares the initial model state

func InitialModel(path string, configs Config) (Model, error) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	fs_list, err := filesystem.CreateSystemFileList(path, configs.ShowHidden, false, false)
	if err != nil {
		return Model{}, err
	}
	var selected_str string
	if len(fs_list) > 0 {
		selected_str = fs_list[0].Path
	}
	return Model{
		Path:        path,
		SystemFiles: fs_list,
		Cursor:      0,
		Selected:    selected_str,
		Settings:    configs,
		TopRow:      0,
		Height:      20,
		Searching: false,
		SearchTerm: "",
	}, nil
}
