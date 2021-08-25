// SPDX-License-Identifier: MIT

// Package events 提供了简单的事件发布订阅功能
//
//  p, e := events.New()
//
//  // 订阅事件
//  e.Attach(func(data interface{}){
//      fmt.Println("subscriber 1")
//  })
//
//  // 订阅事件
//  e.Attach(func(data interface{}){
//      fmt.Println("subscriber 2")
//  })
//
//  p.Publish(true, nil) // 发布事件
package events

import (
	"errors"
	"sync"
)

// ErrStopped 表示发布都已经调用 Destory 销毁了事件处理器
var ErrStopped = errors.New("该事件已经停止发布新内容")

// Subscriber 订阅者函数
//
// 每个订阅函数都是通过 go 异步执行。
//
// data 为事件传递过来的数据，可能存在多个订阅者，
// 用户不应该直接修改 data 数据，否则结果是未知的。
type Subscriber func(data interface{})

type event struct {
	locker      sync.RWMutex
	count       int
	subscribers map[int]Subscriber
}

// Publisher 事件的发布者
type Publisher interface {
	// 触发事件
	//
	// sync 表示订阅者是否以异步的方式执行；
	// data 传递给订阅者的数据；
	Publish(sync bool, data interface{}) error

	// 销毁当前事件处理程序
	Destory()
}

// Eventer 供用户订阅事件的对象接口
type Eventer interface {
	// 注册订阅者
	//
	// 返回唯一 ID，用户可以使用此 ID 取消订阅。
	Attach(Subscriber) (int, error)

	// 取消指定事件的订阅
	Detach(int)
}

// New 声明一个新的事件处理
//
// Publisher 供事件发布者进行发布新事件；
// Event 供订阅者订阅事件。
func New() (Publisher, Eventer) {
	e := &event{
		subscribers: make(map[int]Subscriber, 5),
	}

	return e, e
}

func (e *event) Publish(sync bool, data interface{}) error {
	if e.subscribers == nil { // 初如化时将 subscribers 设置为了 5，所以为 nil 表示已经调用 Destory
		return ErrStopped
	}

	e.locker.RLock()
	defer e.locker.RUnlock()

	if len(e.subscribers) == 0 {
		return nil
	}

	if sync {
		for _, s := range e.subscribers {
			go func(sub Subscriber) {
				sub(data)
			}(s)
		}
	} else {
		for _, s := range e.subscribers {
			s(data)
		}
	}

	return nil
}

func (e *event) Destory() {
	e.locker.Lock()
	e.subscribers = nil
	e.locker.Unlock()
}

func (e *event) Attach(subscriber Subscriber) (int, error) {
	if e.subscribers == nil {
		return 0, ErrStopped
	}

	ret := e.count

	e.locker.Lock()
	e.count++
	e.subscribers[ret] = subscriber
	e.locker.Unlock()

	return ret, nil
}

func (e *event) Detach(id int) {
	e.locker.Lock()
	delete(e.subscribers, id)
	e.locker.Unlock()
}
