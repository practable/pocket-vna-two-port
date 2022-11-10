package rfswitch

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHardware(t *testing.T) {

	if os.Getenv("VNA_HARDWAREPRESENT") != "TRUE" {
		t.Skip("Skipping hardware test (export VNA_HARDWAREPRESENT=TRUE to run test)")
	}

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
		err = r.SetThru()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetDUT1()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetDUT2()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetDUT3()
		assert.NoError(t, err)
		<-time.After(delay)
		err = r.SetDUT4()
		assert.NoError(t, err)
		<-time.After(delay)

	}

}
