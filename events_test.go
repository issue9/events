// SPDX-License-Identifier: MIT

package events

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
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
	a := assert.New(t, false)

	p, e := New()
	a.NotNil(p).NotNil(e)
}

func TestPublisher_Publish(t *testing.T) {
	a := assert.New(t, false)
	p, e := New()
	a.NotNil(p).NotNil(e)

	// 没有订阅者
	a.NotError(p.Publish(true, []byte("123")))

	buf1 := new(bytes.Buffer)
	sub1 := func(data interface{}) {
		buf1.Write(data.([]byte))
	}

	buf2 := new(bytes.Buffer)
	sub2 := func(data interface{}) {
		buf2.Write(data.([]byte))
	}

	id1, err := e.Attach(sub1)
	a.NotError(err)
	p.Publish(true, []byte("p1"))
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p1")
	a.Empty(buf2.Bytes())

	buf1.Reset()
	buf2.Reset()
	e.Attach(sub2)
	a.NotError(p.Publish(false, []byte("p2")))
	time.Sleep(time.Microsecond * 500)
	a.Equal(buf1.String(), "p2")
	a.Equal(buf2.String(), "p2")

	buf1.Reset()
	buf2.Reset()
	e.Detach(id1)
	a.NotError(p.Publish(false, []byte("p3")))
	time.Sleep(time.Microsecond * 500)
	a.Empty(buf1.String())
	a.Equal(buf2.String(), "p3")

	p.Destroy()
	a.Error(p.Publish(false, "p4"))
}

func TestPublisher_Destroy(t *testing.T) {
	a := assert.New(t, false)

	p, e := New()
	a.NotNil(p).NotNil(e)
	p.Destroy()
	ee, ok := e.(*event)
	a.True(ok).NotNil(ee).Nil(ee.subscribers)

	p, e = New()
	a.NotNil(p).NotNil(e)
	e.Attach(s1)
	p.Destroy()
	ee, ok = e.(*event)
	a.True(ok).NotNil(ee).Nil(ee.subscribers)
}

func TestEventer_Attach_Detach(t *testing.T) {
	a := assert.New(t, false)
	p, e := New()
	a.NotNil(p).NotNil(e)

	id1, err := e.Attach(s1)
	a.NotError(err)
	id2, err := e.Attach(s2)
	a.NotError(err)
	ee, ok := e.(*event)
	a.True(ok).NotNil(ee)

	a.Equal(len(ee.subscribers), 2)

	e.Detach(id1)
	a.Equal(len(ee.subscribers), 1)

	e.Detach(id2)
	a.Equal(len(ee.subscribers), 0)

	// Destroy

	p.Destroy()
	e.Attach(s1)
}
