package world

type World interface {
	Name() string
	Equals(other World) bool
}
