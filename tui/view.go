package tui

import(
	"fmt"
	"strings"
)

func (m Model) View() string{
	var str string
	end := m.TopRow + m.Height
	end = min(end, len(m.SystemFiles))
	if !m.Settings.ShowDetails{
		str = normalView(&m, end)
	}
	return str
}

func normalView(m *Model, end int) string{
	var s strings.Builder
	visibleFiles := m.SystemFiles[m.TopRow:end]
	s.WriteString(fmt.Sprintf("Exploring: %s\n", m.Path))
	for i, file := range visibleFiles{
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex{
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, file.Name))
	}
	s.WriteString("\n [tab] enter directory [backspace] parent directory [enter] select  [q] quit\n")
	return s.String()
}

func longView(m *Model, end int)string{
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Exlporing: %s\n", m.Path))

	return s.String()
}
