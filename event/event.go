package event

import (
	"context"
	"errors"
	"log"

	"github.com/panjf2000/ants"
)

type EventHandle interface {
	BeferEvent(i interface{}) error
	HandleEvent(i interface{}) error
	AfterEvent(i interface{}) error
}

type EventData struct {
	Name  string
	Param interface{}
}

type EventBus struct {
	eventByName map[string][]EventHandle
	busChan     chan EventData
	isLive      bool
}

var chanSize = 100
var poolSize = 100

func GetEventBus(ctx context.Context) *EventBus {
	e := &EventBus{
		eventByName: make(map[string][]EventHandle),
		busChan:     make(chan EventData, chanSize),
		isLive:      true,
	}
	go e.listenEvent(ctx)
	return e
}

func (e *EventBus) listenEvent(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("listenEvent panic")
		}
		e.isLive = false
	}()

	p, _ := ants.NewPoolWithFunc(int(poolSize), e.dispatchEvent)
	defer p.Release()

	for {
		select {
		case event := <-e.busChan:
			err := p.Invoke(event)
			if err != nil {
				log.Printf("listenEvent-> %s \n", err)
			}
		case <-ctx.Done():
			log.Printf("listenEvent-> %s \n", "Done")
			return
		}
	}
}

func (e *EventBus) dispatchEvent(eve interface{}) {
	if en, ok := eve.(EventData); ok {
		if _, ok := e.eventByName[en.Name]; !ok {
			err := errors.New("event name not existence")
			log.Printf("dispatchEvent err %v \n", err)
			return
		}
		for _, event := range e.eventByName[en.Name] {
			err := event.HandleEvent(en.Param)
			if err != nil {
				log.Printf("dispatchEvent err %v \n", err)
				return
			}
			err = event.AfterEvent(en.Param)
			if err != nil {
				log.Printf("dispatchEvent err %v \n", err)
				return
			}
		}
	}
}

// 推送事件
func (e *EventBus) pushEventBus(name string, param interface{}) {
	event := EventData{
		Name:  name,
		Param: param,
	}
	e.busChan <- event
}

// 注册事件，提供事件名和回调函数
func (e *EventBus) RegisterEvent(name string, event EventHandle) (err error) {
	if !e.isLive {
		return errors.New("event bus not live")
	}
	if _, ok := e.eventByName[name]; ok {
		e.eventByName[name] = append(e.eventByName[name], event)
	} else {
		e.eventByName[name] = []EventHandle{event}
	}
	return nil
}

// 触发事件
func (e *EventBus) TriggerEvent(name string, param interface{}) (err error) {
	// 检查event bus是否存活
	if !e.isLive {
		return errors.New("event bus not live")
	}
	// 检查事件是否注册过
	if _, ok := e.eventByName[name]; !ok {
		return errors.New("event name not existence")
	}
	// 同步执行所有BeferEvent钩子函数，有可能涉及到参数校验、初始化等，所以要同步执行
	for _, event := range e.eventByName[name] {
		err := event.BeferEvent(param)
		if err != nil {
			return err
		}
	}

	// 将事件发送到bus
	go e.pushEventBus(name, param)
	return nil
}
