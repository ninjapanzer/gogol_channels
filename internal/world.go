package internal

import (
	"fmt"
	glog "gogol2/log"
	"gogol2/renderer"
	"math/rand"
	"time"
)

type ChannelWorld[T ChannelCell] struct {
	r        renderer.Renderer
	s        *Stats
	cells    [][]*ChannelCell
	initProb float64
}

func NewChannelWorld[T ChannelCell](r renderer.Renderer, prob float64) *ChannelWorld[T] {
	y, x := r.Dimensions()
	cells := make([][]*ChannelCell, y)
	for i := range cells {
		cells[i] = make([]*ChannelCell, x)
		for j := range cells[i] {
			cells[i][j] = NewChannelCell(false, fmt.Sprintf("%s-%s", i, j))
		}
	}
	return &ChannelWorld[T]{
		r:        r,
		s:        NewStats(r, "bottom"),
		cells:    cells,
		initProb: prob,
	}
}

func (w *ChannelWorld[T]) Refresh() {
	w.r.Clear()
	w.DrawWorld()
	w.Refresh()
}

func (w *ChannelWorld[T]) ComputeState() {
}

func (w *ChannelWorld[T]) Cells() [][]*ChannelCell {
	return w.cells
}

func (w *ChannelWorld[T]) Bootstrap() {
	w.initializeProbabilisticDistributionOfLife(w.initProb)
	w.setupNeighborhood()
}

func (w *ChannelWorld[T]) initializeProbabilisticDistributionOfLife(prob float64) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i, _ := range w.cells {
		for j, _ := range w.cells[i] {
			target := w.cells[i][j]
			target.SetRenderer(w.DrawCell(i, j))
			target.SetStatsFunc(w.s.AddEvent)
			if rng.Float64() < prob {
				target.SilentSetState(true)
			}
			w.DrawCell(i, j)
		}
	}
	w.r.Refresh()
}

func (w *ChannelWorld[T]) setupNeighborhood() {
	height := len(w.cells)
	width := len(w.cells[0])

	for i, _ := range w.cells {
		for j, _ := range w.cells[i] {
			linkNeighbors(w.cells, w.cells[i][j], i, j, width, height)
		}
	}
}

func linkNeighbors(cells [][]*ChannelCell, cell *ChannelCell, y, x, width, height int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == y && j == x {
				continue
			}
			if y+i < 0 {
				continue
			}
			if y+i >= height {
				continue
			}
			if x+j < 0 {
				continue
			}
			if x+j >= width {
				continue
			}

			glog.GetLogger().Info("Adding Neighbor", "CX", x, "CY", y, "TX", x+i, "TY", y+j)
			cell.AddChannel(cells[y+i][x+j].BroadcastChan())
			if cells[y+i][x+j].State() {
				cell.AddNeighborState(1)
			} else {
				cell.AddNeighborState(0)
			}
		}
	}
	go cell.heartbeat()
	go cell.Live()
}

func (w *ChannelWorld[T]) DrawCell(y, x int) func(bool) {
	return func(state bool) {
		if state {
			w.r.DrawAt(y, x, "0")
		} else {
			w.r.DrawAt(y, x, "-")
		}
		w.r.BufferUpdate()
	}
}

func (w *ChannelWorld[T]) DrawWorld() {
	for y, row := range w.Cells() {
		for x, cell := range row {
			if cell.State() {
				w.r.DrawAt(y, x, "0")
			} else {
				w.r.DrawAt(y, x, "-")
			}
		}
	}
}
