package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"theia/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main (){
	// Define Flag
	longFlag := flag.Bool("l", false, "Show detailed file info")
	allFlag := flag.Bool("a", false, "Shows all files")
	initFlag := flag.Bool("init", false, "Sets up shell integration")
	cdFlag := flag.Bool("cd", false, "Changes directory on selection")
	copyFlag := flag.Bool("cp", false, "Copies file path to clipboard")

	flag.Parse()

	cfg := tui.Config{
		ShowDetails: *longFlag,
		ShowHidden: *allFlag,
		CDMode: *cdFlag,
	}

	// For shell intigration 

	if *initFlag{
		fmt.Print(`
th() {
    local target 
    target=$(command theia "$@")

    if [ -z "$target" ]; then
        return 
    fi
    if [[ "$*" == *"-cd"* ]]; then
        if [ -d "$target" ]; then 
            cd "$target"
        fi 
    else 
        echo "$target"
    fi
}
`) 
		return 
	}
	//Determine the starting dir 
	startPath := "."
	args := flag.Args()
	if len(args) > 0{
		startPath = args[0]
	}
	absPath, err := filepath.Abs(startPath)
	if err != nil{
		fmt.Println("Error Starting: %v\n", err)
		os.Exit(1)
	}
	m, err := tui.InitialModel(absPath, cfg)
	if err != nil{
		fmt.Println("Error Starting: %v\n", err)
		os.Exit(1)
	}
	
	// Run bubbletea with output to Stderr

	prog := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	finalProg, err := prog.Run()
	if err != nil{
		fmt.Println("Error Starting: %v\n", err)
		os.Exit(1)
	}
	m, ok := finalProg.(tui.Model) 
	if ok && m.Selected != ""{
		if *copyFlag{
			command := exec.Command("pbcopy")
			stdin, _ := command.StdinPipe()
			go func(){
				defer stdin.Close()
				fmt.Fprint(stdin, m.Selected)
			}()
			command.Run()
			fmt.Fprintln(os.Stderr, "\n"+lipgloss.NewStyle().Bold(true).Render("Copied to clipboard:"), lipgloss.NewStyle().Underline(true).Render(m.Selected)+"\n")
		} else {
			fmt.Print(m.Selected)
		}
	}
	return 
}
