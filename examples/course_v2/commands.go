package course

import (
	"fmt"
)

var (
	ErrTooManyStudents = fmt.Errorf("too many students")
)

type AddStudentToCourse struct {
	Name string
}

func (cmd AddStudentToCourse) Decide(state Course) ([]Event, error) {
	if len(state.Students) > 10 {
		return nil, ErrTooManyStudents
	}
	if contains(state.Students, cmd.Name) {
		return nil, nil
	}
	return []Event{StudentAddedToCourse{Name: cmd.Name}}, nil
}

type RemoveStudentFromCourse struct {
	Name string
}

func (cmd RemoveStudentFromCourse) Decide(state Course) ([]Event, error) {
	if !contains(state.Students, cmd.Name) {
		return nil, nil
	}
	return []Event{StudentRemovedFromCourse{Name: cmd.Name}}, nil
}

func contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
