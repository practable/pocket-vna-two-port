package pocket

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var verbose bool
var hardware bool

func TestMain(m *testing.M) {
	// Setup  logging
	debug := false
	verbose = true
	hardware = false

	if debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	if hardware {
		err := ForceUnlockDevices()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestMockConnect(t *testing.T) {

	v := NewMock()

	disconnect, err := v.Connect()
	assert.NoError(t, err)
	err = disconnect()
	assert.NoError(t, err)

	ce := errors.New("PVNA_Res_NoDevice")
	v.ConnectError = ce
	disconnect, err = v.Connect()
	assert.Error(t, err)
	assert.Equal(t, ce, err)
	err = disconnect()
	assert.NoError(t, err)

}

func TestMockGetReasonableFrequencyRange(t *testing.T) {

	v := NewMock()

	start := uint64(1000000)
	end := uint64(4000000000)
	id := "rfr00"

	v.ResultReasonableFrequencyRange = Range{
		Start: start,
		End:   end,
	}

	c := ReasonableFrequencyRange{Command: Command{ID: id}}

	err := v.GetReasonableFrequencyRange(&c)

	assert.NoError(t, err)

	assert.Equal(t, id, c.ID)
	assert.Equal(t, start, c.Result.Start)
	assert.Equal(t, end, c.Result.End)

	assert.Equal(t, []interface{}{c}, v.CommandsReceived)

}

func TestMockSingleQuery(t *testing.T) {

	v := NewMock()

	id := "sq00"

	v.ResultSingleQuery = SParam{
		S11: Complex{
			Imag: 1.0,
			Real: 2.0,
		},
		Freq: 3,
	}

	// our mock doesn't check anything in the command
	c := SingleQuery{Command: Command{ID: id}}

	err := v.SingleQuery(&c)

	assert.NoError(t, err)

	assert.Equal(t, id, c.ID)
	assert.Equal(t, 1.0, c.Result.S11.Imag)
	assert.Equal(t, 2.0, c.Result.S11.Real)
	assert.Equal(t, uint64(3), c.Result.Freq)
	assert.Equal(t, []interface{}{c}, v.CommandsReceived)
}

func TestMockRangeQuery(t *testing.T) {

	v := NewMock()

	id := "sq00"

	v.ResultRangeQuery = []SParam{SParam{
		S11: Complex{
			Imag: 1.0,
			Real: 2.0,
		},
		Freq: 3,
	},
		SParam{
			S11: Complex{
				Imag: 4.0,
				Real: 5.0,
			},
			Freq: 6,
		},
	}

	// our mock doesn't check anything in the command
	c := RangeQuery{Command: Command{ID: id}}

	err := v.RangeQuery(&c)

	assert.NoError(t, err)

	assert.Equal(t, id, c.ID)
	assert.Equal(t, 1.0, c.Result[0].S11.Imag)
	assert.Equal(t, 2.0, c.Result[0].S11.Real)
	assert.Equal(t, uint64(3), c.Result[0].Freq)
	assert.Equal(t, 4.0, c.Result[1].S11.Imag)
	assert.Equal(t, 5.0, c.Result[1].S11.Real)
	assert.Equal(t, uint64(6), c.Result[1].Freq)

	assert.Equal(t, []interface{}{c}, v.CommandsReceived)
}

func TestMockHandleCommand(t *testing.T) {

	v := NewMock()

	start := uint64(1000000)
	end := uint64(4000000000)

	v.ResultReasonableFrequencyRange = Range{
		Start: start,
		End:   end,
	}

	v.ResultSingleQuery = SParam{
		S11: Complex{
			Imag: 1.0,
			Real: 2.0,
		},
		Freq: 3,
	}

	v.ResultRangeQuery = []SParam{SParam{
		S11: Complex{
			Imag: 4.0,
			Real: 5.0,
		},
		Freq: 6,
	},
		SParam{
			S11: Complex{
				Imag: 7.0,
				Real: 8.0,
			},
			Freq: 9,
		},
	}

	id0 := "rfr00"
	c0 := ReasonableFrequencyRange{Command: Command{ID: id0, Command: "rfr"}}

	err := v.HandleCommand(&c0)
	assert.NoError(t, err)
	assert.Equal(t, id0, c0.ID)
	assert.Equal(t, start, c0.Result.Start)
	assert.Equal(t, end, c0.Result.End)

	id1 := "rq00"
	c1 := RangeQuery{Command: Command{ID: id1, Command: "rq"}}
	err = v.HandleCommand(&c1)
	assert.NoError(t, err)
	assert.Equal(t, id1, c1.ID)
	assert.Equal(t, 4.0, c1.Result[0].S11.Imag)
	assert.Equal(t, 5.0, c1.Result[0].S11.Real)
	assert.Equal(t, uint64(6), c1.Result[0].Freq)
	assert.Equal(t, 7.0, c1.Result[1].S11.Imag)
	assert.Equal(t, 8.0, c1.Result[1].S11.Real)
	assert.Equal(t, uint64(9), c1.Result[1].Freq)

	id2 := "sq00"
	c2 := SingleQuery{Command: Command{ID: id2, Command: "sq"}}
	err = v.HandleCommand(&c2)
	assert.NoError(t, err)
	assert.Equal(t, id2, c2.ID)
	assert.Equal(t, 1.0, c2.Result.S11.Imag)
	assert.Equal(t, 2.0, c2.Result.S11.Real)
	assert.Equal(t, uint64(3), c2.Result.Freq)

	assert.Equal(t, []interface{}{c0, c1, c2}, v.CommandsReceived)
}

/* TODO refactor tests
func TestGetReleaseHandleHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)
	err = releaseHandle(handle)
	assert.NoError(t, err)
}

func TestGetReasonableFrequencyHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, to, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	if verbose {
		fmt.Printf("Reasonable frequency range: [%d, %d]\n", from, to)
	}

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestSingleQueryHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, _, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	s, err := singleQuery(handle, from, 1, SParamSelect{true, true, true, true})

	assert.NoError(t, err)

	if verbose {
		fmt.Println(s)
	}

	assert.Equal(t, from, s.Freq)

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestRangeQueryHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)

	from, to, err := getReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	s, err := rangeQuery(handle, from, to, 3, 1, 1, SParamSelect{true, true, true, true})

	assert.NoError(t, err)

	if verbose {
		fmt.Println(s)
	}

	//make int for easier interpretation of debug info during test, if wrong
	assert.Equal(t, int(from), int(s[0].Freq))
	assert.Equal(t, int(from+(to-from)/2), int(s[1].Freq))
	assert.Equal(t, int(to), int(s[2].Freq))

	err = releaseHandle(handle)
	assert.NoError(t, err)

}

func TestNewServiceHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	timeout := time.Millisecond * 100

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	v := New(ctx)

	// Do GetReasonableFrequencyRange command

	reasonable := Range{}

	id := "123xyz"
	v.Request <- ReasonableFrequencyRange{Command: Command{ID: id}}

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-v.Response:

		if actual, ok := ri.(ReasonableFrequencyRange); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			// weak test - with real kit attached, we should get non-zero numbers
			assert.True(t, actual.Result.Start > 0)
			assert.True(t, actual.Result.End > actual.Result.Start)
			reasonable = actual.Result //save for RangeQuery
			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	// Do SingleQuery command

	id = "456abc"
	v.Request <- SingleQuery{
		Command: Command{ID: id},
		Freq:    200000,
		Avg:     1,
		Select:  SParamSelect{true, true, true, true},
	}

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-v.Response:

		if actual, ok := ri.(SingleQuery); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			// weak test - with real kit attached, we should get non-zero numbers
			assert.True(t, actual.Result.S11.Real != 0)
			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	// Do RangeQuery command

	id = "789def"
	N := 7 // number of samples
	v.Request <- RangeQuery{
		Command:         Command{ID: id},
		Range:           reasonable,
		Size:            N,
		Avg:             1,
		LogDistribution: true,
		Select:          SParamSelect{true, true, true, true},
	}

	timeout = time.Second //need more time for this than a single query

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-v.Response:

		if actual, ok := ri.(RangeQuery); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			// weak test - with real kit attached, we should get non-zero numbers
			assert.Equal(t, len(actual.Result), N)

			assert.Equal(t, reasonable.Start, actual.Result[0].Freq)
			assert.Equal(t, reasonable.End, actual.Result[N-1].Freq)

			expectedFreq := LogFrequency(reasonable.Start, reasonable.End, N)

			for i := 0; i < N; i++ {
				if verbose {
					fmt.Printf("%d: %d %d\n", i, int(expectedFreq[i]), int(actual.Result[i].Freq))
				}
				assert.Equal(t, int(expectedFreq[i]), int(actual.Result[i].Freq))
			}

			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	cancel()

}

func TestRunHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	timeout := time.Millisecond * 100

	v := NewHardware()

	command := make(chan interface{})
	result := make(chan interface{})

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	go v.Run(ctx, command, result)

	// Do GetReasonableFrequencyRange command

	reasonable := Range{}

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
			// weak test - with real kit attached, we should get non-zero numbers
			assert.True(t, actual.Result.Start > 0)
			assert.True(t, actual.Result.End > actual.Result.Start)
			reasonable = actual.Result //save for RangeQuery
			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	// Do SingleQuery command

	id = "456abc"
	command <- SingleQuery{
		Command: Command{ID: id},
		Freq:    200000,
		Avg:     1,
		Select:  SParamSelect{true, true, true, true},
	}

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-result:

		if actual, ok := ri.(SingleQuery); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			// weak test - with real kit attached, we should get non-zero numbers
			assert.True(t, actual.Result.S11.Real != 0)
			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	// Do RangeQuery command

	id = "789def"
	N := 7 // number of samples
	command <- RangeQuery{
		Command:         Command{ID: id},
		Range:           reasonable,
		Size:            N,
		Avg:             1,
		LogDistribution: true,
		Select:          SParamSelect{true, true, true, true},
	}

	timeout = time.Second //need more time for this than a single query

	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case ri := <-result:

		if actual, ok := ri.(RangeQuery); !ok {
			t.Error("Wrong type returned")
		} else {

			assert.Equal(t, actual.ID, id)
			// weak test - with real kit attached, we should get non-zero numbers
			assert.Equal(t, len(actual.Result), N)

			assert.Equal(t, reasonable.Start, actual.Result[0].Freq)
			assert.Equal(t, reasonable.End, actual.Result[N-1].Freq)

			expectedFreq := LogFrequency(reasonable.Start, reasonable.End, N)

			for i := 0; i < N; i++ {
				if verbose {
					fmt.Printf("%d: %d %d\n", i, int(expectedFreq[i]), int(actual.Result[i].Freq))
				}
				assert.Equal(t, int(expectedFreq[i]), int(actual.Result[i].Freq))
			}

			if verbose {
				fmt.Println(actual.Result)
			}
		}
	}

	cancel()

}

func TestFrequencyHW(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	var start, end uint64
	start = 1000000
	end = 500000000
	size := 11

	expectedLinear := []uint64{
		1000000,
		50900000,
		100800000,
		150700000,
		200600000,
		250500000,
		300400000,
		350300000,
		400200000,
		450100000,
		500000000,
	}

	expectedLog := []uint64{
		1000000,
		1861646,
		3465724,
		6451950,  //-1 cf native app
		12011244, //-1 cf native app
		22360680,
		41627660,  //-8 cf native app
		77495949,  //+5 cf native app
		144269991, //-9 cf native app
		268579588, //+36 cf native app
		500000000,
	}

	flin := LinFrequency(start, end, size)
	flog := LogFrequency(start, end, size)

	for i := 0; i < size; i++ {
		assert.Equal(t, int(expectedLinear[i]), int(flin[i]))
		assert.Equal(t, int(expectedLog[i]), int(flog[i]))
	}

}
*/
