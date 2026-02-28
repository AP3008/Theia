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
		// Padding(0, 1).
		MarginBottom(1)

	// cursor style
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c4a7e7")).
			Bold(true)

	// Directory Names
	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ccfd8"))

	// Symlink Names
	symlinkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56B6C2"))

	// Executable name

	execStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#908caa"))

	// Regular File name

	regStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0def4"))

	// Design for colour button info
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#82827E"))
)

func (m Model) View() string {
	var str string
	end := m.TopRow + m.Height
	end = min(end, len(m.SystemFiles))
	if m.Settings.ShowDetails {
		str = longView(&m, end)
	} else {
		str = normalView(&m, end)
	}
	return str
}

func normalView(m *Model, end int) string {
	var s strings.Builder
	visibleFiles := m.SystemFiles[m.TopRow:end]
	header := headerStyle.Render(fmt.Sprintf("  Exploring: " + tildaPath(m.Path)))
	s.WriteString(header + "\n")
	for i, file := range visibleFiles {
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex {
			cursor = cursorStyle.Render("> ")
		}
		name := regStyle.Render(file.Name)
		if file.IsDir {
			name = dirStyle.Render(file.Name + "/")
		}
		if file.IsSymLink {
			name = symlinkStyle.Render(name)
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, name))
	}
	info := infoStyle.Render("\n [tab] enter directory [backspace] parent directory [enter] select  [q] quit\n [ctrl+o] return current directory")
	s.WriteString(info)
	return s.String()
}

var (
	// Column headings
	columnHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Underline(true).
				Bold(true)

	// Rose Pine colours
	rose  = lipgloss.Color("#ebbcba")
	gold  = lipgloss.Color("#f6c177")
	iris  = lipgloss.Color("#c4a7e7")
	pine  = lipgloss.Color("#31748f")
	love  = lipgloss.Color("#eb6f92")
	muted = lipgloss.Color("#6e6a86")

	// Style Objects
	sizeStyle = lipgloss.NewStyle().Foreground(rose)
	dateStyle = lipgloss.NewStyle().Foreground(gold)
)

// Dimensions
const (
	modeW = 12
	sizeW = 10
	modW  = 16
)

// Handles -l flag
func longView(m *Model, end int) string {
	var s strings.Builder
	explorer := headerStyle.Render(fmt.Sprintf("  Exploring: " + tildaPath(m.Path)))
	s.WriteString(explorer + "\n")

	// Header
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		"  ",
		columnHeaderStyle.Width(modeW).Render("MODE"),
		columnHeaderStyle.Width(sizeW).Render("SIZE"),
		columnHeaderStyle.Width(modW).Render("MODIFIED"),
		columnHeaderStyle.Render("NAME"),
	)
	s.WriteString(header + "\n")

	visibleFiles := m.SystemFiles[m.TopRow:end]

	for i, file := range visibleFiles {
		actualIndex := i + m.TopRow
		cursor := " "
		if m.Cursor == actualIndex {
			cursor = cursorStyle.Render("> ")
		}

		// perm handling

		perm := permStyler(file.Permission.String())
		permBlock := lipgloss.NewStyle().Width(modeW).Render(perm)

		// size handling

		sizeStr := "-"
		if !file.IsDir {
			sizeStr = formatBytes(file.Size)
		}
		size := sizeStyle.Width(sizeW).Align(lipgloss.Right).Render(sizeStr)

		// Date handling

		date := dateStyle.Width(modW).Render(file.ModifiedTime.Format("Jan 02 15:04"))

		// Name handling

		var name string
		if file.IsDir {
			name = dirStyle.Render(file.Name + "/")
		} else {
			name = regStyle.Render(file.Name)
		}
		if file.IsSymLink {
			name = symlinkStyle.Render(name)
		}

		// line creation

		line := lipgloss.JoinHorizontal(lipgloss.Top,
			cursor,
			permBlock,
			" ",
			size,
			" ",
			date,
			" ",
			name,
		)
		s.WriteString(line + "\n")
	}
	info := infoStyle.Render("\n [tab] enter directory [backspace] parent directory [enter] select  [q] quit [ctrl+o] return this directory\n")
	s.WriteString(info)
	return s.String()
}

// Private helper to format bytes

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	divisor := int64(unit)
	exponent := 0
	for n := bytes / unit; n >= unit; n /= unit {
		divisor *= unit
		exponent++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(divisor), "KMGT"[exponent])
}

// Private helper to remove home dir

func tildaPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
	return path
}

// Private helper to style the permissions

func permStyler(perm string) string {
	var s strings.Builder
	for _, char := range perm {
		c := string(char)
		switch c {
		case "r":
			s.WriteString(lipgloss.NewStyle().Foreground(iris).Render(c))
		case "w":
			s.WriteString(lipgloss.NewStyle().Foreground(pine).Render(c))
		case "x":
			s.WriteString(lipgloss.NewStyle().Foreground(love).Render(c))
		case "d", "l":
			s.WriteString(lipgloss.NewStyle().Foreground(love).Render(c))
		default:
			s.WriteString(lipgloss.NewStyle().Foreground(muted).Render(c))
		}
	}
	return s.String()
}
