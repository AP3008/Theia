// Handles logic for updating bubbletea
package tui

import (
	"theia/filesystem"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg)(tea.Model, tea.Cmd){
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0{
				m.Cursor--
			} else {
				m.Cursor = len(m.SystemFiles) - 1
			}
		case "down", "j":
			if m.Cursor < len(m.SystemFiles)-1{
				m.Cursor++
			} else {
				m.Cursor = 0 
			}
		case "tab":
			if len(m.SystemFiles) == 0 {
				return m, nil
			}
			curr := m.SystemFiles[m.Cursor]
			newFiles := m.SystemFiles
			if curr.IsDir{
				var err error
				newFiles, err = filesystem.CreateSystemFileList(curr.Path)	
				if err != nil{
					return m, nil
				}
			}
			m.Path = curr.Path
			m.SystemFiles = newFiles
			m.Cursor = 0
			if len(newFiles) > 0{
				m.Selected = newFiles[0].Path
			}
		}
	}
	return m, nil
}	

