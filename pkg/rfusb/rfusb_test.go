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
	debug := true
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

	r := New()

	err := r.Open(port, baud, timeout)

	assert.NoError(t, err)

	err = r.Close()

	assert.NoError(t, err)

}

func TestSetPort(t *testing.T) {
	if !hardware {
		t.Skip("no hardware")
	}

	r := New()

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

	r := New()

	err := r.Open(port, baud, timeout)
	assert.NoError(t, err)

	time.Sleep(time.Second)

	err = r.SetShort()
	assert.NoError(t, err)

	err = r.SetOpen()
	assert.NoError(t, err)

	err = r.SetLoad()
	assert.NoError(t, err)

	err = r.SetThru()
	assert.NoError(t, err)

	err = r.SetDUT1()
	assert.NoError(t, err)

	err = r.SetDUT2()
	assert.NoError(t, err)

	err = r.SetDUT3()
	assert.NoError(t, err)

	err = r.SetDUT4()
	assert.NoError(t, err)

	err = r.SetDUT3()
	assert.NoError(t, err)

	err = r.SetDUT2()
	assert.NoError(t, err)

	err = r.SetDUT1()
	assert.NoError(t, err)

	err = r.SetThru()
	assert.NoError(t, err)

	err = r.SetLoad()
	assert.NoError(t, err)

	err = r.SetOpen()
	assert.NoError(t, err)

	err = r.SetShort()
	assert.NoError(t, err)

	err = r.Close()

	assert.NoError(t, err)

}
