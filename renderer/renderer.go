package renderer

// Key represents a keyboard or mouse input
type Key int

// Define key constants that match goncurses keys for compatibility
const (
	KEY_MOUSE = 409 // Same as goncurses.KEY_MOUSE
	KEY_MOUSE_RELEASE = 410 // Custom key for mouse release events
)

// MouseEvent represents a mouse event
type MouseEvent struct {
	X, Y int
}

type Renderer interface {
	Beep()
	Draw(string)
	DrawAt(int, int, string)
	Dimensions() (y int, x int)
	Start()
	End()
	Refresh()
	BufferUpdate()
	Clear()
	GetChar() Key
	GetMouse() MouseEvent
	MouseSupport() bool
	CreateStatsWindow(height, width, y, x int) StatsWindow
	// Methods for rate control
	GetReadRate() int64
	GetBroadcastRate() int64
	SetRateChangeCallback(func(readRate, broadcastRate int64))
	SetInitialRates(readRate, broadcastRate int64)
}

// StatsWindow represents a window for displaying statistics
type StatsWindow interface {
	MovePrint(y, x int, str string)
	Clear()
	NoutRefresh()
	Delete() error
}
