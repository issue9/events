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

// SubscribeFunc 订阅者函数
//
// data 为事件传递过来的数据，可能存在多个订阅者，
// 用户不应该直接修改 data 数据，否则结果是未知的。
type SubscribeFunc[T any] func(data T)

type event[T any] struct {
	funcs *sync.Map
}

// Publisher 事件的发布者
type Publisher[T any] interface {
	// Publish 触发事件
	//
	// sync 表示订阅者是否以异步的方式执行；
	// data 传递给订阅者的数据；
	Publish(sync bool, data T)
}

// Subscriber 供用户订阅事件的对象接口
type Subscriber[T any] interface {
	// Subscribe 注册订阅事件
	//
	// 返回用于注销此订阅事件的方法。
	Subscribe(SubscribeFunc[T]) (context.CancelFunc, error)
}

type Eventer[T any] interface {
	Publisher[T]
	Subscriber[T]
}

// New 声明一个新的事件处理
//
// T 为事件传递过程的参数类型；
func New[T any]() Eventer[T] {
	return &event[T]{
		funcs: &sync.Map{},
	}
}

func (e *event[T]) Publish(sync bool, data T) {
	if sync {
		e.funcs.Range(func(key, value any) bool {
			go func(sub SubscribeFunc[T]) { sub(data) }(value.(SubscribeFunc[T]))
			return true
		})
	} else {
		e.funcs.Range(func(key, value any) bool {
			value.(SubscribeFunc[T])(data)
			return true
		})
	}
}

func (e *event[T]) Subscribe(subscriber SubscribeFunc[T]) (context.CancelFunc, error) {
	ptr := reflect.ValueOf(subscriber).Pointer()
	e.funcs.Store(ptr, subscriber)

	return func() { e.funcs.Delete(ptr) }, nil
}
