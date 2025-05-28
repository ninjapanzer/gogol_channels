package main

import (
	"context"
	"flag"
	"github.com/gbin/goncurses"
	"github.com/ninjapanzer/gogol_channels/internal"
	glog "github.com/ninjapanzer/gogol_channels/log"
	"github.com/ninjapanzer/gogol_channels/renderer"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Parse command line arguments
	rendererType := flag.String("renderer", "ncurses", "Renderer to use (ncurses or ebiten)")
	readRate := flag.Int64("read-rate", 500, "Initial read rate in milliseconds")
	broadcastRate := flag.Int64("broadcast-rate", 500, "Initial broadcast rate in milliseconds")
	flag.Parse()

	closer := glog.InitLogger()
	defer closer()
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// Set the initial global rate variables from command line arguments
	internal.GlobalReadRate = *readRate
	internal.GlobalBroadcastRate = *broadcastRate

	// Create the appropriate renderer based on the command line argument
	var r renderer.Renderer
	switch *rendererType {
	case "ebiten":
		r = renderer.NewEbitenRenderer(0)
		// Initialize the renderer with the command line values
		r.SetInitialRates(*readRate, *broadcastRate)
		// Set the rate change callback to update the global rate variables
		r.SetRateChangeCallback(func(readRate, broadcastRate int64) {
			internal.GlobalReadRate = readRate
			internal.GlobalBroadcastRate = broadcastRate
		})
	default:
		r = renderer.NewShellRenderer(0)
		// Initialize the renderer with the command line values
		r.SetInitialRates(*readRate, *broadcastRate)
	}
	defer r.End()

	cWorld := internal.NewChannelWorld[internal.ChannelCell](r, 0.13)
	cWorld.Bootstrap()

	// Only call goncurses.Update() if using the ncurses renderer
	if *rendererType == "ncurses" {
		goncurses.Update()
	}

	// Track last mouse position for drag detection
	var lastMouseX, lastMouseY int
	var isMouseDragging bool

	go func() {
		for {
			ch := r.GetChar() // Wait for input

			// Check for mouse event
			glog.GetLogger().Debug("event", "int", ch, "mousekey", renderer.KEY_MOUSE)
			if ch == renderer.KEY_MOUSE {
				// Capture mouse event
				mevent := r.GetMouse()
				mx := mevent.X
				my := mevent.Y

				//cWorld.DrawCell(my, mx)
				glog.GetLogger().Debug("mouse event", "y", my, "x", mx)

				// Check if this is a new click or a drag
				if !isMouseDragging {
					isMouseDragging = true
					lastMouseX, lastMouseY = mx, my
				}

				// Always generate cells on first click, and on drag when position changes
				if !isMouseDragging || mx != lastMouseX || my != lastMouseY {
					// Create a randomized cluster of cells in a circle of radius 5
					radius := 2 // Radius of 2 will create a circle of about 5 cells
					for dy := -radius; dy <= radius; dy++ {
						for dx := -radius; dx <= radius; dx++ {
							// Check if the point is within the circle
							if dx*dx + dy*dy <= radius*radius {
								// Add randomness - only set some cells to alive
								if rand.Float64() < 0.7 { // 70% chance of becoming alive
									ny, nx := my+dy, mx+dx
									if ny >= 0 && ny < len(cWorld.Cells()) && nx >= 0 && nx < len(cWorld.Cells()[ny]) {
										target := cWorld.Cells()[ny][nx]
										target.SilentSetState(true) // Set cells to alive
									}
								}
							}
						}
					}

					// Update last position
					lastMouseX, lastMouseY = mx, my
				}
			} else if ch == renderer.KEY_MOUSE_RELEASE { // Mouse button released
				// Reset drag state when mouse is released
				isMouseDragging = false
				glog.GetLogger().Debug("mouse released")
			} else if ch == 'q' { // Quit on 'q' press
				cancel()
				return
			} else {
				// Other key pressed
				isMouseDragging = false
			}
		}
	}()

	go func() {
		glog.GetLogger().Info("Rendering")
		for {
			select {
			case <-ctx.Done():
				r.End()
				return
			default:
				//time.Sleep(30 * time.Millisecond)
				//Makes render choppy this one routine can run at cpu time
				//cWorld.DrawWorld()
				r.Refresh()
			}
		}
	}()

	select {
	case <-cancelChan:
		cancel()
		r.End()
		println("Cancelled")
	case <-ctx.Done():
		r.End()
		println("Done")
	}
}
