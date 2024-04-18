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
	s1 SubscribeFunc[string] = func(data string) {
		println("s1")
	}

	s2 SubscribeFunc[string] = func(data string) {
		println("s2")
	}
)

func TestPublisher_Publish(t *testing.T) {
	a := assert.New(t, false)
	e := New[string]()
	a.NotNil(e)

	// 没有订阅者
	e.Publish(true, "123")

	buf1 := new(bytes.Buffer)
	sub1 := func(data string) {
		buf1.WriteString(data)
	}

	c1, err := e.Subscribe(sub1)
	a.NotError(err).NotNil(c1)
	e.Publish(true, "p1")
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p1")

	buf1.Reset()
	buf2 := new(bytes.Buffer)
	sub2 := func(data string) {
		buf2.WriteString(data)
	}
	a.Empty(buf2.Bytes())
	e.Subscribe(sub2)
	e.Publish(false, "p2")
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p2")
	a.Equal(buf2.String(), "p2")

	buf1.Reset()
	buf2.Reset()
	c1()
	e.Publish(false, "p3")
	time.Sleep(time.Microsecond * 500)
	a.Empty(buf1.String())
	a.Equal(buf2.String(), "p3")
}

func TestPublisher_Destroy(t *testing.T) {
	a := assert.New(t, false)

	e := New[string]()
	a.NotNil(e)
	ee, ok := e.(*(event[string]))
	a.True(ok).NotNil(ee).Zero(ee.len())

	e = New[string]()
	a.NotNil(e)
	e.Subscribe(s1)
	ee, ok = e.(*(event[string]))
	a.True(ok).NotNil(ee).Equal(ee.len(), 1)
}

func TestSubscriber_Attach_Detach(t *testing.T) {
	a := assert.New(t, false)
	e := New[string]()
	a.NotNil(e)

	c1, err := e.Subscribe(s1)
	a.NotError(err)
	c2, err := e.Subscribe(s2)
	a.NotError(err)
	ee, ok := e.(*(event[string]))
	a.True(ok).NotNil(ee)

	a.Equal(ee.len(), 2)

	c1()
	a.Equal(ee.len(), 1)

	c2()
	a.Equal(ee.len(), 0)
}

func (e *event[T]) len() (c int) {
	e.funcs.Range(func(key, value any) bool {
		c++
		return true
	})
	return
}
