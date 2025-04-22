package tui

import "github.com/charmbracelet/lipgloss"

// red   = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
// green = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
// blue  = lipgloss.AdaptiveColor{Light: "#7246FF", Dark: "#7A56ED"}
var (
	ColorRed    = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"} // lipgloss.Color("#f54242")
	ColorYellow = lipgloss.Color("#b0ad09")
	ColorBlue   = lipgloss.Color("#347aeb")
	ColorGray   = lipgloss.Color("#636363")
	ColorGreen  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"} // lipgloss.Color("#1fb009")
	ColorWhite  = lipgloss.Color("#FFFDF5")
)
