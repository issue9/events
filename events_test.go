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
	a.NotError(e.Publish(true, "123"))

	buf1 := new(bytes.Buffer)
	sub1 := func(data string) {
		buf1.WriteString(data)
	}

	buf2 := new(bytes.Buffer)
	sub2 := func(data string) {
		buf2.WriteString(data)
	}

	id1, err := e.Attach(sub1)
	a.NotError(err)
	e.Publish(true, "p1")
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p1")
	a.Empty(buf2.Bytes())

	buf1.Reset()
	buf2.Reset()
	e.Attach(sub2)
	a.NotError(e.Publish(false, "p2"))
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p2")
	a.Equal(buf2.String(), "p2")

	buf1.Reset()
	buf2.Reset()
	e.Detach(id1)
	a.NotError(e.Publish(false, "p3"))
	time.Sleep(time.Microsecond * 500)
	a.Empty(buf1.String())
	a.Equal(buf2.String(), "p3")

	e.Destroy()
	a.Error(e.Publish(false, "p4"))
}

func TestPublisher_Destroy(t *testing.T) {
	a := assert.New(t, false)

	e := New[string]()
	a.NotNil(e)
	e.Destroy()
	ee, ok := e.(*(event[string]))
	a.True(ok).NotNil(ee).Nil(ee.funcs)

	e = New[string]()
	a.NotNil(e)
	e.Attach(s1)
	e.Destroy()
	ee, ok = e.(*(event[string]))
	a.True(ok).NotNil(ee).Nil(ee.funcs)
}

func TestSubscriber_Attach_Detach(t *testing.T) {
	a := assert.New(t, false)
	e := New[string]()
	a.NotNil(e)

	id1, err := e.Attach(s1)
	a.NotError(err)
	id2, err := e.Attach(s2)
	a.NotError(err)
	ee, ok := e.(*(event[string]))
	a.True(ok).NotNil(ee)

	a.Equal(len(ee.funcs), 2)

	e.Detach(id1)
	a.Equal(len(ee.funcs), 1)

	e.Detach(id2)
	a.Equal(len(ee.funcs), 0)

	// Destroy

	e.Destroy()
	e.Attach(s1)
}
