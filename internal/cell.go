package internal

import "gogol2/game"

type ChannelCell struct {
	game.Life
	state     bool
	neighbors []<-chan bool
	broadcast chan bool
}

func NewChannelCell(state bool) *ChannelCell {
	b := &ChannelCell{
		state:     state,
		neighbors: make([]<-chan bool, 0),
		broadcast: make(chan bool),
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
	c.neighbors = append(c.neighbors, ch)
}

func (c *ChannelCell) BroadcastChan() <-chan bool {
	return c.broadcast
}
