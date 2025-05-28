package renderer

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	glog "github.com/ninjapanzer/gogol_channels/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// EbitenStatsWindow implements the StatsWindow interface using Ebiten
type EbitenStatsWindow struct {
	x, y, width, height int
	buffer              []string
	renderer            *EbitenRenderer
}

func (sw *EbitenStatsWindow) MovePrint(y, x int, str string) {
	if y >= len(sw.buffer) {
		// Resize buffer if needed
		newBuffer := make([]string, y+1)
		copy(newBuffer, sw.buffer)
		sw.buffer = newBuffer
	}
	sw.buffer[y] = str
}

func (sw *EbitenStatsWindow) Clear() {
	for i := range sw.buffer {
		sw.buffer[i] = ""
	}
}

func (sw *EbitenStatsWindow) NoutRefresh() {
	// Nothing to do, Ebiten will draw in the next frame
}

func (sw *EbitenStatsWindow) Delete() error {
	// Nothing to do, Ebiten handles cleanup
	return nil
}

// EbitenGame implements the ebiten.Game interface
type EbitenGame struct {
	renderer *EbitenRenderer
}

func (g *EbitenGame) Update() error {
	// Handle keyboard input
	g.renderer.charMutex.Lock()
	defer g.renderer.charMutex.Unlock()

	// Check for 'q' key press
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.renderer.charBuffer = append(g.renderer.charBuffer, Key('q'))
	}

	// Check for mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.renderer.mousePressed = true
		g.renderer.mouseX, g.renderer.mouseY = ebiten.CursorPosition()
		g.renderer.charBuffer = append(g.renderer.charBuffer, KEY_MOUSE)
	} else {
		g.renderer.mousePressed = false
	}

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.Black)

	// Draw the buffer
	for y, row := range g.renderer.buffer {
		for x, cell := range row {
			if cell != "" {
				text.Draw(screen, cell, g.renderer.fontFace, x*g.renderer.cellSize, (y+1)*g.renderer.cellSize, color.White)
			}
		}
	}

	// Draw stats windows
	for _, sw := range g.renderer.statsWindows {
		for y, line := range sw.buffer {
			if line != "" {
				text.Draw(screen, line, g.renderer.fontFace, sw.x, sw.y+(y+1)*g.renderer.cellSize, color.White)
			}
		}
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.width * g.renderer.cellSize, g.renderer.height * g.renderer.cellSize
}

// EbitenRenderer implements the Renderer interface using Ebiten
type EbitenRenderer struct {
	width, height int
	cellSize      int
	buffer        [][]string
	statsWindows  []*EbitenStatsWindow
	charBuffer    []Key
	charMutex     sync.Mutex
	mouseX, mouseY int
	mousePressed  bool
	game          *EbitenGame
	fontFace      font.Face
}

// NewEbitenRenderer creates a new Ebiten renderer
func NewEbitenRenderer(padding int) Renderer {
	// Get desktop dimensions
	desktopWidth, desktopHeight := ebiten.ScreenSizeInFullscreen()

	// Calculate 90% of desktop dimensions
	windowWidth := int(float64(desktopWidth) * 0.9)
	windowHeight := int(float64(desktopHeight) * 0.9)

	// Calculate cell size and grid dimensions to fit the window
	cellSize := 13 // Size of each cell in pixels
	width := windowWidth / cellSize
	height := windowHeight / cellSize

	// Log the dimensions
	glog.GetLogger().Info("Desktop dimensions", "width", desktopWidth, "height", desktopHeight)
	glog.GetLogger().Info("Window dimensions", "width", windowWidth, "height", windowHeight)
	glog.GetLogger().Info("Grid dimensions", "width", width, "height", height)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Game of Life - Ebiten")

	r := &EbitenRenderer{
		width:       width,
		height:      height,
		cellSize:    cellSize,
		buffer:      make([][]string, height),
		statsWindows: make([]*EbitenStatsWindow, 0),
		charBuffer:  make([]Key, 0),
		fontFace:    basicfont.Face7x13,
	}

	// Initialize buffer
	for i := range r.buffer {
		r.buffer[i] = make([]string, width)
	}

	r.game = &EbitenGame{renderer: r}
	r.Start()

	// Start the Ebiten game loop in a goroutine
	go func() {
		if err := ebiten.RunGame(r.game); err != nil {
			log.Fatal(err)
		}
	}()

	return r
}

func (r *EbitenRenderer) Start() {
	// Reset the buffer with the current dimensions
	r.buffer = make([][]string, r.height)
	for i := range r.buffer {
		r.buffer[i] = make([]string, r.width)
	}
	glog.GetLogger().Info("Starting Ebiten Window", "height", r.height, "width", r.width)
}

func (r *EbitenRenderer) End() {
	// Nothing to do, Ebiten handles cleanup
}

func (r *EbitenRenderer) Dimensions() (y int, x int) {
	return r.height, r.width
}

func (r *EbitenRenderer) Draw(str string) {
	// Not implemented for Ebiten
}

func (r *EbitenRenderer) DrawAt(y, x int, ach string) {
	if y >= 0 && y < len(r.buffer) && x >= 0 && x < len(r.buffer[y]) {
		r.buffer[y][x] = ach
	}
}

func (r *EbitenRenderer) Beep() {
	// Not implemented for Ebiten
}

func (r *EbitenRenderer) Refresh() {
	// Nothing to do, Ebiten refreshes automatically
}

func (r *EbitenRenderer) BufferUpdate() {
	// Nothing to do, Ebiten updates the buffer automatically
}

func (r *EbitenRenderer) Clear() {
	for i := range r.buffer {
		for j := range r.buffer[i] {
			r.buffer[i][j] = ""
		}
	}
}

func (r *EbitenRenderer) GetChar() Key {
	r.charMutex.Lock()
	defer r.charMutex.Unlock()

	if len(r.charBuffer) > 0 {
		ch := r.charBuffer[0]
		r.charBuffer = r.charBuffer[1:]
		return ch
	}
	return 0
}

func (r *EbitenRenderer) GetMouse() MouseEvent {
	return MouseEvent{
		X: r.mouseX / r.cellSize,
		Y: r.mouseY / r.cellSize,
	}
}

func (r *EbitenRenderer) MouseSupport() bool {
	return true
}

func (r *EbitenRenderer) CreateStatsWindow(height, width, y, x int) StatsWindow {
	sw := &EbitenStatsWindow{
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		buffer:   make([]string, height),
		renderer: r,
	}
	r.statsWindows = append(r.statsWindows, sw)
	return sw
}
