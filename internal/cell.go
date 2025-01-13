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
}

func NewChannelCell(state bool) *ChannelCell {
	b := &ChannelCell{
		state:          state,
		neighborChans:  make([]<-chan bool, 0),
		neighborStates: 0,
		broadcast:      make(chan bool, 1),
	}

	return b
}

func (c *ChannelCell) State() bool {
	return c.state
}

func (c *ChannelCell) SetState(state bool) {
	c.state = state
	c.broadcast <- state
}

func (c *ChannelCell) SilentSetState(state bool) {
	c.state = state
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

func (c *ChannelCell) listenAndUpdate() {
	// Map to store the latest states from the channels
	var latestStates uint = 0

	for {
		// Check each neighbor's channel for new state updates
		for _, neighborChan := range c.neighborChans {
			select {
			case _ = <-neighborChan: // If a message is received on this channel
				// Update the state for this particular neighbor
				latestStates <<= 1
				latestStates |= 1
				//latestStates <<= 1
				//if state {
				//	latestStates |= 1 // Neighbor is alive
				//} else {
				//	latestStates |= 0 // Neighbor is dead
				//}
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
			c.SetState(newState)
			glog.GetLogger().Info("Cell Updated", "New State", c.state, "Reason", reason, "Old State", oldState)
		}

		// Sleep to simulate waiting for the next update
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
			return true, "Porrige Just Right"
		}
	} else {
		if aliveCount == 3 {
			return true, "Nobody Expects the Cellular Resurrection"
		}
	}

	return false, "Still Mostly Dead"
}
