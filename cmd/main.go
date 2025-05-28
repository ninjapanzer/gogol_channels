package main

import (
	"context"
	"flag"
	"github.com/gbin/goncurses"
	"github.com/ninjapanzer/gogol_channels/internal"
	glog "github.com/ninjapanzer/gogol_channels/log"
	"github.com/ninjapanzer/gogol_channels/renderer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command line arguments
	rendererType := flag.String("renderer", "ncurses", "Renderer to use (ncurses or ebiten)")
	flag.Parse()

	closer := glog.InitLogger()
	defer closer()
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// Create the appropriate renderer based on the command line argument
	var r renderer.Renderer
	switch *rendererType {
	case "ebiten":
		r = renderer.NewEbitenRenderer(0)
	default:
		r = renderer.NewShellRenderer(0)
	}
	defer r.End()

	cWorld := internal.NewChannelWorld[internal.ChannelCell](r, 0.13)
	cWorld.Bootstrap()

	// Only call goncurses.Update() if using the ncurses renderer
	if *rendererType == "ncurses" {
		goncurses.Update()
	}

	//go func() {
	//	for {
	//		ch := r.GetChar() // Wait for input
	//
	//		// Check for mouse event
	//		glog.GetLogger().Debug("event", "int", ch, "mousekey", goncurses.KEY_MOUSE)
	//		if ch == goncurses.KEY_MOUSE {
	//			// Capture mouse event
	//			mevent := goncurses.GetMouse()
	//			mx := int(mevent.X)
	//			my := int(mevent.Y)
	//
	//			//cWorld.DrawCell(my, mx)
	//			glog.GetLogger().Debug("mouse event", "y", my, "x", mx)
	//			target := cWorld.Cells()[my][mx]
	//			target.SilentSetState(!target.State())
	//		} else if ch == 'q' { // Quit on 'q' press
	//			break
	//		}
	//	}
	//}()

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
