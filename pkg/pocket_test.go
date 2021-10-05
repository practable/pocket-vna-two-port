package pocket

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	/*
		result = C.pocketvna_get_first_device_handle(&handle)
			fmt.Println(Decode(result))
			result = C.pocketvna_release_handle(&handle)
			fmt.Println(Decode(result))
	*/
}
