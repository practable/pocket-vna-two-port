package rfusb

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var hardware bool

var port string
var baud int
var timeout time.Duration

func TestMain(m *testing.M) {
	// Setup  logging
	debug := false
	hardware = true

	if debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&log.JSONFormatter{})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	port = "/dev/ttyUSB0"
	baud = 57600
	timeout = time.Duration(time.Second)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestOpen(t *testing.T) {
	if !hardware {
		t.Skip("no hardware")
	}

	r := NewRFUSB()

	err := r.Open(port, baud, timeout)

	assert.NoError(t, err)

	err = r.Close()

	assert.NoError(t, err)

}

func TestSetPort(t *testing.T) {
	if !hardware {
		t.Skip("no hardware")
	}

	r := NewRFUSB()

	err := r.Open(port, baud, timeout)

	assert.NoError(t, err)

	time.Sleep(time.Second)

	err = r.SetPort("short")

	assert.NoError(t, err)

	err = r.Close()

	assert.NoError(t, err)

}

func TestSettingPorts(t *testing.T) {
	if !hardware {
		t.Skip("no hardware")
	}

	r := NewRFUSB()

	err := r.Open(port, baud, timeout)
	assert.NoError(t, err)

	time.Sleep(time.Second)

	err = r.SetShort()
	assert.NoError(t, err)
	assert.Equal(t, "short", r.Get())

	err = r.SetOpen()
	assert.NoError(t, err)
	assert.Equal(t, "open", r.Get())

	err = r.SetLoad()
	assert.NoError(t, err)
	assert.Equal(t, "load", r.Get())

	err = r.SetThru()
	assert.NoError(t, err)
	assert.Equal(t, "thru", r.Get())

	err = r.SetDUT1()
	assert.NoError(t, err)
	assert.Equal(t, "dut1", r.Get())

	err = r.SetDUT2()
	assert.NoError(t, err)
	assert.Equal(t, "dut2", r.Get())

	err = r.SetDUT3()
	assert.NoError(t, err)
	assert.Equal(t, "dut3", r.Get())

	err = r.SetDUT4()
	assert.NoError(t, err)
	assert.Equal(t, "dut4", r.Get())

	err = r.SetDUT3()
	assert.NoError(t, err)
	assert.Equal(t, "dut3", r.Get())

	err = r.SetDUT2()
	assert.NoError(t, err)
	assert.Equal(t, "dut2", r.Get())

	err = r.SetDUT1()
	assert.NoError(t, err)
	assert.Equal(t, "dut1", r.Get())

	err = r.SetThru()
	assert.NoError(t, err)
	assert.Equal(t, "thru", r.Get())

	err = r.SetLoad()
	assert.NoError(t, err)
	assert.Equal(t, "load", r.Get())

	err = r.SetOpen()
	assert.NoError(t, err)
	assert.Equal(t, "open", r.Get())

	err = r.SetShort()
	assert.NoError(t, err)
	assert.Equal(t, "short", r.Get())

	err = r.Close()

	assert.NoError(t, err)

}

func TestSettingPortsMock(t *testing.T) {

	r := NewMock()

	err := r.Open(port, baud, timeout)
	assert.NoError(t, err)

	err = r.SetShort()
	assert.NoError(t, err)
	assert.Equal(t, "short", r.Get())

	err = r.SetOpen()
	assert.NoError(t, err)
	assert.Equal(t, "open", r.Get())

	err = r.SetLoad()
	assert.NoError(t, err)
	assert.Equal(t, "load", r.Get())

	err = r.SetThru()
	assert.NoError(t, err)
	assert.Equal(t, "thru", r.Get())

	err = r.SetDUT1()
	assert.NoError(t, err)
	assert.Equal(t, "dut1", r.Get())

	err = r.SetDUT2()
	assert.NoError(t, err)
	assert.Equal(t, "dut2", r.Get())

	err = r.SetDUT3()
	assert.NoError(t, err)
	assert.Equal(t, "dut3", r.Get())

	err = r.SetDUT4()
	assert.NoError(t, err)
	assert.Equal(t, "dut4", r.Get())

	err = r.SetDUT3()
	assert.NoError(t, err)
	assert.Equal(t, "dut3", r.Get())

	err = r.SetDUT2()
	assert.NoError(t, err)
	assert.Equal(t, "dut2", r.Get())

	err = r.SetDUT1()
	assert.NoError(t, err)
	assert.Equal(t, "dut1", r.Get())

	err = r.SetThru()
	assert.NoError(t, err)
	assert.Equal(t, "thru", r.Get())

	err = r.SetLoad()
	assert.NoError(t, err)
	assert.Equal(t, "load", r.Get())

	err = r.SetOpen()
	assert.NoError(t, err)
	assert.Equal(t, "open", r.Get())

	err = r.SetShort()
	assert.NoError(t, err)
	assert.Equal(t, "short", r.Get())

	err = r.Close()

	assert.NoError(t, err)

}