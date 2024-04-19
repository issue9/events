// SPDX-FileCopyrightText: 2019-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package events 提供了简单的事件发布订阅功能
//
//	e := events.New[string]()
//
//	// 订阅事件
//	e.Subscribe(func(data string){
//	    fmt.Println("subscriber 1:", data)
//	})
//
//	// 订阅事件
//	e.Subscribe(func(data string){
//	    fmt.Println("subscriber 2:", data)
//	})
//
//	e.Publish(true, "test") // 发布事件
package events

import (
	"context"
	"reflect"
	"sync"
)

type (
	// SubscribeFunc 订阅者函数
	//
	// data 为事件传递过来的数据，可能存在多个订阅者，
	// 用户不应该直接修改 data 数据，否则结果是未知的。
	SubscribeFunc[T any] func(data T)

	// Publisher 事件的发布者
	Publisher[T any] interface {
		// Publish 触发事件
		//
		// sync 表示订阅者是否以异步的方式执行；
		// data 传递给订阅者的数据；
		Publish(sync bool, data T)
	}

	// Subscriber 供用户订阅事件的对象接口
	Subscriber[T any] interface {
		// Subscribe 注册订阅事件
		//
		// 返回用于注销此订阅事件的方法。
		Subscribe(SubscribeFunc[T]) context.CancelFunc
	}

	// Event 事件处理对象
	//
	// 同时实现了 [Subscriber] 和 [Publisher] 两个接口。
	Event[T any] struct {
		subscribers *sync.Map
	}
)

// New 声明一个新的事件处理
//
// T 为事件传递过程的参数类型；
func New[T any]() *Event[T] {
	return &Event[T]{
		subscribers: &sync.Map{},
	}
}

func (e *Event[T]) Publish(sync bool, data T) {
	if sync {
		e.subscribers.Range(func(key, value any) bool {
			go func(sub SubscribeFunc[T]) { sub(data) }(value.(SubscribeFunc[T]))
			return true
		})
	} else {
		e.subscribers.Range(func(key, value any) bool {
			value.(SubscribeFunc[T])(data)
			return true
		})
	}
}

func (e *Event[T]) Subscribe(subscriber SubscribeFunc[T]) context.CancelFunc {
	ptr := reflect.ValueOf(subscriber).Pointer()
	e.subscribers.Store(ptr, subscriber)
	return func() { e.subscribers.Delete(ptr) }
}

// Reset 重置对象
func (e *Event[T]) Reset() {
	e.subscribers.Range(func(key, _ any) bool {
		e.subscribers.Delete(key)
		return true
	})
}

// Len 订阅者的数量
func (e *Event[T]) Len() (c int) {
	e.subscribers.Range(func(key, value any) bool {
		c++
		return true
	})
	return
}
