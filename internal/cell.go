package internal

import (
	"github.com/ninjapanzer/gogol_channels/game"
	glog "github.com/ninjapanzer/gogol_channels/log"
	"math/bits"
	"sync/atomic"
	"time"
)

// Global variables for cell read and broadcast rates
var (
	GlobalReadRate     int64 = 500 // Default read rate in milliseconds
	GlobalBroadcastRate int64 = 500 // Default broadcast rate in milliseconds
)

type ChannelCell struct {
	game.Life
	state          bool
	location       string
	readSpeed      time.Duration
	broadcastSpeed time.Duration
	ticker         *time.Ticker
	neighborChans  []<-chan bool
	neighborStates uint
	broadcast      chan bool
	renderFunc     func(bool)
	statsFunc      func(event CellEvent)
}

func NewChannelCell(state bool, location string) *ChannelCell {
	readRate := time.Duration(atomic.LoadInt64(&GlobalReadRate))
	broadcastRate := time.Duration(atomic.LoadInt64(&GlobalBroadcastRate))

	b := &ChannelCell{
		state:          state,
		location:       location,
		readSpeed:      readRate,
		broadcastSpeed: broadcastRate,
		ticker:         time.NewTicker(broadcastRate * time.Millisecond),
		neighborChans:  make([]<-chan bool, 0),
		neighborStates: 0,
		broadcast:      make(chan bool, 1),
		statsFunc:      func(event CellEvent) {},
	}

	return b
}

func (c *ChannelCell) State() bool {
	return c.state
}

func (c *ChannelCell) SetState(state bool) {
	c.state = state
	c.renderFunc(c.state)
	c.statsBroadcast()
	c.broadcast <- state
}

func (c *ChannelCell) SilentSetState(state bool) {
	c.state = state
	glog.GetLogger().Debug("Silent Set State:", "name", c.location, "state", c.state)
	c.renderFunc(c.state)
}

func (c *ChannelCell) AddChannel(ch <-chan bool) {
	c.neighborChans = append(c.neighborChans, ch)
}

func (c *ChannelCell) AddNeighborState(s uint) {
	c.neighborStates <<= 1
	c.neighborStates |= s
}

func (c *ChannelCell) BroadcastChan() <-chan bool {
	return c.broadcast
}

func (c *ChannelCell) Live() {
	c.listenAndUpdate()
}

func (c *ChannelCell) SetRenderer(r func(bool)) {
	c.renderFunc = r
}

func (c *ChannelCell) SetStatsFunc(s func(event CellEvent)) {
	c.statsFunc = s
}

func (c *ChannelCell) heartbeat() {
	for {
		// Use the current global broadcast rate
		broadcastRate := time.Duration(atomic.LoadInt64(&GlobalBroadcastRate))
		c.broadcastSpeed = broadcastRate

		time.Sleep(c.broadcastSpeed * time.Millisecond)
		if c.state {
			c.statsHeartbeat()
			glog.GetLogger().Debug("Heartbeat", "name", c.location)
			c.statsBroadcast()
			c.broadcast <- true
		}
	}
}

func (c *ChannelCell) listenAndUpdate() {
	// Map to store the latest states from the channels
	var latestStates uint = 0

	for {
		// Use the current global read rate
		readRate := time.Duration(atomic.LoadInt64(&GlobalReadRate))
		c.readSpeed = readRate

		time.Sleep(c.readSpeed * time.Millisecond)
		for _, neighborChan := range c.neighborChans {
			select {
			case _ = <-neighborChan: // If a message is received on this channel
				latestStates <<= 1
				latestStates |= 1
				//gotUpdate = true
			default:
				latestStates <<= 1
				latestStates |= 0
			}
		}
		glog.GetLogger().Debug("Consumed", "Latest", latestStates)

		// After checking all channels, update state based on the latest values
		// For example, we can just print the new state based on the latest updates:
		oldState := c.state
		newState, reason := c.computeStateFromNeighbors(latestStates)
		if oldState != newState {
			if newState {
				c.SilentSetState(newState)
			} else {
				c.statsDied()
				c.SetState(newState)
			}
			glog.GetLogger().Debug("Cell Updated", "name", c.location, "New State", c.state, "Reason", reason, "Old State", oldState)
		}
	}
}

func (c *ChannelCell) computeStateFromNeighbors(latestStates uint) (bool, string) {

	//Basic Rules:
	//
	//Any live cell with fewer than two live neighbors dies (underpopulation).
	//Any live cell with two or three live neighbors continues to live (survival).
	//Any live cell with more than three live neighbors dies (overpopulation).
	//Any dead cell with exactly three live neighbors becomes a live cell (reproduction).
	//
	// Simple rule: a cell becomes alive if any of its neighbors are alive (just an example)
	aliveCount := 0
	c.neighborStates = c.neighborStates ^ latestStates

	aliveCount = bits.OnesCount(c.neighborStates)

	if c.state {
		if aliveCount < 2 || aliveCount > 3 {
			return false, "Under or Over Population"
		} else {
			return true, "Porridge Just Right"
		}
	} else {
		if aliveCount == 3 {
			c.statsResurrected()
			return true, "Nobody Expects the Cellular Resurrection"
		}
	}

	return false, "Still Mostly Dead"
}

func (c *ChannelCell) statsBroadcast() {
	c.statsFunc(CellEvent{
		name:  Broadcast,
		count: 1,
	})
}

func (c *ChannelCell) statsHeartbeat() {
	c.statsFunc(CellEvent{
		name:  Heartbeat,
		count: 1,
	})
}

func (c *ChannelCell) statsDied() {
	c.statsFunc(CellEvent{
		name:  Died,
		count: 1,
	})
}

func (c *ChannelCell) statsResurrected() {
	c.statsFunc(CellEvent{
		name:  Resurrected,
		count: 1,
	})
}
