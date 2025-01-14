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
		r.Clear()

		cWorld := internal.NewChannelWorld[internal.ChannelCell](r, 0.1)
		// If you want the whole empty world drawn so the grid is filled
		//cWorld.Refresh()
		cWorld.Bootstrap()
		glog.GetLogger().Info("Setup")
		for {
			select {
			case <-ctx.Done():
				r.End()
				return
			default:
				//Makes render choppy this one routine can run at cpu time
				//time.Sleep(100 * time.Millisecond)
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
