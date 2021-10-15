package pocket

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup  logging
	debug := false

	if debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	err := ForceUnlockDevices()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestGetReleaseHandle(t *testing.T) {

	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)
	err = releaseHandle(handle)
	assert.NoError(t, err)
}

func TestGetReasonableFrequency(t *testing.T) {

	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, to, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	fmt.Printf("Reasonable frequency range: [%d, %d]\n", from, to)

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestSingleQuery(t *testing.T) {

	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, _, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	s, err := singleQuery(handle, from, 1, SParamSelect{true, true, true, true})

	assert.NoError(t, err)

	fmt.Println(s)

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestRangeQuery(t *testing.T) {

	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, to, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	s, err := rangeQuery(handle, from, to, 3, 1, 1, SParamSelect{true, true, true, true})

	assert.NoError(t, err)

	fmt.Println(s)

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestRun(t *testing.T) {

	timeout := time.Millisecond * 100

	v := NewVNA()

	command := make(chan interface{})
	result := make(chan interface{})

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	go v.Run(command, result, ctx)

	id := "123xyz"
	command <- ReasonableFrequencyRange{Command: Command{ID: id}}

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-result:

		if actual, ok := ri.(ReasonableFrequencyRange); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			assert.True(t, actual.Result.Start > 0)
			assert.True(t, actual.Result.End > actual.Result.Start)
		}

	}
	cancel()

}
