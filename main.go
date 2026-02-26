package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"theia/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main (){
	// Define Flag
	longFlag := flag.Bool("l", false, "Show detailed file info")
	allFlag := flag.Bool("a", false, "Shows all files")
	//editorFlag := flag.String("c", "nvim", "Opens code editor")
	flag.Parse()

	cfg := tui.Config{
		ShowDetails: *longFlag,
		ShowHidden: *allFlag,
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
		fmt.Print(m.Selected)
	}
}
