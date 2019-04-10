// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package events

import (
	"testing"

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

func TestPublisher_Destory(t *testing.T) {
	a := assert.New(t)

	p, e := New()
	a.NotNil(p).NotNil(e)

	p.Destory()
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
