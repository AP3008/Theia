package main

import (
	"flag"
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"theia/tui"
)

func main (){
	// Define Flag
	longFlag := flag.Bool("l", false, "Show detailed file info")
	flag.Parse()

	//Determine the starting dir 
	startPath := "."
	args := flag.Args()
	if len(args) > 0{
		startPath = args[0]
	}
	m, err := tui.InitialModel(startPath, *longFlag)
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
