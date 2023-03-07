package course

import (
	"github.com/IAmRadek/decider/v1"
)

type Event = decider.Event

type StudentAddedToCourse struct {
	decider.Event `json:"-"`

	Name string
}

type StudentRemovedFromCourse struct {
	decider.Event `json:"-"`

	Name string
}

func evolve(state Course, event decider.Event) Course {
	switch v := event.(type) {
	case StudentAddedToCourse:
		return Course{
			Students: append(state.Students, v.Name),
		}
	case StudentRemovedFromCourse:
		for i, name := range state.Students {
			if name == v.Name {
				state.Students = append(state.Students[:i], state.Students[i+1:]...)
				break
			}
		}
		return state
	}

	panic("unreachable")
}
