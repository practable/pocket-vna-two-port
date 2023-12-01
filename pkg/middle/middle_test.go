package middle

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/practable/pocket-vna-two-port/pkg/drain"
	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	"github.com/practable/pocket-vna-two-port/pkg/reconws"
	"github.com/practable/pocket-vna-two-port/pkg/stream"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var verbose bool
var debug bool

/*
to setup for the test:
connect rfswitch
sessionrelay host&
cd pkg/rfswitch
./connectlocalswitch.sh /dev/ttyUSB1 #runs in background (check correct /dev/? for RF switch with dmesg)
cd ../../py
python client.py & #run the calibration service
*/

func TestMain(m *testing.M) {
	// Setup  logging
	debug = true
	verbose = false

	if debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else if !debug && verbose {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)
	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	err := pocket.ForceUnlockDevices()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

var upgrader = websocket.Upgrader{}

func fakeMiddle(u string, ctx context.Context) stream.Stream {
	return stream.New(ctx, u)
}

// This test demonstrates draining the fromClient channel
// if you don't, the handler throws a test error
// TODO check if this is still relevant after upgrading with pocket/measure/rfusb?
func TestFakeMiddle(t *testing.T) {

	timeout := 100 * time.Millisecond

	toClient := make(chan reconws.WsMessage)
	fromClient := make(chan reconws.WsMessage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create test server with the channel handler.
	s := httptest.NewServer(http.HandlerFunc(userChannelHandler(t, toClient, fromClient, ctx)))

	//s := httptest.NewServer(http.HandlerFunc(reasonableRange))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	stream := fakeMiddle(u, ctx) //this works

	mt := int(websocket.TextMessage)

	message := []byte("{\"cmd\":\"rr\"}")

	ws := reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case toClient <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case request := <-stream.Request:

		v, ok := request.(pocket.ReasonableFrequencyRange)

		assert.True(t, ok)
		if debug {
			t.Log(v)
		}
		assert.Equal(t, "rr", v.Command.Command)
	}

	// Send something back to avoid mock Handler stalling on readmessage
	select {
	case stream.Response <- pocket.Command{ID: "1"}:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send response")
	}

	// must drain the fromClient channel to avoid the test Handler closing early.
	// Check the behaviour by setting this false, and not drain channel
	if true {
		go func() {
			for {

				select {

				case <-ctx.Done():
					return
				case <-fromClient:
					// carry on
				}

			}
		}()
	}

	time.Sleep(1000 * time.Millisecond)
}

func TestMiddle(t *testing.T) {
	if debug {
		t.Log("This test requires external dependencies and succeeds most of the time")
		t.Log("Hardware required: VNA, RFSwitch")
		t.Log("Services required:, gRPC calibration service py/server.py")
	}

	timeout := 2 * time.Second

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	streamWrite := make(chan reconws.WsMessage)
	streamRead := make(chan reconws.WsMessage)

	mds := drain.NewWs(streamRead, ctx)

	// Create test server with the channel handler.
	ss := httptest.NewServer(http.HandlerFunc(userChannelHandler(t, streamWrite, streamRead, ctx)))
	defer ss.Close()

	// URL only assigned after starting
	// Convert http://127.0.0.1 to ws://127.0.0.
	topic := "ws" + strings.TrimPrefix(ss.URL, "http")

	addr := "localhost:9001" //gRPC calibration service
	port := "/dev/ttyUSB0"
	baud := 57600
	timeoutUSB := time.Duration(time.Minute)
	timeoutRequest := time.Duration(time.Minute) //2min in production for large calibrated scans?

	v, disconnect, err := pocket.NewHardware()

	defer disconnect()

	assert.NoError(t, err)

	m := New(ctx, addr, port, baud, timeoutUSB, timeoutRequest, topic, &v)

	go m.Run()

	// Test ReasonableFrequencyRange

	mt := int(websocket.TextMessage)
	message := []byte("{\"cmd\":\"rr\"}")

	ws := reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	t0 := time.Now()
RFR:
	for {
		select {

		case <-time.After(timeout): //this times out if there are no heartbeats
			t.Error("RFR: timeout awaiting response")
			break RFR

		case response := <-mds.Next():

			//fmt.Printf(string((request.(reconws.WsMessage)).Data))

			m, ok := response.(reconws.WsMessage)

			assert.True(t, ok)

			var rr pocket.ReasonableFrequencyRange

			err := json.Unmarshal(m.Data, &rr)

			assert.NoError(t, err)

			assert.True(t, ok)

			///ignore heartbeat, but timeout if only get heartbeats for too long
			if rr.Command.Command == "hb" {
				t.Log("RFR: heartbeat")
				if time.Now().After(t0.Add(timeout)) {
					t.Fatal("RFR timeout")
					break RFR
				}
				continue
			}

			assert.Equal(t, "rr", rr.Command.Command)

			// These are hardware dependent tests and can be changed to suit the hardware you are testing with
			assert.Equal(t, uint64(500000), rr.Result.Start)
			assert.Equal(t, uint64(4000000000), rr.Result.End)
			t.Log("RFR test completed")
			break RFR
		}
	}

	// Test rangeQuery

	rq := pocket.RangeQuery{
		Command:         pocket.Command{Command: "rq"},
		Range:           pocket.Range{Start: 100000, End: 4000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
		What:            "short",
	}

	message, err = json.Marshal(rq)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	t0 = time.Now()

RQ:
	for {
		select {

		case <-time.After(timeout): // timeout if no heartbeats or responses
			t.Error("timeout awaiting response")
			break RQ
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)

			assert.True(t, ok)

			var rq pocket.RangeQuery

			err := json.Unmarshal(m.Data, &rq)

			log.Debugf("rq result: %s", m.Data)

			assert.NoError(t, err)

			//ignore heartbeats but timeout if only get heartbeats for too long
			if rq.Command.Command == "hb" {
				if time.Now().After(t0.Add(timeout)) {
					t.Fatal("RQ timeout")
					break RQ
				}
				t.Log("RQ: heartbeat")
				continue
			}

			assert.Equal(t, "rq", rq.Command.Command)

			// TODO check message contents are ok

			// cast to int to make human readable in assert error message
			if len(rq.Result) == 2 {
				assert.Equal(t, 100000, int(rq.Result[0].Freq))
				assert.Equal(t, 4000000, int(rq.Result[1].Freq))
			} else {
				t.Fatal("RQ wrong length results")
			}
			t.Log("RQ test completed")
			break RQ

		}
	}

	// Test rangeCal

	timeout = 8 * time.Second //needs extra time compared to a single measurement

	rc := pocket.RangeQuery{
		Command:         pocket.Command{Command: "rc"},
		Range:           pocket.Range{Start: 100000, End: 4000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
		What:            "short",
	}

	message, err = json.Marshal(rc)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	t0 = time.Now()

RC:
	for {
		select {

		case <-time.After(timeout): // timeout if no heartbeats or responses
			t.Error("timeout awaiting response")
			break RC
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)

			assert.True(t, ok)

			var rq pocket.RangeQuery

			err := json.Unmarshal(m.Data, &rq)

			log.Debugf("rq result: %s", m.Data)

			assert.NoError(t, err)

			//ignore heartbeats but timeout if only get heartbeats for too long
			if rq.Command.Command == "hb" {
				if time.Now().After(t0.Add(timeout)) {
					t.Fatal("RC timeout")
					break RC
				}
				t.Log("RC: heartbeat")
				continue
			}

			assert.Equal(t, "rc", rq.Command.Command)

			// TODO check message contents are ok

			// cast to int to make human readable in assert error message
			if len(rq.Result) == 2 {
				assert.Equal(t, 100000, int(rq.Result[0].Freq))
				assert.Equal(t, 4000000, int(rq.Result[1].Freq))
			} else {
				t.Fatal("RC wrong length results")
			}
			t.Log("RC test completed")
			break RC

		}
	}

	// Test calibratedRangeQuery

	crq := pocket.RangeQuery{
		Command: pocket.Command{Command: "crq"},
		Avg:     1,
		Size:    2,
		Select:  pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
		What:    "dut2",
	}

	message, err = json.Marshal(crq)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	t0 = time.Now()

CRQ:
	for {
		select {

		case <-time.After(timeout): // timeout if no heartbeats or responses
			t.Error("timeout awaiting response")
			break CRQ
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)

			assert.True(t, ok)

			var crq pocket.CalibratedRangeQuery

			err := json.Unmarshal(m.Data, &crq)

			log.Debugf("rq result: %s", m.Data)

			assert.NoError(t, err)

			//ignore heartbeats but timeout if only get heartbeats for too long
			if crq.Command.Command == "hb" {
				if time.Now().After(t0.Add(timeout)) {
					t.Fatal("CRQ timeout")
					break CRQ
				}
				t.Log("CRQ: heartbeat")
				continue
			}

			assert.Equal(t, "crq", crq.Command.Command)

			// TODO check message contents are ok

			// cast to int to make human readable in assert error message
			if len(crq.Result) == 2 {
				assert.Equal(t, 100000, int(crq.Result[0].Freq))
				assert.Equal(t, 4000000, int(crq.Result[1].Freq))
			} else {
				t.Fatal("CRQ wrong length results")
			}
			t.Log("CRQ test completed")
			break CRQ

		}
	}

}

func userChannelHandler(t *testing.T, toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// needs more than 1ms when there are multiple servers to set up, else miss first message - never writes.
		timeout := 10 * time.Millisecond

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Cannot upgrade")
			return
		}
		defer c.Close()

		for { //backwards to normal server: we send a message to the client, get a response, repeat until cancel
			// write our first message to the websocket client....
			select {

			case <-ctx.Done():
				return

			case <-time.After(timeout):
				//carry on

			case msg := <-toClient:

				err = c.WriteMessage(msg.Type, msg.Data)
				if err != nil {
					break
				}
				if verbose {
					log.Infof("userChannelHandler: writing message %s\n", string(msg.Data))
				}

			} //select

			// read from the Client's websocket connection
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			// timeout if we don't manage to write to the fromClient channel
			select {
			case <-time.After(timeout):
				// throw an error to avoid silently dropping messages (harder to debug)
				t.Errorf("userChannelHandler: is your test correctly draining the fromClient channel? Timed out trying to write message: %s\n", message)
			case fromClient <- reconws.WsMessage{Data: message, Type: mt}:
			}

		} // for

	} //anon func

}

func serviceChannelHandler(t *testing.T, toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		timeout := 1 * time.Millisecond

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Cannot upgrade")
			return
		}
		defer c.Close()

		for { //normal order - receive message over websocket, and reply

			// read from the Client's websocket connection
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			// timeout if we don't manage to write to the fromClient channel
			select {
			case <-time.After(timeout):
				// throw an error to avoid silently dropping messages (harder to debug)
				t.Errorf("serviceChannelHandler: is your test correctly draining the fromClient channel? Timed out trying to write message: %s\n", message)
			case fromClient <- reconws.WsMessage{Data: message, Type: mt}:
			}

			select {

			case <-ctx.Done():
				return

			case <-time.After(timeout):
				//carry on

			case msg := <-toClient:

				err = c.WriteMessage(msg.Type, msg.Data)
				if err != nil {
					break
				}

			} //select

		} // for

	} //anon func

}
