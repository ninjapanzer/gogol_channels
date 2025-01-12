package main

import (
	"context"
	"gogol2/renderer"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		r := renderer.NewShellRenderer()
		defer r.End()
		r.Beep()
		r.Draw("Hello, world.go!")
		r.Refresh()
		//r.Display.GetChar()

		//y, x := r.Dimensions()

		//gWorld := internal.NewWorld(r)
		//gWorld.ref
		//gWorld.Cells()

		//gWorld := game.NewWorld(x, y)
		//
		//world.Bootstrap()
		//
		//worldRenderer := internal.NewWorld(world, *r)
		//
		for {
			select {
			case <-ctx.Done():
				r.End()
				return
			default:
				time.Sleep(100 * time.Millisecond)
				//world.ComputeState()
				//r.Clear()
				//worldRenderer.Refresh()
				r.Refresh()
			}
		}
	}()

	select {
	case <-cancelChan:
		cancel()
		println("Cancelled")
	case <-ctx.Done():
		println("Done")
	}
}
