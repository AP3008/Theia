package tui

import (
	"fmt"
	"os"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss defined styles 
var (
	// header styling 
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		//Padding(0, 1).
		MarginBottom(1)
	
	// cursor style
	cursorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#EE6FF8")).
        Bold(true)
	
	// Directory Names
	dirStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#61AFEF"))
	
	// Symlink Names
	symlinkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#56B6C2"))
	
	// Executable name
	
	execStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#98C379"))

	// Regular File name
	
	regStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ABB2BF"))

	// Longview underlined 
	columnHeaderStyle = lipgloss.NewStyle().
        Underline(true).
        Bold(true).
        Foreground(lipgloss.Color("#ABB2BF"))

	// Design for colour button info
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#82827E"))
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
	header := headerStyle.Render(fmt.Sprintf("Exploring: " + tildaPath(m.Path)))
	s.WriteString(header + "\n")
	for i, file := range visibleFiles{
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex{
			cursor = cursorStyle.Render("> ")
		}
		name := regStyle.Render(file.Name)
		if file.IsDir{
			name = dirStyle.Render(file.Name + "/")
		}
		if file.IsSymLink{
			name = symlinkStyle.Render(name) 
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, name))
	}
	info := infoStyle.Render("\n [tab] enter directory [backspace] parent directory [enter] select  [q] quit\n")
	s.WriteString(info)
	return s.String()
}

// Handles -l flag 
func longView(m *Model, end int)string{
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Exlporing: %s\n\n", tildaPath(m.Path)))
	
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

// Private helper to remove home dir

func tildaPath(path string) string{
	home, err := os.UserHomeDir()
	if err != nil{
		return path 
	}
	if strings.HasPrefix(path, home){
		return strings.Replace(path, home, "~", 1)
	}
	return path
}
