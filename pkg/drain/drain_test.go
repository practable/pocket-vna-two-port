package drain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Str string
	Num int
}

func TestDrain(t *testing.T) {

	ch := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New(ch, ctx)

	n, err := s.NextNoWait()
	assert.Error(t, err)

	idx, err := s.LastReadIndex()
	assert.Error(t, err)
	assert.Equal(t, "empty store", err.Error())

	assert.Equal(t, 0, s.Count())

	a := foo{Str: "hello", Num: 0}
	b := foo{Str: "there", Num: 1}
	c := foo{Str: "friend", Num: 2}

	ch <- a
	ch <- b
	ch <- c

	assert.Equal(t, 3, s.Count())
	n, err = s.PeekLatest()
	assert.NoError(t, err)
	assert.Equal(t, c, n)

	idx, err = s.LastReadIndex()
	assert.Error(t, err)
	assert.Equal(t, "no reads yet", err.Error())
	assert.Equal(t, -1, idx)

	n, err = s.NextNoWait()
	assert.NoError(t, err)
	assert.Equal(t, a, n)

	n, err = s.PeekLastRead()
	assert.NoError(t, err)
	assert.Equal(t, a, n)

	n, err = s.NextNoWait()
	assert.NoError(t, err)
	assert.Equal(t, b, n)

	n, err = s.NextNoWait()
	assert.NoError(t, err)
	assert.Equal(t, c, n)

	n, err = s.NextNoWait()
	assert.Error(t, err)

	all := s.All()
	assert.Equal(t, []interface{}{a, b, c}, all)

}
