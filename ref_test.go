// Copyright 2018 vogo.
// Author: wongoo
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package hessian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type circularT struct {
	Num      int
	Previous *circularT
	Next     *circularT
}

func doTestRef(t *testing.T, c interface{}, name string) interface{} {
	bytes, err := ToBytes(c, NameMapFrom(c))
	assert.Nil(t, err)

	t.Logf("%s ref bytes: %s", name, string(bytes))
	t.Logf("%s ref bytes: %x", name, bytes)

	decoded, err := ToObject(bytes, TypeMapFrom(c))
	assert.Nil(t, err)
	t.Logf("%s ref decoded: %v", name, decoded)
	return decoded
}

func doTestCircularRef(t *testing.T, c *circularT, name string) *circularT {
	decoded := doTestRef(t, c, name)
	c2, ok := decoded.(*circularT)
	assert.True(t, ok)

	return c2
}

func TestBasicCircularRef(t *testing.T) {
	c1 := &circularT{}
	c1.Num = 12345
	c1.Previous = c1
	c1.Next = c1

	d1 := doTestCircularRef(t, c1, "basic")

	assert.Equal(t, c1.Num, d1.Num)
	assert.Equal(t, d1, d1.Previous)
	assert.Equal(t, d1, d1.Next)
}

func TestComplexCircularRef(t *testing.T) {
	c1 := &circularT{Num: 111}
	c2 := &circularT{Num: 222}
	c3 := &circularT{Num: 333}
	c4 := &circularT{Num: 444}

	c1.Previous = c4
	c1.Next = c2

	c2.Previous = c1
	c2.Next = c3

	c3.Previous = c2
	c3.Next = c4

	c4.Previous = c3
	c4.Next = c1

	d1 := doTestCircularRef(t, c1, "complex")
	d2 := d1.Next
	d3 := d2.Next
	d4 := d3.Next

	assert.Equal(t, c1.Num, d1.Num)
	assert.Equal(t, c2.Num, d2.Num)
	assert.Equal(t, c3.Num, d3.Num)
	assert.Equal(t, c4.Num, d4.Num)

	assert.True(t, AddrEqual(d4, d1.Previous))
	assert.True(t, AddrEqual(d1, d2.Previous))
	assert.True(t, AddrEqual(d2, d3.Previous))
	assert.True(t, AddrEqual(d3, d4.Previous))

	assert.True(t, AddrEqual(d2, d1.Next))
	assert.True(t, AddrEqual(d3, d2.Next))
	assert.True(t, AddrEqual(d4, d3.Next))
	assert.True(t, AddrEqual(d1, d4.Next))
}

type personsT []*personT

type personT struct {
	Name      string
	Likes     *personsT
	Relations []*personT
	Parent    *personT
	Marks     map[string]*personT
}

func TestComplexLevelRef(t *testing.T) {
	p1 := &personT{Name: "p1"}
	p2 := &personT{Name: "p2"}
	p3 := &personT{Name: "p3"}
	p4 := &personT{Name: "p4"}
	p5 := &personT{Name: "p5"}
	p6 := &personT{Name: "p6"}

	likes1 := &personsT{p2, p3}
	likes2 := &personsT{p4, p5, p6}

	p1.Likes = likes1
	p2.Likes = likes2

	p1.Parent = p2
	p2.Parent = p3

	relations := []*personT{p5, p6}
	p3.Relations = relations
	p4.Relations = relations

	marks := map[string]*personT{
		"beautiful": p1,
		"tall":      p2,
		"fat":       p3,
	}
	p4.Marks = marks
	p5.Marks = marks

	decoded := doTestRef(t, p1, "person")

	d1, ok := decoded.(*personT)
	assert.True(t, ok)

	d2 := d1.Parent
	assert.NotNil(t, d2)

	d3 := d2.Parent
	assert.NotNil(t, d3)

	assert.NotNil(t, d2.Likes)
	assert.Equal(t, 3, len(*d2.Likes))

	d4 := (*d2.Likes)[0]
	d5 := (*d2.Likes)[1]
	d6 := (*d2.Likes)[2]
	assert.NotNil(t, d4)
	assert.NotNil(t, d5)
	assert.NotNil(t, d6)

	assert.Equal(t, p1.Name, d1.Name)
	assert.Equal(t, p2.Name, d2.Name)
	assert.Equal(t, p3.Name, d3.Name)
	assert.Equal(t, p4.Name, d4.Name)
	assert.Equal(t, p5.Name, d5.Name)
	assert.Equal(t, p6.Name, d6.Name)

	assert.Equal(t, 2, len(*d1.Likes))
	assert.True(t, AddrEqual(d2, (*d1.Likes)[0]))
	assert.True(t, AddrEqual(d3, (*d1.Likes)[1]))

	assert.Equal(t, 2, len(d3.Relations))
	assert.True(t, AddrEqual(d5, d3.Relations[0]))
	assert.True(t, AddrEqual(d6, d3.Relations[1]))

	//assert.True(t, AddrEqual(d3.Relations, d4.Relations))
	//
	//assert.Equal(t, 3, len(d4.Marks))
	//assert.True(t, AddrEqual(p1, d4.Marks["beautiful"]))
	//assert.True(t, AddrEqual(p2, d4.Marks["tall"]))
	//assert.True(t, AddrEqual(p3, d4.Marks["fat"]))
	//assert.True(t, AddrEqual(d4.Marks, d5.Marks))
}
