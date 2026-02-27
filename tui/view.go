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
func longView(m *Model, end int)string{
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Exlporing: %s\n", m.Path))
	
	// Header
	//Escape Code for underlined 
	u := "\033[4m"
	r := "\033[0m"

	// Column widths
	const (
		modeW = 10
		sizeW = 8
		modW  = 12
	)

	s.WriteString("  ") 
	
	s.WriteString(fmt.Sprintf("%s%s%s%*s ", u, "MODE", r, modeW-len("MODE"), ""))
	
	s.WriteString(fmt.Sprintf("%s%s%s%*s ", u, "SIZE", r, sizeW-len("SIZE"), ""))
	
	s.WriteString(fmt.Sprintf("%s%s%s%*s ", u, "MODIFIED", r, modW-len("MODIFIED"), ""))
	
	s.WriteString(fmt.Sprintf("%s%s%s\n", u, "NAME", r))

	visibleFiles := m.SystemFiles[m.TopRow:end]

	for i, file := range visibleFiles{
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex{
			cursor = ">"
		}

		sizeStr := "-"
		if !file.IsDir{
			sizeStr = formatBytes(file.Size)
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

// Private helper to format bytes

func formatBytes(bytes int64) string{
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	divisor := int64(unit)
	exponent := 0
	for n := bytes/unit; n>=unit; n/=unit{
		divisor *= unit
		exponent++
	}
	return fmt.Sprintf("%.1f %cB",float64(bytes)/float64(divisor), "KMGT"[exponent])
}
