package course

import (
	"github.com/IAmRadek/decider/v2"
)

type Event = decider.Event[Course]

type StudentAddedToCourse struct {
	Name string
}

func (event StudentAddedToCourse) Evolve(state Course) Course {
	state.Students = append(state.Students, event.Name)
	return state
}

type StudentRemovedFromCourse struct {
	Name string
}

func (event StudentRemovedFromCourse) Evolve(state Course) Course {
	for i, name := range state.Students {
		if name == event.Name {
			state.Students = append(state.Students[:i], state.Students[i+1:]...)
			break
		}
	}
	return state
}
