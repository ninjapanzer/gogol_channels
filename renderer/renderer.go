package renderer

// Key represents a keyboard or mouse input
type Key int

// Define key constants that match goncurses keys for compatibility
const (
	KEY_MOUSE = 409 // Same as goncurses.KEY_MOUSE
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
}

// StatsWindow represents a window for displaying statistics
type StatsWindow interface {
	MovePrint(y, x int, str string)
	Clear()
	NoutRefresh()
	Delete() error
}
