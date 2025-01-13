package main

import (
	"context"
	"gogol2/internal"
	"gogol2/renderer/mock"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	//slog.SetLogLoggerLevel(slog.LevelDebug)
	go func() {
		r := mock.NewMockRenderer()
		defer r.End()
		r.Beep()
		r.Draw("Hello, world.go!")
		r.Refresh()
		//r.Display.GetChar()

		//y, x := r.Dimensions()

		cWorld := internal.NewChannelWorld[internal.ChannelCell](r)
		//gWorld.ref
		//var a game.Life = cWorld.Cells()[0][0]
		cWorld.Bootstrap()

		//gWorld := game.NewChannelWorld(x, y)
		//
		//world.Bootstrap()
		//
		//worldRenderer := internal.NewChannelWorld(world, *r)
		//
		slog.Info("Setup")
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
