package decider

type Command interface {
	command()
}

type Event interface {
	event()
}

type Decider[State any] struct {
	Decide          func(state State, command Command) ([]Event, error)
	Evolve          func(state State, event Event) State
	GetInitialState func() State
}
