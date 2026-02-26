package tui

import(
	"fmt"
	"strings"
)

func (m Model) View() string{
	var str string
	end := m.TopRow + m.Height
	end = min(end, len(m.SystemFiles))
	if m.Settings.ShowDetails{
		str = longView(&m, end)
	} else {
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

// Handles -l flag 
// TODO: ADD BYTE CONVERSION PRIVATE FUNC & BETTER FORMATTING W/ SPACING
func longView(m *Model, end int)string{
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Exlporing: %s\n", m.Path))
	
	// Header
	s.WriteString(fmt.Sprintf("  %-10s %-8s %-12s %s\n", "MODE", "SIZE", "MODIFIED", "NAME"))
    s.WriteString(strings.Repeat("-", 50) + "\n")

	visibleFiles := m.SystemFiles[m.TopRow:end]

	for i, file := range visibleFiles{
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex{
			cursor = ">"
		}

		sizeStr := fmt.Sprintf("%d B", file.Size)
		if file.IsDir{
			sizeStr = "-"
		}
		line := fmt.Sprintf("%s %-10s %8s %-12s %s\n",
            cursor,
            file.Permission.String(), 
            sizeStr,
            file.ModifiedTime.Format("Jan 02 15:04"), 
            file.Name,
        )
        s.WriteString(line)
	}
	s.WriteString("\n [tab] enter directory [backspace] parent directory [enter] select  [q] quit\n")
	return s.String()
}
