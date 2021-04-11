package eventsource

import "context"

//Command Encapsulates the ID of the aggregate to mutate
type Command interface {
	AggregateID() string
}

//CommandHandler consume a command and produces an event
type CommandHandler interface {
	Apply(ctx context.Context, command Command) ([]Event, error)
}

//Command model to represent a command
type CommandModel struct {
	ID string
}

//Implements the Command interface
func (m CommandModel) AggregateID() string {
	return m.ID
}
