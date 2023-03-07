package course

import (
	"encoding/json"
	"fmt"

	"github.com/IAmRadek/decider/v1"
)

type ES struct {
	events map[string][][]byte
}

func (es *ES) ReadStream(streamID string) (decider.Events[Course, int], error) {
	out := make([]Event, len(es.events[streamID]))

	var err error
	for i, e := range es.events[streamID] {
		out[i], err = unmarshalEvent(e)
		if err != nil {
			return decider.Events[Course, int]{}, err
		}
	}

	return decider.Events[Course, int]{Events: out}, nil
}

func (es *ES) AppendToStream(streamID string, version int, events []Event) error {
	for _, e := range events {
		bytes, err := marshalEvent(e)
		if err != nil {
			return err
		}

		es.events[streamID] = append(es.events[streamID], bytes)
	}

	return nil
}

func unmarshalEvent(bytes []byte) (Event, error) {
	var out struct {
		Type string
		Data json.RawMessage
	}
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}

	switch out.Type {
	case "StudentAddedToCourse":
		return unmarshalGeneric[StudentAddedToCourse](out.Data)
	case "StudentRemovedFromCourse":
		return unmarshalGeneric[StudentRemovedFromCourse](out.Data)
	default:
		return nil, fmt.Errorf("unknown event type: %s", out.Type)
	}
}

func unmarshalGeneric[T any](bytes []byte) (T, error) {
	var t T
	if err := json.Unmarshal(bytes, &t); err != nil {
		return t, err
	}
	return t, nil
}

func marshalEvent(v Event) ([]byte, error) {
	type outT struct {
		Type string
		Data Event
	}

	switch v := v.(type) {
	case StudentAddedToCourse:
		return json.Marshal(outT{Type: "StudentAddedToCourse", Data: v})
	case StudentRemovedFromCourse:
		return json.Marshal(outT{Type: "StudentRemovedFromCourse", Data: v})
	}
	return nil, fmt.Errorf("unknown event type: %T", v)
}
