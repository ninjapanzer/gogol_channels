package renderer

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"strconv"
	"sync"
	"time"

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

	// Update fading cells
	now := time.Now()
	delta := now.Sub(g.renderer.lastFrameTime).Seconds()
	g.renderer.lastFrameTime = now

	// Fade rate: how much opacity decreases per second
	fadeRate := 1.0 // Reduced fade rate to make the fading effect more visible

	// Update fading cells
	for y := range g.renderer.fadingCells {
		for x := range g.renderer.fadingCells[y] {
			if g.renderer.fadingCells[y][x] > 0 {
				// Reduce opacity based on time delta
				g.renderer.fadingCells[y][x] -= float64(fadeRate) * delta

				// Ensure opacity doesn't go below 0
				if g.renderer.fadingCells[y][x] < 0 {
					g.renderer.fadingCells[y][x] = 0
				}
			}
		}
	}

	// Check for 'q' key press
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.renderer.charBuffer = append(g.renderer.charBuffer, Key('q'))
	}

	// Check for window close
	if ebiten.IsWindowBeingClosed() {
		g.renderer.charBuffer = append(g.renderer.charBuffer, Key('q'))
		return ebiten.Termination
	}

	// Slider constants
	sliderX := 50
	sliderY := 30
	sliderWidth := 300
	sliderHeight := 20

	// Get current mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Check for mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Check if click is on read rate slider
		if mouseX >= sliderX && mouseX <= sliderX+sliderWidth &&
			mouseY >= sliderY-5 && mouseY <= sliderY+sliderHeight+5 {
			g.renderer.sliderDragging = "read"
			g.renderer.sliderReadX = mouseX - sliderX
			g.renderer.readRate = int64(g.renderer.sliderReadX * 1000 / sliderWidth)
			if g.renderer.rateCallback != nil {
				g.renderer.rateCallback(g.renderer.readRate, g.renderer.broadcastRate)
			}
		} else if mouseX >= sliderX && mouseX <= sliderX+sliderWidth &&
			mouseY >= sliderY+40 && mouseY <= sliderY+45+sliderHeight+5 {
			// Check if click is on broadcast rate slider
			g.renderer.sliderDragging = "broadcast"
			g.renderer.sliderBroadcastX = mouseX - sliderX
			g.renderer.broadcastRate = int64(g.renderer.sliderBroadcastX * 1000 / sliderWidth)
			if g.renderer.rateCallback != nil {
				g.renderer.rateCallback(g.renderer.readRate, g.renderer.broadcastRate)
			}
		} else {
			// Regular mouse click
			g.renderer.mousePressed = true
			g.renderer.mouseX, g.renderer.mouseY = mouseX, mouseY
			g.renderer.charBuffer = append(g.renderer.charBuffer, KEY_MOUSE)
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Handle slider dragging
		if g.renderer.sliderDragging == "read" {
			g.renderer.sliderReadX = mouseX - sliderX
			if g.renderer.sliderReadX < 0 {
				g.renderer.sliderReadX = 0
			}
			if g.renderer.sliderReadX > sliderWidth {
				g.renderer.sliderReadX = sliderWidth
			}
			g.renderer.readRate = int64(g.renderer.sliderReadX * 1000 / sliderWidth)
			if g.renderer.rateCallback != nil {
				g.renderer.rateCallback(g.renderer.readRate, g.renderer.broadcastRate)
			}
		} else if g.renderer.sliderDragging == "broadcast" {
			g.renderer.sliderBroadcastX = mouseX - sliderX
			if g.renderer.sliderBroadcastX < 0 {
				g.renderer.sliderBroadcastX = 0
			}
			if g.renderer.sliderBroadcastX > sliderWidth {
				g.renderer.sliderBroadcastX = sliderWidth
			}
			g.renderer.broadcastRate = int64(g.renderer.sliderBroadcastX * 1000 / sliderWidth)
			if g.renderer.rateCallback != nil {
				g.renderer.rateCallback(g.renderer.readRate, g.renderer.broadcastRate)
			}
		} else {
			// Regular mouse drag
			g.renderer.mousePressed = true
			g.renderer.mouseX, g.renderer.mouseY = mouseX, mouseY
			g.renderer.charBuffer = append(g.renderer.charBuffer, KEY_MOUSE)
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.renderer.sliderDragging = ""
		g.renderer.mousePressed = false
		g.renderer.charBuffer = append(g.renderer.charBuffer, KEY_MOUSE_RELEASE)
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
	grayLineColor := color.RGBA{128, 128, 128, 255} // Gray color for broadcasts to dead cells
	sliderColor := color.RGBA{200, 200, 200, 255} // Color for sliders
	sliderHandleColor := color.RGBA{100, 100, 255, 255} // Color for slider handles

	// Draw the buffer
	for y, row := range g.renderer.buffer {
		for x, cell := range row {
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

			// Draw fading cells
			fadeOpacity := g.renderer.fadingCells[y][x]
			if fadeOpacity > 0 && cell != "0" {
				// Create a color with the appropriate opacity
				alpha := uint8(fadeOpacity * 255)
				fadeColor := color.RGBA{240, 240, 240, alpha} // Off-white with fading alpha

				// Draw the fading cell
				ebitenutil.DrawRect(screen, cellX + cellPadding, cellY + cellPadding, 
					cellSize - (cellPadding * 2), cellSize - (cellPadding * 2), fadeColor)
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

							// Calculate start and end points for the line
							startX := cellX + cellSize/2
							startY := cellY + cellSize/2
							endX := float64(adjX * g.renderer.cellSize) + cellSize/2
							endY := float64(adjY * g.renderer.cellSize) + cellSize/2

							// Check if the adjacent cell is alive
							if g.renderer.buffer[adjY][adjX] == "0" {
								// Draw multiple lines to create a thicker appearance
								// Main line
								ebitenutil.DrawLine(screen, startX, startY, endX, endY, blueLineColor)
								// Additional lines with slight offsets to create thickness
								ebitenutil.DrawLine(screen, startX+1, startY, endX+1, endY, blueLineColor)
								ebitenutil.DrawLine(screen, startX-1, startY, endX-1, endY, blueLineColor)
								ebitenutil.DrawLine(screen, startX, startY+1, endX, endY+1, blueLineColor)
								ebitenutil.DrawLine(screen, startX, startY-1, endX, endY-1, blueLineColor)
							} else if g.renderer.buffer[adjY][adjX] == "-" {
								// If the adjacent cell is dead, mark it as receiving a broadcast
								g.renderer.deadCellBroadcasts[adjY][adjX] = true

								// Draw gray lines to dead cells that are receiving broadcasts
								// Main line
								ebitenutil.DrawLine(screen, startX, startY, endX, endY, grayLineColor)
								// Additional lines with slight offsets to create thickness
								ebitenutil.DrawLine(screen, startX+1, startY, endX+1, endY, grayLineColor)
								ebitenutil.DrawLine(screen, startX-1, startY, endX-1, endY, grayLineColor)
								ebitenutil.DrawLine(screen, startX, startY+1, endX, endY+1, grayLineColor)
								ebitenutil.DrawLine(screen, startX, startY-1, endX, endY-1, grayLineColor)
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

	// Draw sliders for rate control
	sliderY := 30
	sliderWidth := 300
	sliderHeight := 20
	sliderX := 50

	// Draw read rate slider
	text.Draw(screen, "Read Rate: " + strconv.FormatInt(g.renderer.readRate, 10) + "ms", g.renderer.fontFace, sliderX, sliderY - 5, color.White)
	ebitenutil.DrawRect(screen, float64(sliderX), float64(sliderY), float64(sliderWidth), float64(sliderHeight), sliderColor)
	ebitenutil.DrawRect(screen, float64(sliderX + g.renderer.sliderReadX - 5), float64(sliderY - 5), 10, float64(sliderHeight + 10), sliderHandleColor)

	// Draw broadcast rate slider
	text.Draw(screen, "Broadcast Rate: " + strconv.FormatInt(g.renderer.broadcastRate, 10) + "ms", g.renderer.fontFace, sliderX, sliderY + 40, color.White)
	ebitenutil.DrawRect(screen, float64(sliderX), float64(sliderY + 45), float64(sliderWidth), float64(sliderHeight), sliderColor)
	ebitenutil.DrawRect(screen, float64(sliderX + g.renderer.sliderBroadcastX - 5), float64(sliderY + 40), 10, float64(sliderHeight + 10), sliderHandleColor)
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
	deadCellBroadcasts [][]bool // Tracks broadcasts to dead cells
	fadingCells    [][]float64 // Tracks cells that are fading out (0.0 to 1.0, where 0.0 is fully faded)
	lastFrameTime  time.Time // Used to calculate time delta for fading
	statsWindows   []*EbitenStatsWindow
	charBuffer     []Key
	charMutex      sync.Mutex
	mouseX, mouseY int
	mousePressed   bool
	game           *EbitenGame
	fontFace       font.Face

	// Rate control
	readRate       int64
	broadcastRate  int64
	rateCallback   func(readRate, broadcastRate int64)

	// Slider state
	sliderReadX    int
	sliderBroadcastX int
	sliderDragging string
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
		deadCellBroadcasts: make([][]bool, height),
		fadingCells:    make([][]float64, height),
		lastFrameTime:  time.Now(),
		statsWindows:   make([]*EbitenStatsWindow, 0),
		charBuffer:     make([]Key, 0),
		fontFace:       basicfont.Face7x13,
		readRate:       500, // Default read rate in milliseconds
		broadcastRate:  500, // Default broadcast rate in milliseconds
		sliderReadX:    150, // Initial slider position
		sliderBroadcastX: 150, // Initial slider position
	}

	// Initialize buffer, communications, deadCellBroadcasts, and fadingCells
	for i := range r.buffer {
		r.buffer[i] = make([]string, width)
		r.communications[i] = make([]bool, width)
		r.deadCellBroadcasts[i] = make([]bool, width)
		r.fadingCells[i] = make([]float64, width)
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
	// Reset the buffer, communications, deadCellBroadcasts, and fadingCells with the current dimensions
	r.buffer = make([][]string, r.height)
	r.communications = make([][]bool, r.height)
	r.deadCellBroadcasts = make([][]bool, r.height)
	r.fadingCells = make([][]float64, r.height)
	for i := range r.buffer {
		r.buffer[i] = make([]string, r.width)
		r.communications[i] = make([]bool, r.width)
		r.deadCellBroadcasts[i] = make([]bool, r.width)
		r.fadingCells[i] = make([]float64, r.width)
	}
	r.lastFrameTime = time.Now()
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
		// Check if the cell was previously alive and is now not alive
		if r.buffer[y][x] == "0" && ach != "0" {
			// Start fading out the cell
			r.fadingCells[y][x] = 1.0 // Start with full opacity
		}

		r.buffer[y][x] = ach
		// Mark cell as communicating if it's alive (represented by "0")
		if ach == "0" {
			r.communications[y][x] = true
			// Reset fading opacity when a cell becomes alive
			r.fadingCells[y][x] = 0.0
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
			r.deadCellBroadcasts[i][j] = false
			r.fadingCells[i][j] = 0.0 // Reset fading cells
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

// GetReadRate returns the current read rate
func (r *EbitenRenderer) GetReadRate() int64 {
	return r.readRate
}

// GetBroadcastRate returns the current broadcast rate
func (r *EbitenRenderer) GetBroadcastRate() int64 {
	return r.broadcastRate
}

// SetRateChangeCallback sets the callback function for rate changes
func (r *EbitenRenderer) SetRateChangeCallback(callback func(readRate, broadcastRate int64)) {
	r.rateCallback = callback
}

// SetInitialRates sets the initial read and broadcast rates
func (r *EbitenRenderer) SetInitialRates(readRate, broadcastRate int64) {
	r.readRate = readRate
	r.broadcastRate = broadcastRate

	// Update slider positions based on the rates
	sliderWidth := 300 // Must match the width used in the Draw method
	r.sliderReadX = int(readRate * int64(sliderWidth) / 1000)
	r.sliderBroadcastX = int(broadcastRate * int64(sliderWidth) / 1000)
}
