// Contains the state that bubble tea requires that will define Theia's TUI
package tui

import (
	"theia/filesystem"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Path string 
	SystemFiles []filesystem.SystemFile
	Cursor int
	Selected string
	ShowDetails bool
}

func (m Model) Init() tea.Cmd{
	return nil
}

// Declares the initial model state

func InitialModel(path string, details bool) (Model, error){
	fs_list, err := filesystem.CreateSystemFileList(path)
	if err != nil{
		return Model{}, err
	}
	var selected_str string
	if len(fs_list) > 0 {
		selected_str = fs_list[0].Path
	}
	return Model {
		Path : path,
		SystemFiles : fs_list,
		Cursor : 0,
		Selected : selected_str,
		ShowDetails: details,
	}, nil
}
