package decider

type Command[State any] interface {
	Decide(state State) ([]Event[State], error)
}

type Event[State any] interface {
	Evolve(state State) State
}
