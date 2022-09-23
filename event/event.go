package event

type Handler func()

type GourdEvent struct {
	Boot  Handler
	Init  Handler
	Start Handler
}
