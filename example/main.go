package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/newbietao/goEventBus/event"
)

// 实现event.EventHandle
type MyEvent struct {
}

// 事件前，可以做一些参数校验或者其他工作
func (m MyEvent) BeferEvent(i interface{}) error {
	log.Println("BeferEvent", i)
	if _, ok := i.(string); ok {
		return nil
	}
	return errors.New("name err")
}
func (m MyEvent) HandleEvent(i interface{}) error {
	log.Println("hello", i)
	return nil
}

func (m MyEvent) AfterEvent(i interface{}) error {
	log.Println("AfterEvent", i)
	return nil
}

func main() {

	e := event.GetEventBus()
	defer e.DestoryEventBus()

	myEvent := MyEvent{}
	e.RegisterEvent("sayHello", myEvent)
	// 为一个事件添加多个事件处理函数
	e.RegisterEvent("sayHello", myEvent)

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
