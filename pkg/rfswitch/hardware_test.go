package rfswitch

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHardware(t *testing.T) {

	ur := "ws://localhost:8888/ws/rfswitch"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := New(ur, ctx)

	delay := 100 * time.Millisecond

	for i := 0; i < 2; i++ {

		err := r.SetShort()
		assert.NoError(t, err)

		<-time.After(delay)
		err = r.SetOpen()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetLoad()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetDUT()
		assert.NoError(t, err)
		<-time.After(delay)
	}

}
