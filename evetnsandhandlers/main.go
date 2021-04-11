package main

import (
	"EventSourcing/eventsource"
	"context"
	"fmt"
	"time"
)

//Create event strucuts
type OrderCreated struct {
	eventsource.Model
}

type OrderApproved struct {
	eventsource.Model
}

type OrderShipped struct {
	eventsource.Model
}

//Create command structs
type CreateOrder struct {
	eventsource.CommandModel
}

type ApproveOrder struct {
	eventsource.CommandModel
}

type ShipOrder struct {
	eventsource.CommandModel
}

type Order struct {
	ID        string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
}

func (order *Order) Apply(ctx context.Context, command eventsource.Command) ([]eventsource.Event, error) {
	switch cmd := command.(type) {
	case *CreateOrder:
		orderCreated := &OrderCreated{
			Model: eventsource.Model{ID: command.AggregateID(), Version: order.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCreated}, nil

	case *ApproveOrder:
		if order.State != "Created" {
			return nil, fmt.Errorf("Cannot approved order, %v, if it has not been created yet", command.AggregateID())
		}
		orderApproved := &OrderApproved{
			Model: eventsource.Model{ID: command.AggregateID(), Version: order.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderApproved}, nil

	case *ShipOrder:
		if order.State != "Approved" {
			return nil, fmt.Errorf("Cannot shipped order, %v, if it has not been approved yet", command.AggregateID())
		}
		orderShipped := &OrderShipped{
			Model: eventsource.Model{ID: command.AggregateID(), Version: order.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderShipped}, nil

	default:
		return nil, fmt.Errorf("Unknown Command: %v", cmd)
	}
}

func (order *Order) OnEvent(event eventsource.Event) {
	switch ev := event.(type) {
	case *OrderCreated:
		order.State = "Created"
	case *OrderApproved:
		order.State = "Approved"
	case *OrderShipped:
		order.State = "Shipped"
	default:
		fmt.Printf("Unknown Event: %T", ev)
	}

	order.ID = event.AggregateID()
	order.Version = event.EventVersion()
	order.UpdatedAt = event.EventAt()
}

func main() {
	id := "1"
	ctx := context.Background()
	order := &Order{}

	//Create the order
	createOrder := &CreateOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	events, _ := order.Apply(ctx, createOrder)

	for _, e := range events {
		order.OnEvent(e)
	}

	//Approve the order
	approveOrder := &ApproveOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	events, _ = order.Apply(ctx, approveOrder)

	for _, e := range events {
		order.OnEvent(e)
	}

	//Ship the order
	shipOrder := &ShipOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	events, _ = order.Apply(ctx, shipOrder)

	for _, e := range events {
		order.OnEvent(e)
	}

	fmt.Printf("The order state is %v and date is %v", order.State, order.UpdatedAt)
}
