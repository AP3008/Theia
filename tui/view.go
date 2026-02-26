package tui

import(
	"fmt"
	"strings"
)

func (m Model) View() string{
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Exploring: %s\n", m.Path))
	for i, file := range m.SystemFiles{
		cursor := " "
		if m.Cursor == i{
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, file.Name))
	}
	s.WriteString("\n [tab] enter  [backspace] back  [enter] select  [q] quit\n")
	return s.String()
}

