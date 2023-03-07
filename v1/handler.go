package decider

type Events[State any, V any] struct {
	Events  []Event
	Version V
}

type EventStore[State any, ID any, V any] interface {
	ReadStream(streamID ID) (Events[State, V], error)
	AppendToStream(streamID ID, version V, events []Event) error
}

type CommandHandler[S any, ID any, EV any] struct {
	EventStore EventStore[S, ID, EV]
	Decider    Decider[S]
}

func (c *CommandHandler[S, ID, EV]) HandleCommand(id ID, command Command) ([]Event, error) {
	es := c.EventStore

	events, readErr := es.ReadStream(id)
	if readErr != nil {
		return nil, readErr
	}

	state := c.Decider.GetInitialState()
	for _, event := range events.Events {
		state = c.Decider.Evolve(state, event)
	}

	newEvents, decideErr := c.Decider.Decide(state, command)
	if decideErr != nil {
		return nil, decideErr
	}

	appendErr := es.AppendToStream(id, events.Version, newEvents)
	if appendErr != nil {
		return nil, appendErr
	}

	return newEvents, nil
}
