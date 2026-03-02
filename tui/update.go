// Handles logic for updating bubbletea
package tui

import (
	"path/filepath"

	"theia/filesystem"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Searching{
			switch msg.String(){
				case "enter", "esc":
					//returns search goes (turns searching off)
					m.Searching = false
					m.SearchInput.Blur()
					if len(m.SystemFiles) == 0 || msg.String() == "esc"{
						originalList, err := filesystem.CreateSystemFileList(m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
						if err != nil{
							return m, tea.Quit
						}
						m.SystemFiles = originalList
					}
				default:
					// append this char to the search string 
					m.SearchInput, _ = m.SearchInput.Update(msg)
					fullList, err := filesystem.CreateSystemFileList(m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)

					if err != nil{
						return m, tea.Quit
					}
					newList := filesystem.SearchSystemList(m.SearchInput.Value(), fullList)
					if m.SearchInput.Value() == ""{
						m.SystemFiles = fullList
					} else {
						m.SystemFiles = newList
					}
					m.Cursor = 0
					m.TopRow = 0
			}		
		} else {
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.Selected = ""
				return m, tea.Quit
			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
					if m.Cursor < m.TopRow {
						m.TopRow--
					}
				} else {
					tr := max(0,len(m.SystemFiles) - 20)
					m.TopRow = tr
					m.Cursor = len(m.SystemFiles) - 1
				}
			case "down", "j":
				if m.Cursor < len(m.SystemFiles)-1 {
					m.Cursor++
					if m.Cursor >= m.TopRow+m.Height {
						m.TopRow++
					}
				} else {
					m.TopRow = 0
					m.Cursor = 0
				}

			case "tab":
				if len(m.SystemFiles) == 0 {
					return m, nil
				}
				curr := m.SystemFiles[m.Cursor]
				newFiles := m.SystemFiles
				if curr.IsDir {
					var err error
					newFiles, err = filesystem.CreateSystemFileList(curr.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
					if err != nil {
						return m, nil
					}

					m.Path = curr.Path
					m.SystemFiles = newFiles
					m.Cursor = 0
					m.TopRow = 0
					m.Selected = m.Path
				}

			case "d":
				m.Settings.DirMode = true
				m.Settings.FileMode = false
				newFiles, err := filesystem.CreateSystemFileList(m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
				if err != nil {
					return m, nil
				}
				m.SystemFiles = newFiles
				m.TopRow = 0 
				m.Cursor = 0

			case "f":
				m.Settings.DirMode = false
				m.Settings.FileMode = true
				newFiles, err := filesystem.CreateSystemFileList(m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
				if err != nil {
					return m, nil
				}
				m.SystemFiles = newFiles
				m.TopRow = 0 
				m.Cursor = 0

			case "n":
				m.Settings.DirMode = false
				m.Settings.FileMode = false
				newFiles, err := filesystem.CreateSystemFileList(m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
				if err != nil {
					return m, nil
				}
				m.SystemFiles = newFiles
				m.TopRow = 0 
				m.Cursor = 0

			case "backspace":
				parent := filepath.Dir(m.Path)
				newFiles, err := filesystem.CreateSystemFileList(parent, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode)
				if err != nil {
					return m, nil
				}
				m.SystemFiles = newFiles
				m.Path = parent
				m.Selected = m.Path
				m.Cursor = 0
				m.TopRow = 0
			case "enter":
				if len(m.SystemFiles) > 0 {
					m.Selected = m.SystemFiles[m.Cursor].Path
				}
				return m, tea.Quit

			case "/":
				m.Searching = true
				m.SearchInput.Focus()
				m.Cursor = 0
				m.TopRow = 0

			case "ctrl+o", "alt+enter":
				m.Selected = m.Path
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
