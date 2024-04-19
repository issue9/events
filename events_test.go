// SPDX-FileCopyrightText: 2019-2024 caixw
//
// SPDX-License-Identifier: MIT

package events

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
)

var (
	_ Publisher[int]  = &Event[int]{}
	_ Subscriber[int] = &Event[int]{}
)

func s1(data string) { println("s1") }

func s2(data string) { println("s2") }

func TestPublisher_Publish(t *testing.T) {
	a := assert.New(t, false)
	e := New[string]()
	a.NotNil(e)

	// 没有订阅者
	e.Publish(true, "123")

	buf1 := new(bytes.Buffer)
	sub1 := func(data string) { buf1.WriteString(data) }

	c1 := e.Subscribe(sub1)
	a.NotNil(c1)
	e.Publish(true, "p1")
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p1")

	buf1.Reset()
	buf2 := new(bytes.Buffer)
	sub2 := func(data string) { buf2.WriteString(data) }
	a.Empty(buf2.Bytes())
	e.Subscribe(sub2)
	e.Publish(false, "p2")
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p2").
		Equal(buf2.String(), "p2")

	buf1.Reset()
	buf2.Reset()
	c1()
	e.Publish(false, "p3")
	time.Sleep(time.Microsecond * 500)
	a.Empty(buf1.String())
	a.Equal(buf2.String(), "p3")

}

func TestPublisher_Reset(t *testing.T) {
	a := assert.New(t, false)

	e := New[string]()
	a.NotNil(e)
	a.Zero(e.Len())

	e = New[string]()
	a.NotNil(e)
	e.Subscribe(s1)
	a.Equal(e.Len(), 1)

	e.Reset()
	a.Zero(e.Len())
}

func TestSubscriber_Attach_Detach(t *testing.T) {
	a := assert.New(t, false)
	e := New[string]()
	a.NotNil(e)

	c1 := e.Subscribe(s1)
	c2 := e.Subscribe(s2)

	a.Equal(e.Len(), 2)

	c1()
	a.Equal(e.Len(), 1)

	c2()
	a.Equal(e.Len(), 0)
}
