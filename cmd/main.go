package main

import (
	"context"
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
	r := renderer.NewShellRenderer()
	defer r.End()

	go func() {
		r.Beep()
		r.Draw("Hello, world.go!")
		r.Refresh()
		//r.Display.GetChar()

		//y, x := r.Dimensions()

		cWorld := internal.NewChannelWorld[internal.ChannelCell](r, 0.1)
		cWorld.Bootstrap()

		glog.GetLogger().Info("Setup")
		for {
			select {
			case <-ctx.Done():
				r.End()
				return
			default:
				time.Sleep(100 * time.Millisecond)
				cWorld.Refresh()
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
