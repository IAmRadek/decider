package decider

type Events[State any, V any] struct {
	Events  []Event[State]
	Version V
}

type EventStore[State any, ID any, V any] interface {
	ReadStream(streamID ID) (Events[State, V], error)
	AppendToStream(streamID ID, version V, events []Event[State]) error
}

type CommandHandler[State any, ID any, EV any] struct {
	EventStore      EventStore[State, ID, EV]
	GetInitialState func() State
}

func (c *CommandHandler[S, ID, EV]) HandleCommand(id ID, command Command[S]) ([]Event[S], error) {
	es := c.EventStore

	events, readErr := es.ReadStream(id)
	if readErr != nil {
		return nil, readErr
	}

	state := c.GetInitialState()
	for _, event := range events.Events {
		state = event.Evolve(state)
	}

	newEvents, decideErr := command.Decide(state)
	if decideErr != nil {
		return nil, decideErr
	}

	appendErr := es.AppendToStream(id, events.Version, newEvents)
	if appendErr != nil {
		return nil, appendErr
	}

	return newEvents, nil
}
