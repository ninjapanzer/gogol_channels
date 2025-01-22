package main

import (
	"context"
	"github.com/gbin/goncurses"
	"gogol2/internal"
	glog "gogol2/log"
	"gogol2/renderer"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	closer := glog.InitLogger()
	defer closer()
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	r := renderer.NewShellRenderer(0)
	defer r.End()

	cWorld := internal.NewChannelWorld[internal.ChannelCell](r, 0.13)
	cWorld.Bootstrap()
	goncurses.Update()

	go func() {
		for {
			ch := r.GetChar() // Wait for input

			// Check for mouse event
			glog.GetLogger().Debug("event", "int", ch, "mousekey", goncurses.KEY_MOUSE)
			if ch == goncurses.KEY_MOUSE {
				// Capture mouse event
				mevent := goncurses.GetMouse()
				mx := int(mevent.X)
				my := int(mevent.Y)

				//cWorld.DrawCell(my, mx)
				glog.GetLogger().Debug("mouse event", "y", my, "x", mx)
				target := cWorld.Cells()[my][mx]
				target.SilentSetState(!target.State())
			} else if ch == 'q' { // Quit on 'q' press
				break
			}
		}
	}()

	go func() {
		r.Beep()
		r.Draw("Hello, world.go!")
		r.Refresh()
		r.Clear()
		// If you want the whole empty world drawn so the grid is filled
		//cWorld.Refresh()
		glog.GetLogger().Info("Rendering")
		for {
			select {
			case <-ctx.Done():
				r.End()
				return
			default:
				time.Sleep(30 * time.Millisecond)
				//Makes render choppy this one routine can run at cpu time
				goncurses.Update()
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
