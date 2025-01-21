package internal

import (
	"fmt"
	glog "gogol2/log"
	"gogol2/renderer"
	"time"
)

const (
	Heartbeat   = "heartbeat"
	Broadcast   = "broadcast"
	Died        = "died"
	Resurrected = "ressurected"
)

type CellEvent struct {
	name  string
	count int
}

type Stats struct {
	r          renderer.Renderer
	x, y       int
	heartbeats int64
	broadcasts int64
	died       int64
	eventChan  chan CellEvent
}

func NewStats(r renderer.Renderer, location string) *Stats {
	y, _ := r.Dimensions()
	s := &Stats{
		r:          r,
		y:          y - 1,
		x:          1,
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
		for {
			time.Sleep(1 * time.Second)
			s.Update()
		drainLoop:
			for {
				consumed := 0
				select {
				case e := <-s.eventChan:
					glog.GetLogger().Debug("Consuming", "event", e.name)
					if e.name == Heartbeat {
						s.heartbeats++
					} else if e.name == Broadcast {
						s.broadcasts++
					} else if e.name == Died {
						s.died++
					} else if e.name == Resurrected {
						s.died--
					}
					consumed++
					if consumed > 10 {
						consumed = 0
						s.Update()
					}
				default:
					break drainLoop
				}
			}
		}
	}()
}

func (s *Stats) String() string {
	return fmt.Sprintf(
		"Died: %v "+
			"Broadcasts: %v "+
			"Heartbeats: %v",
		s.died,
		s.broadcasts,
		s.heartbeats)
}

func (s *Stats) Update() {
	s.r.DrawAt(s.y-1, s.x+2, s.String())
	glog.GetLogger().Debug("stats update", "Data", s.String())
	s.r.BufferUpdate()
}
