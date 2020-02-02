// SPDX-License-Identifier: MIT

// Package events 提供了简单的事件发布订阅功能
package events

import (
	"errors"
	"sync"
)

// ErrStopped 表示发布都已经调用 Destory 销毁了事件处理器。
var ErrStopped = errors.New("该事件已经停止发布新内容")

// Subscriber 订阅者函数
//
// 当存在多个订阅者时，通过 go 异步执行每个函数。
//
// data 为事件传递过来的数据，可能存在多个订阅者，
// 最好不要直接修改 data 数据，否则结果是未知的。
type Subscriber func(data interface{})

// Event 事件
type Event struct {
	locker      sync.RWMutex
	count       int
	subscribers map[int]Subscriber
}

// Publisher 事件的发布者
type Publisher struct {
	e *Event
}

// New 声明一个新的事件处理
//
// Publisher 供事件发布者进行发布新事件；
// Event 供订阅者订阅事件。
func New() (*Publisher, *Event) {
	e := &Event{
		subscribers: make(map[int]Subscriber, 5),
	}

	p := &Publisher{
		e: e,
	}

	return p, e
}

// Publish 触发事件
func (p *Publisher) Publish(data interface{}) error {
	if p.e == nil {
		return ErrStopped
	}

	p.e.locker.RLock()
	for _, s := range p.e.subscribers {
		go func(sub Subscriber) {
			sub(data)
		}(s)
	}
	p.e.locker.RUnlock()

	return nil
}

// Destory 销毁当前事件处理程序
func (p *Publisher) Destory() {
	p.e.locker.Lock()
	p.e.subscribers = nil
	p.e.locker.Unlock()

	p.e = nil
}

// Attach 注册订阅者
//
// 返回一个唯一 ID，用户可以使用此 ID 取消订阅
func (e *Event) Attach(subscriber Subscriber) int {
	ret := e.count

	e.locker.Lock()
	e.count++
	e.subscribers[ret] = subscriber
	e.locker.Unlock()

	return ret
}

// Detach 取消订阅者
func (e *Event) Detach(id int) {
	e.locker.Lock()
	delete(e.subscribers, id)
	e.locker.Unlock()
}
