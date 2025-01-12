package game

type Life interface {
	State() bool
	SetState(bool)
}
