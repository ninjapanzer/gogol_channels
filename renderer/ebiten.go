package renderer

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	// Check for window close
	if ebiten.IsWindowBeingClosed() {
		g.renderer.charBuffer = append(g.renderer.charBuffer, Key('q'))
		return ebiten.Termination
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

	// Define colors
	offWhite := color.RGBA{240, 240, 240, 255} // Off-white color for cubes
	blueLineColor := color.RGBA{0, 0, 255, 255} // Blue color for communication lines

	// Draw the buffer
	for y, row := range g.renderer.buffer {
		for x, cell := range row {
			if cell != "" {
				// Calculate cell position and size
				cellX := float64(x * g.renderer.cellSize)
				cellY := float64(y * g.renderer.cellSize)
				cellSize := float64(g.renderer.cellSize - 1) // Leave a small gap between cells

				// Define cell padding to make cells smaller and leave space for communication lines
				cellPadding := 3.0 // Padding around cells to make them smaller

				// Draw a cube (rectangle) for the cell if it's alive
				if cell == "0" {
					// Draw the cell as a smaller white square with padding
					ebitenutil.DrawRect(screen, cellX + cellPadding, cellY + cellPadding, 
						cellSize - (cellPadding * 2), cellSize - (cellPadding * 2), offWhite)
				}

				// Draw blue indicator for communication
				if g.renderer.communications[y][x] {
					// Check adjacent cells and draw blue lines between them if they're alive
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							// Skip the cell itself
							if dx == 0 && dy == 0 {
								continue
							}

							// Check if the adjacent cell is within bounds
							adjY, adjX := y+dy, x+dx
							if adjY >= 0 && adjY < len(g.renderer.buffer) && 
								adjX >= 0 && adjX < len(g.renderer.buffer[0]) {

								// Check if the adjacent cell is alive
								if g.renderer.buffer[adjY][adjX] == "0" {
									// Calculate start and end points for the line
									startX := cellX + cellSize/2
									startY := cellY + cellSize/2
									endX := float64(adjX * g.renderer.cellSize) + cellSize/2
									endY := float64(adjY * g.renderer.cellSize) + cellSize/2

									// Draw multiple lines to create a thicker appearance
									// Main line
									ebitenutil.DrawLine(screen, startX, startY, endX, endY, blueLineColor)
									// Additional lines with slight offsets to create thickness
									ebitenutil.DrawLine(screen, startX+1, startY, endX+1, endY, blueLineColor)
									ebitenutil.DrawLine(screen, startX-1, startY, endX-1, endY, blueLineColor)
									ebitenutil.DrawLine(screen, startX, startY+1, endX, endY+1, blueLineColor)
									ebitenutil.DrawLine(screen, startX, startY-1, endX, endY-1, blueLineColor)
								}
							}
						}
					}
				}
			}
		}
	}

	// Draw stats windows
	for _, sw := range g.renderer.statsWindows {
		// Calculate the background rectangle dimensions
		maxLineLength := 0
		lineCount := 0
		for _, line := range sw.buffer {
			if line != "" {
				lineCount++
				if len(line) > maxLineLength {
					maxLineLength = len(line)
				}
			}
		}

		// Add some padding around the text
		padding := 5
		bgWidth := maxLineLength*7 + padding*2 // Approximate width based on font size
		bgHeight := lineCount*g.renderer.cellSize + padding*2

		// Create a semi-transparent black background (20% alpha)
		bgColor := color.RGBA{0, 0, 0, 51} // 51 is 20% of 255

		// Draw the background
		ebitenutil.DrawRect(screen, float64(sw.x-padding), float64(sw.y-padding), float64(bgWidth), float64(bgHeight), bgColor)

		// Draw a border around the background
		borderColor := color.RGBA{255, 255, 255, 128}                                                                          // Semi-transparent white
		ebitenutil.DrawRect(screen, float64(sw.x-padding), float64(sw.y-padding), float64(bgWidth), 1, borderColor)            // Top
		ebitenutil.DrawRect(screen, float64(sw.x-padding), float64(sw.y-padding), 1, float64(bgHeight), borderColor)           // Left
		ebitenutil.DrawRect(screen, float64(sw.x-padding), float64(sw.y-padding+bgHeight-1), float64(bgWidth), 1, borderColor) // Bottom
		ebitenutil.DrawRect(screen, float64(sw.x-padding+bgWidth-1), float64(sw.y-padding), 1, float64(bgHeight), borderColor) // Right

		// Draw the text
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
	width, height  int
	cellSize       int
	buffer         [][]string
	communications [][]bool // Tracks which cells have communicated
	statsWindows   []*EbitenStatsWindow
	charBuffer     []Key
	charMutex      sync.Mutex
	mouseX, mouseY int
	mousePressed   bool
	game           *EbitenGame
	fontFace       font.Face
}

// NewEbitenRenderer creates a new Ebiten renderer
func NewEbitenRenderer(padding int) Renderer {
	// Get desktop dimensions
	desktopWidth, desktopHeight := ebiten.ScreenSizeInFullscreen()

	// Calculate 90% of desktop dimensions
	windowWidth := int(float64(desktopWidth) * 0.9)
	windowHeight := int(float64(desktopHeight) * 0.9)

	// Calculate cell size and grid dimensions to fit the window
	cellSize := 20 // Size of each cell in pixels
	width := windowWidth / cellSize
	height := windowHeight / cellSize

	// Log the dimensions
	glog.GetLogger().Info("Desktop dimensions", "width", desktopWidth, "height", desktopHeight)
	glog.GetLogger().Info("Window dimensions", "width", windowWidth, "height", windowHeight)
	glog.GetLogger().Info("Grid dimensions", "width", width, "height", height)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Game of Life - Ebiten")

	r := &EbitenRenderer{
		width:          width,
		height:         height,
		cellSize:       cellSize,
		buffer:         make([][]string, height),
		communications: make([][]bool, height),
		statsWindows:   make([]*EbitenStatsWindow, 0),
		charBuffer:     make([]Key, 0),
		fontFace:       basicfont.Face7x13,
	}

	// Initialize buffer and communications
	for i := range r.buffer {
		r.buffer[i] = make([]string, width)
		r.communications[i] = make([]bool, width)
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
	// Reset the buffer and communications with the current dimensions
	r.buffer = make([][]string, r.height)
	r.communications = make([][]bool, r.height)
	for i := range r.buffer {
		r.buffer[i] = make([]string, r.width)
		r.communications[i] = make([]bool, r.width)
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
		// Mark cell as communicating if it's alive (represented by "0")
		if ach == "0" {
			r.communications[y][x] = true
		} else {
			r.communications[y][x] = false
		}
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
			r.communications[i][j] = false
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
