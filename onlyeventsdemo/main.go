package main

import (
	"EventSourcing/eventsource"
	"fmt"
	"time"
)

type OrderCreated struct {
	eventsource.Model
}

type OrderApproved struct {
	eventsource.Model
}

type OrderShipped struct {
	eventsource.Model
}

type Order struct {
	ID        string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
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

	orderCreated := &OrderCreated{
		Model: eventsource.Model{ID: id, Version: 1, At: time.Now()},
	}
	orderAppoved := &OrderCreated{
		Model: eventsource.Model{ID: id, Version: 2, At: time.Now()},
	}
	orderShipped := &OrderCreated{
		Model: eventsource.Model{ID: id, Version: 3, At: time.Now()},
	}

	//Create the order
	order := &Order{}

	order.OnEvent(orderCreated)
	order.OnEvent(orderAppoved)
	order.OnEvent(orderShipped)

	fmt.Printf("The order state is %v and date is %v", order.State, order.UpdatedAt)
}
