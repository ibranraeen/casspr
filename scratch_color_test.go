package main

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
)

func main() {
	c1 := lipgloss.Color("#60A5FA")
	c2 := lipgloss.Color("#2563EB")
	fmt.Printf("COLORTERM: %q\n", os.Getenv("COLORTERM"))
	fmt.Printf("TERM: %q\n", os.Getenv("TERM"))
	fmt.Printf("Color profile: %v\n", lipgloss.ColorProfile())
	
	style1 := lipgloss.NewStyle().Foreground(c1)
	style2 := lipgloss.NewStyle().Foreground(c2)
	
	fmt.Println(style1.Render("This is starlightBlue (#60A5FA)"))
	fmt.Println(style2.Render("This is celestialBlue (#2563EB)"))
}
