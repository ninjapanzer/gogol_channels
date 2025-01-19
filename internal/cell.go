package internal

import (
	"gogol2/game"
	glog "gogol2/log"
	"math/bits"
	"time"
)

type ChannelCell struct {
	game.Life
	state          bool
	neighborChans  []<-chan bool
	neighborStates uint
	broadcast      chan bool
	renderFunc     func(bool)
}

func NewChannelCell(state bool) *ChannelCell {
	b := &ChannelCell{
		state:          state,
		neighborChans:  make([]<-chan bool, 0),
		neighborStates: 0,
		broadcast:      make(chan bool, 10),
	}

	return b
}

func (c *ChannelCell) State() bool {
	return c.state
}

func (c *ChannelCell) SetState(state bool) {
	c.state = state
	glog.GetLogger().Info("Broadcasting", "NewState", state)
	c.broadcast <- state
	glog.GetLogger().Info("Broadcasted", "NewState", state)
	c.renderFunc(c.state)
	//time.Sleep(1000 * time.Millisecond)
}

func (c *ChannelCell) SilentSetState(state bool) {
	c.state = state
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

func (c *ChannelCell) listenAndUpdate() {
	// Map to store the latest states from the channels

	for {
		var latestStates uint = 0
		// Check each neighbor's channel for new state updates
		for _, neighborChan := range c.neighborChans {
			select {
			case _ = <-neighborChan: // If a message is received on this channel
				latestStates <<= 1
				latestStates |= 1
			default:
				latestStates <<= 1
				latestStates |= 0
			}
		}
		glog.GetLogger().Info("Consumed", "Latest", latestStates)

		// After checking all channels, update state based on the latest values
		// For example, we can just print the new state based on the latest updates:
		oldState := c.state
		newState, reason := c.computeStateFromNeighbors(latestStates)
		if oldState != newState {
			c.SetState(newState)
			glog.GetLogger().Info("Cell Updated", "New State", c.state, "Reason", reason, "Old State", oldState)
		}

		time.Sleep(100 * time.Millisecond)
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
			return true, "Nobody Expects the Cellular Resurrection"
		}
	}

	return false, "Still Mostly Dead"
}
