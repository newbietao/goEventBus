package event

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/panjf2000/ants"
)

type EventHandle func(data interface{}) error
type EventName = string

type EventData struct {
	Name  string
	Param interface{}
}
type EventBus struct {
	eventHandles map[EventName][]EventHandle
	busChan      chan interface{}
	isLive       bool
	cancel       context.CancelFunc
}

// var chanSize = 0

var chanSize = 100
var poolSize = 100
var once = sync.Once{}

func GetEventBus() *EventBus {
	eCtx, eCancle := context.WithCancel(context.Background())
	e := &EventBus{
		eventHandles: make(map[string][]EventHandle),
		busChan:      make(chan interface{}, chanSize),
		isLive:       true,
		cancel:       eCancle,
	}
	go e.listenEvent(eCtx)
	return e
}

func (e *EventBus) listenEvent(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("listenEvent panic")
		}
		e.DestoryEventBus()
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
		if _, ok := e.eventHandles[en.Name]; !ok {
			err := errors.New("event name not existence")
			log.Printf("dispatchEvent err %v \n", err)
			return
		}
		var wg sync.WaitGroup
		for _, eventHandle := range e.eventHandles[en.Name] {
			wg.Add(1)
			go func(eventHandle EventHandle) {
				defer wg.Done()
				err := eventHandle(en.Param)
				if err != nil {
					log.Printf("eventHandle err %v \n", err)
				}
			}(eventHandle)
		}
		wg.Wait()
	}
}

// 推送事件
func (e *EventBus) pushEventBus(name string, param interface{}) {
	event := EventData{
		Name:  name,
		Param: param,
	}
	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		case e.busChan <- event:
			return
		case <-timer.C:
			log.Printf("pushEventBus %v \n", "done")
			return
		}
	}
}

func (e *EventBus) doDestory() {
	e.cancel()
	e.isLive = false
	close(e.busChan)
	for k, _ := range e.eventHandles {
		delete(e.eventHandles, k)
	}
	e.eventHandles = nil
}

// 注销event bus
func (e *EventBus) DestoryEventBus() {
	once.Do(e.doDestory)
}

// 注册事件，提供事件名和回调函数
func (e *EventBus) RegisterEvent(name string, event EventHandle) (err error) {
	if !e.isLive {
		err = errors.New("event bus not live")
		log.Printf("RegisterEvent err: %v \n", err)
		return
	}
	if _, ok := e.eventHandles[name]; ok {
		e.eventHandles[name] = append(e.eventHandles[name], event)
	} else {
		e.eventHandles[name] = []EventHandle{event}
	}
	return nil
}

// 触发事件
func (e *EventBus) TriggerEvent(name string, param interface{}) (err error) {
	// 检查event bus是否存活
	if !e.isLive {
		err = errors.New("event bus not live")
		log.Printf("TriggerEvent err: %v \n", err)
		return
	}
	// 检查事件是否注册过
	if _, ok := e.eventHandles[name]; !ok {
		err = errors.New("event name not existence")
		log.Printf("TriggerEvent err: %v \n", err)
		return
	}

	// 将事件发送到bus
	go e.pushEventBus(name, param)
	return nil
}
