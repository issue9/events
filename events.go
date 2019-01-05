// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package events

// Subscriber 订阅者函数
type Subscriber func(data interface{})

// Event 事件
type Event struct {
	subscribers []Subscriber
}

// Notify 触发事件
func (e *Event) Notify(data interface{}) {
	for _, subscriber := range e.subscribers {
		subscriber(data)
	}
}

// New 创建新的事件
func New(name string) *Event {
	return &Event{
		subscribers: make([]Subscriber, 0, 5),
	}
}

// Attach 注册订阅者
func (e *Event) Attach(name string, subscriber Subscriber) {
	e.subscribers = append(e.subscribers, subscriber)
}
