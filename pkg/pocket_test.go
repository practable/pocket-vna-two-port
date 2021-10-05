package pocket

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

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

	handle, err := GetFirstDeviceHandle()
	assert.NoError(t, err)
	err = ReleaseHandle(handle)
	assert.NoError(t, err)
}

func TestPrintDeviceDetails(t *testing.T) {
	handle, err := GetFirstDeviceHandle()
	assert.NoError(t, err)

	err = ReleaseHandle(handle)
	assert.NoError(t, err)
	/*

			typedef struct PocketVnaDeviceDesc {
		    const char * path;
		    PVNA_Access access;

		    const wchar_t * serial_number;

		    const wchar_t * manufacturer_string;
		    const wchar_t * product_string;

		    uint16_t release_number;

		    uint16_t pid;
		    uint16_t vid;
		    uint16_t ciface_code; //value from ConnectionInterfaceCode

		    struct PocketVnaDeviceDesc * next;
	*/
}
