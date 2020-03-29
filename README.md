# goEventBus

event bus for Go's
+ High performance
+ asynchronous
+ non blocking

# Getting started

## Install

Import package:

```
go get github.com/newbietao/goEventBus/event
```

```
import (
    "github.com/newbietao/goEventBus/event"
)
```

# Usage

```
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
	e.TriggerEvent("sayHello", name)
}
```
See the complete example: [https://github.com/newbietao/goEventBus/tree/master/example](https://github.com/newbietao/goEventBus/tree/master/example)
