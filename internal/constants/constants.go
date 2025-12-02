package constants

// Animation frames for thinking indicator
var AnimationFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Color codes for lipgloss
const (
	ColorMagenta  = "205"
	ColorGreen    = "46"
	ColorCyan     = "51"
	ColorLimeGreen = "40"
	ColorYellow   = "226"
	ColorDarkGray = "240"
	ColorRed      = "196"
)

// UI dimensions
const (
	DefaultTerminalWidth = 80
	MemoryBarLength      = 10
	ScrollOffsetPageSize = 5
	ScrollOffsetLineSize = 1
	MinimumMessageWidth  = 20
	InputMaxWidthOffset  = 6
	MessageWidthOffset   = 3
)

// Animation timing
const (
	ThinkingInterval = 200 // milliseconds
)
