// SPDX-License-Identifier: MIT

// Package events 提供了简单的事件发布订阅功能
//
//	e := events.New[string]()
//
//	// 订阅事件
//	e.Attach(func(data string){
//	    fmt.Println("subscriber 1:", data)
//	})
//
//	// 订阅事件
//	e.Attach(func(data string){
//	    fmt.Println("subscriber 2:", data)
//	})
//
//	e.Publish(true, "test") // 发布事件
package events

import (
	"errors"
	"sync"
)

var errStopped = errors.New("events: stopped")

// SubscribeFunc 订阅者函数
//
// 每个订阅函数都是通过 go 异步执行。
//
// data 为事件传递过来的数据，可能存在多个订阅者，
// 用户不应该直接修改 data 数据，否则结果是未知的。
type SubscribeFunc[T any] func(data T)

type event[T any] struct {
	locker sync.RWMutex
	count  int
	funcs  map[int]SubscribeFunc[T]
}

// Publisher 事件的发布者
type Publisher[T any] interface {
	// Publish 触发事件
	//
	// sync 表示订阅者是否以异步的方式执行；
	// data 传递给订阅者的数据；
	Publish(sync bool, data T) error

	// Destroy 销毁当前事件处理程序
	Destroy()
}

// Subscriber 供用户订阅事件的对象接口
type Subscriber[T any] interface {
	// Attach 注册订阅者
	//
	// 返回唯一 ID，用户可以使用此 ID 取消订阅。
	Attach(SubscribeFunc[T]) (int, error)

	// Detach 取消指定事件的订阅
	Detach(int)
}

type Eventer[T any] interface {
	Publisher[T]
	Subscriber[T]
}

// errStopped 表示发布都已经调用 [Publisher.Destroy] 销毁了事件处理器
func ErrStopped() error { return errStopped }

// New 声明一个新的事件处理
//
// T 为事件传递过程的参数类型；
func New[T any]() Eventer[T] {
	return &event[T]{
		funcs: make(map[int]SubscribeFunc[T], 5),
	}
}

func (e *event[T]) Publish(sync bool, data T) error {
	// 初如化时将 e.funcs 设置为了非 nil 状态，
	// 所以为 nil 表示已经调用 [Publisher.Destroy]
	if e.funcs == nil {
		return ErrStopped()
	}

	e.locker.RLock()
	defer e.locker.RUnlock()

	if len(e.funcs) == 0 {
		return nil
	}

	if sync {
		for _, s := range e.funcs {
			go func(sub SubscribeFunc[T]) {
				sub(data)
			}(s)
		}
	} else {
		for _, s := range e.funcs {
			s(data)
		}
	}

	return nil
}

func (e *event[T]) Destroy() {
	e.locker.Lock()
	e.funcs = nil
	e.locker.Unlock()
}

func (e *event[T]) Attach(subscriber SubscribeFunc[T]) (int, error) {
	if e.funcs == nil {
		return 0, ErrStopped()
	}

	ret := e.count

	e.locker.Lock()
	e.count++
	e.funcs[ret] = subscriber
	e.locker.Unlock()

	return ret, nil
}

func (e *event[T]) Detach(id int) {
	e.locker.Lock()
	delete(e.funcs, id)
	e.locker.Unlock()
}
