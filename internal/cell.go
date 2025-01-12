package internal

import "gogol2/game"

type ChannelCell struct {
	game.Life
	state    bool
	channels []<-chan bool
}

func (c *ChannelCell) State() bool {
	return c.state
}

func (c *ChannelCell) SetState(state bool) {
	c.state = state
}

func (c *ChannelCell) AddChannel(ch <-chan bool) {
	c.channels = append(c.channels, ch)
}
