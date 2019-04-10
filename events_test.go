// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package events

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert"
)

var (
	s1 Subscriber = func(data interface{}) {
		println("s1")
	}

	s2 Subscriber = func(data interface{}) {
		println("s2")
	}
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	p, e := New()
	a.NotNil(p).NotNil(e)
}

func TestPublisher_Publich(t *testing.T) {
	a := assert.New(t)
	p, e := New()
	a.NotNil(p).NotNil(e)

	buf1 := new(bytes.Buffer)
	sub1 := func(data interface{}) {
		buf1.Write(data.([]byte))
	}

	buf2 := new(bytes.Buffer)
	sub2 := func(data interface{}) {
		buf2.Write(data.([]byte))
	}

	id1 := e.Attach(sub1)
	p.Publish([]byte("p1"))
	time.Sleep(time.Microsecond * 500)
	a.NotError(a.Equal(buf1.String(), "p1"))
	a.Empty(buf2.Bytes())

	buf1.Reset()
	buf2.Reset()
	e.Attach(sub2)
	a.NotError(p.Publish([]byte("p2")))
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p2")
	a.Equal(buf2.String(), "p2")

	buf1.Reset()
	buf2.Reset()
	e.Detach(id1)
	a.NotError(p.Publish([]byte("p3")))
	time.Sleep(time.Microsecond * 500)
	a.Empty(buf1.String())
	a.Equal(buf2.String(), "p3")

	p.Destory()
	a.Error(p.Publish("p4"))
}

func TestPublisher_Destory(t *testing.T) {
	a := assert.New(t)

	p, e := New()
	a.NotNil(p).NotNil(e)
	p.Destory()
	a.Nil(p.e).
		Nil(e.subscribers)

	p, e = New()
	a.NotNil(p).NotNil(e)
	e.Attach(s1)
	p.Destory()
	a.Nil(p.e).
		Nil(e.subscribers)
}

func TestEvent_Attach_Detach(t *testing.T) {
	a := assert.New(t)
	p, e := New()
	a.NotNil(p).NotNil(e)

	id1 := e.Attach(s1)
	id2 := e.Attach(s2)
	a.Equal(len(e.subscribers), 2)

	e.Detach(id1)
	a.Equal(len(e.subscribers), 1)

	e.Detach(id2)
	a.Equal(len(e.subscribers), 0)
}
