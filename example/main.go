package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/newbietao/goEventBus/event"
)

// 实现event.EventHandle

// 事件前，可以做一些参数校验或者其他工作
func BeferEvent(i interface{}) error {
	log.Println("BeferEvent", i)
	if _, ok := i.(string); ok {
		return nil
	}
	return errors.New("name err")
}
func HandleEvent(i interface{}) error {
	log.Println("hello", i)
	return nil
}

func AfterEvent(i interface{}) error {
	log.Println("AfterEvent", i)
	return nil
}

func main() {

	e := event.GetEventBus()
	defer e.DestoryEventBus()

	e.RegisterEvent("sayHello", BeferEvent)
	e.RegisterEvent("sayHello", HandleEvent)
	e.RegisterEvent("sayHello", AfterEvent)

	name := ""
	for {
		fmt.Println("pleace entry name:")
		fmt.Scanf("%s", &name)
		if name == "end" {
			return
		} else if name == "cancle" {
			e.DestoryEventBus()
		} else {
			e.TriggerEvent("sayHello", name)
		}

		time.Sleep(2 * time.Second)
	}
}
