package eventsource

import "time"

//Events describe a change that happened to an Aggregate (Encapsulation of a group of domain objects)
//Always in PastTense e.g EmailChanged

type Event interface {
	//Returns the ID of the Aggregate referenced by the Event
	AggregateID() string

	//Returns the version of the Event
	EventVersion() int

	//Returns the time when the event happened
	EventAt() time.Time
}

//Model represents the implementation of the event.
type Model struct {
	//The aggregates ID
	ID string

	//The version of the event
	Version int

	//The time of the event
	At time.Time
}

//Make Model implements Event interface, implementing all it's functions.
func (model Model) AggregateID() string {
	return model.ID
}

func (model Model) EventVersion() int {
	return model.Version
}

func (model Model) EventAt() time.Time {
	return model.At
}
