package main

import (
	"context"
	"log"

	"github.com/newbietao/goEventBus/event"
)

type MyEvent struct {
}

func (m MyEvent) BeferEvent(i interface{}) error {
	log.Println("BeferEvent", i)
	return nil
}
func (m MyEvent) HandleEvent(i interface{}) error {
	log.Println("HandleEvent", i)
	return nil
}

func (m MyEvent) AfterEvent(i interface{}) error {
	log.Println("AfterEvent", i)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e := event.GetEventBus(ctx)
	myEvent := MyEvent{}
	e.RegisterEvent("myEvent", myEvent)
	e.TriggerEvent("myEvent", "test1")
	e.TriggerEvent("myEvent", "test2")
}
