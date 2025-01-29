package game

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1, 0, 1)

	// Table formatting
	singleColumnStyle = borderStyle.Width(80)
	dualColumnStyle   = borderStyle.Width(39)

	// Text styles
	headerStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	attrNameStyle      = lipgloss.NewStyle().Bold(true).Foreground(white)
	attrTextValueStyle = lipgloss.NewStyle().Bold(false).Foreground(white)
	attrValueStyle     = lipgloss.NewStyle().Bold(false).Foreground(cyan)
	attrPosModStyle    = lipgloss.NewStyle().Bold(true).Foreground(green)
	attrNegModStyle    = lipgloss.NewStyle().Bold(true).Foreground(red)
	// blackText          = lipgloss.NewStyle().Foreground(black)
	// boldBlackText      = lipgloss.NewStyle().Foreground(black).Bold(true)
	// boldRedText        = lipgloss.NewStyle().Foreground(red).Bold(true)
	boldGreenText = lipgloss.NewStyle().Foreground(green).Bold(true)
	// boldYellowText     = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	// boldBlueText       = lipgloss.NewStyle().Foreground(blue).Bold(true)
	// boldMagentaText    = lipgloss.NewStyle().Foreground(magenta).Bold(true)
	// boldCyanText       = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	boldWhiteText = lipgloss.NewStyle().Foreground(white).Bold(true)
	redText       = lipgloss.NewStyle().Foreground(red)
	greenText     = lipgloss.NewStyle().Foreground(green)
	// yellowText         = lipgloss.NewStyle().Foreground(yellow)
	// blueText           = lipgloss.NewStyle().Foreground(blue)
	// magentaText        = lipgloss.NewStyle().Foreground(magenta)
	// cyanText           = lipgloss.NewStyle().Foreground(cyan)
	whiteText = lipgloss.NewStyle().Foreground(white)

	// Colors
	// black   = lipgloss.Color("0")
	red   = lipgloss.Color("1")
	green = lipgloss.Color("2")
	// yellow  = lipgloss.Color("3")
	// blue    = lipgloss.Color("4")
	// magenta = lipgloss.Color("5")
	cyan  = lipgloss.Color("6")
	white = lipgloss.Color("7")
)
