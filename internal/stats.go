package internal

import (
	"fmt"
	glog "github.com/ninjapanzer/gogol_channels/log"
	"github.com/ninjapanzer/gogol_channels/renderer"
	"strings"
	"time"
)

const (
	Heartbeat   = "heartbeat"
	Broadcast   = "broadcast"
	Died        = "died"
	Resurrected = "resurrected"
)

type CellEvent struct {
	name  string
	count int
}

type Stats struct {
	r                  renderer.Renderer
	st                 renderer.StatsWindow
	x, y               int
	heartbeats         int64
	heartbeatPerSecond int64
	broadcasts         int64
	broadcastPerSecond int64
	died               int64
	diedPerSecond      int64
	eventChan          chan CellEvent
}

func NewStats(r renderer.Renderer, location string) *Stats {
	_, x := r.Dimensions()
	// Position the stats window in the top right with padding
	padding := 20 // Padding from the right edge
	// Use a smaller width for the stats window to make it fit the text better
	statsWidth := 80 // Approximate width based on the text length
	st := r.CreateStatsWindow(3, statsWidth, 1, x-statsWidth-padding)

	s := &Stats{
		r:          r,
		st:         st,
		y:          1,
		x:          x-statsWidth-padding,
		heartbeats: 0,
		broadcasts: 0,
		died:       0,
		eventChan:  make(chan CellEvent, 10000),
	}

	s.collectStats()

	return s
}

func (s *Stats) AddEvent(event CellEvent) {
	glog.GetLogger().Debug("AddEvent", "Event", event.name)
	s.eventChan <- event
}

func (s *Stats) collectStats() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		hps := 0
		bps := 0
		dps := 0
		for {
			select {
			case <-ticker.C:
				s.heartbeatPerSecond = int64(hps)
				s.broadcastPerSecond = int64(bps)
				s.diedPerSecond = int64(dps)
				hps = 0
				bps = 0
				dps = 0
			case e := <-s.eventChan:
				if e.name == Heartbeat {
					hps += e.count
					s.heartbeats += int64(e.count)
				} else if e.name == Broadcast {
					bps += e.count
					s.broadcasts += int64(e.count)
				} else if e.name == Died {
					dps += e.count
					s.died += int64(e.count)
				} else if e.name == Resurrected {
					dps -= e.count
					s.died -= int64(e.count)
				}
			default:
			}
		}
	}()
	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			s.Update()
		}
	}()
}

func (s *Stats) String() string {
	return fmt.Sprintf(
		"Died: %v "+
			"d/s: %v "+
			"Broadcasts: %v "+
			"b/s %v "+
			"Heartbeats: %v "+
			"h/s %v",
		s.died,
		s.diedPerSecond,
		s.broadcasts,
		s.broadcastPerSecond,
		s.heartbeats,
		s.heartbeatPerSecond)
}

func (s *Stats) Update() {
	m := s.String()
	// Clear the entire line
	s.st.MovePrint(0, 0, strings.Repeat(" ", len(m)+20))
	// Position the text at the right edge of the window
	s.st.MovePrint(1, 0, m)
	glog.GetLogger().Debug("stats update", "Data", s.String())
	s.st.NoutRefresh()
}
