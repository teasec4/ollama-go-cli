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
	DefaultTerminalWidth  = 80
	DefaultMemoryGB       = 120
	MemoryBarLength       = 10
	ScrollOffsetPageSize  = 5
	ScrollOffsetLineSize  = 1
	MinimumMessageWidth   = 20
	InputMaxWidthOffset   = 6
	MessageWidthOffset    = 3
)

// Animation timing
const (
	AnimationFrameCount = 10
	ThinkingInterval    = 200 // milliseconds
)

// API endpoints
const (
	OllamaAPIChat = "/api/chat"
)

// Message roles
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Token estimation
const (
	CharactersPerToken   = 4
	TokensPerWordFactor  = 4
	TokensPerWordDivisor = 3 // ~1.3 tokens per word
)
