package course

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/IAmRadek/decider/v1"
)

func TestCourse(t *testing.T) {
	d := decider.Decider[Course]{
		Evolve: evolve,
		Decide: decide,
		GetInitialState: func() Course {
			return Course{}
		},
	}

	es := &ES{
		events: map[string][][]byte{},
	}

	eh := decider.CommandHandler[Course, string, int]{
		EventStore: es,
		Decider:    d,
	}

	evts, err := eh.HandleCommand("course_v2-1", AddStudentToCourse{Name: "John"})
	if err != nil {
		t.Fatal(err)
	}
	pp(evts)

	evts, err = eh.HandleCommand("course_v2-1", AddStudentToCourse{Name: "John2"})
	if err != nil {
		t.Fatal(err)
	}
	pp(evts)

	events, err := es.ReadStream("course_v2-1")
	if err != nil {
		t.Fatal(err)
	}
	pp(events)

	evts, err = eh.HandleCommand("course_v2-1", RemoveStudentFromCourse{Name: "John2"})
	if err != nil {
		t.Fatal(err)
	}
	pp(evts)
	println("----")
	for _, e := range es.events["course_v2-1"] {
		fmt.Println(string(e))
	}
}

func pp(d interface{}) {
	dd, _ := json.MarshalIndent(d, "", "\t")
	fmt.Println(string(dd))
}
