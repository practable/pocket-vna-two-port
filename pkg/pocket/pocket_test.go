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
	hardware = true

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

func TestLinLogFrequency(t *testing.T) {

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

func TestHardwareGetReleaseHandle(t *testing.T) {
	if !hardware {
		t.Skip("hardware not present")
	}
	handle, err := getFirstDeviceHandle()
	assert.NoError(t, err)
	err = releaseHandle(handle)
	assert.NoError(t, err)
}

func TestHardwareGetReasonableFrequency(t *testing.T) {
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

func TestHardwareSingleQuery(t *testing.T) {
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

func TestHardwareRangeQuery(t *testing.T) {
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
func TestNewHardware(t *testing.T) {

	if !hardware {
		t.Skip("hardware not present")
	}

	_, disconnect, err := NewHardware()

	assert.NoError(t, err)

	defer disconnect()
}

func TestHardwareHandleCommand(t *testing.T) {

	if !hardware {
		t.Skip("hardware not present")
	}

	v, disconnect, err := NewHardware()

	assert.NoError(t, err)

	defer disconnect()

	// ReasonableFrequencyRange

	start := uint64(500000)
	end := uint64(4000000000)

	id0 := "rfr00"
	c0 := ReasonableFrequencyRange{
		Command: Command{
			ID:      id0,
			Command: "rfr",
		},
	}

	err = v.HandleCommand(&c0)
	assert.NoError(t, err)
	assert.Equal(t, id0, c0.ID)
	//note these values are hardware dependent and are taken from our current hardware
	if verbose {
		fmt.Println(c0.Result)
	}
	assert.Equal(t, start, c0.Result.Start)
	assert.Equal(t, end, c0.Result.End)

	// RangeQuery
	id1 := "rq00"
	c1 := RangeQuery{
		Command: Command{
			ID:      id1,
			Command: "rq",
		},
		Range: Range{
			Start: 500000,
			End:   4000000000,
		},
		Size:            7,
		Avg:             1,
		LogDistribution: true,
		Select:          SParamSelect{true, true, true, true},
	}

	err = v.HandleCommand(&c1)
	assert.NoError(t, err)
	assert.Equal(t, id1, c1.ID)

	// weak check that first and last results are not zero
	// there could be anything connected to the VNA
	// so we can't check for specific values here
	// so we just check they're not zero which
	// implies something got put in the array
	assert.NotEqual(t, 0, c1.Result[0].Freq)
	assert.NotEqual(t, 0, c1.Result[0].S11.Imag)
	assert.NotEqual(t, 0, c1.Result[0].S11.Real)
	assert.NotEqual(t, 0, c1.Result[0].S12.Imag)
	assert.NotEqual(t, 0, c1.Result[0].S12.Real)
	assert.NotEqual(t, 0, c1.Result[0].S21.Imag)
	assert.NotEqual(t, 0, c1.Result[0].S21.Real)
	assert.NotEqual(t, 0, c1.Result[0].S22.Imag)
	assert.NotEqual(t, 0, c1.Result[0].S22.Real)
	assert.NotEqual(t, 0, c1.Result[6].Freq)
	assert.NotEqual(t, 0, c1.Result[6].S11.Imag)
	assert.NotEqual(t, 0, c1.Result[6].S11.Real)
	assert.NotEqual(t, 0, c1.Result[6].S12.Imag)
	assert.NotEqual(t, 0, c1.Result[6].S12.Real)
	assert.NotEqual(t, 0, c1.Result[6].S21.Imag)
	assert.NotEqual(t, 0, c1.Result[6].S21.Real)
	assert.NotEqual(t, 0, c1.Result[6].S22.Imag)
	assert.NotEqual(t, 0, c1.Result[6].S22.Real)

	// SingleQuery
	id2 := "sq00"
	c2 := SingleQuery{
		Command: Command{
			ID:      id2,
			Command: "sq",
		},
		Freq:   1000000000,
		Avg:    1,
		Select: SParamSelect{true, true, true, true},
	}

	err = v.HandleCommand(&c2)
	assert.NoError(t, err)
	assert.Equal(t, id2, c2.ID)

	// weak check that first and last results are not zero
	// there could be anything connected to the VNA
	// so we can't check for specific values here
	// so we just check they're not zero which
	// implies something got put in the array
	assert.NotEqual(t, 0, c2.Result.Freq)
	assert.NotEqual(t, 0, c2.Result.S11.Imag)
	assert.NotEqual(t, 0, c2.Result.S11.Real)
	assert.NotEqual(t, 0, c2.Result.S12.Imag)
	assert.NotEqual(t, 0, c2.Result.S12.Real)
	assert.NotEqual(t, 0, c2.Result.S21.Imag)
	assert.NotEqual(t, 0, c2.Result.S21.Real)
	assert.NotEqual(t, 0, c2.Result.S22.Imag)
	assert.NotEqual(t, 0, c2.Result.S22.Real)
}
