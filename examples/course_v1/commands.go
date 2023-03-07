package course

import (
	"fmt"

	"github.com/IAmRadek/decider/v1"
)

var (
	ErrTooManyStudents = fmt.Errorf("too many students")
)

type AddStudentToCourse struct {
	decider.Command

	Name string
}

type RemoveStudentFromCourse struct {
	decider.Command

	Name string
}

type MakeSnapshot struct {
	decider.Command
	State Course
}

func decide(state Course, command decider.Command) ([]decider.Event, error) {
	switch v := command.(type) {
	case AddStudentToCourse:
		if len(state.Students) > 10 {
			return nil, ErrTooManyStudents
		}
		if contains(state.Students, v.Name) {
			return nil, nil
		}
		return []decider.Event{StudentAddedToCourse{Name: v.Name}}, nil
	case RemoveStudentFromCourse:
		if !contains(state.Students, v.Name) {
			return nil, nil
		}
		return []decider.Event{StudentRemovedFromCourse{Name: v.Name}}, nil
	}
	panic("unreachable")
}

func contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
