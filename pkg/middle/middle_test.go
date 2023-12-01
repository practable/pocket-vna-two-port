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

	timeout := time.Millisecond * 1000

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

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")

	case response := <-mds.Next():

		//fmt.Printf(string((request.(reconws.WsMessage)).Data))

		m, ok := response.(reconws.WsMessage)

		assert.True(t, ok)

		var rr pocket.ReasonableFrequencyRange

		err := json.Unmarshal(m.Data, &rr)

		assert.NoError(t, err)

		assert.True(t, ok)

		assert.Equal(t, "rr", rr.Command.Command)

		// These are hardware dependent tests and can be changed to suit the hardware you are testing with
		assert.Equal(t, uint64(500000), rr.Result.Start)
		assert.Equal(t, uint64(4000000000), rr.Result.End)
	}

} /*

	// Test rangeQuery

	rq := pocket.RangeQuery{
		Command:         pocket.Command{Command: "rq"},
		Range:           pocket.Range{Start: 100000, End: 4000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
	}

	message, err := json.Marshal(rq)

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

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case response := <-mds.Next():

		m, ok := response.(reconws.WsMessage)

		assert.True(t, ok)

		var rq pocket.RangeQuery

		err := json.Unmarshal(m.Data, &rq)

		log.Debugf("rq result: %s", m.Data)

		assert.NoError(t, err)
		assert.Equal(t, "rq", rq.Command.Command)

		// TODO check message contents are ok

		// cast to int to make human readable in assert error message
		if len(rq.Result) == 2 {
			assert.Equal(t, 100000, int(rq.Result[0].Freq))
			assert.Equal(t, 4000000, int(rq.Result[1].Freq))
		} else {
			t.Fatal("wrong length results")
		}

	}

	// Test rangeCal

	// TODO - consider replacment test. Now won't fail on incorrect Select though
	// so do test later to avoid messing up next test ...

    // Test calibratedRangeQuery - will fail because there is no cal in place

	crq := pocket.CalibratedRangeQuery{
		Command: pocket.Command{Command: "crq", ID: "crq-before-cal"},
		What:    "dut1",
		Avg:     1,
		Select:  pocket.SParamSelect{S11: true, S12: false, S21: false, S22: false},
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

	idx, err := mds.LastReadIndex()
	assert.NoError(t, err)
	<-time.After(time.Second)

	msgs := mds.All()

	var responseFiltered reconws.WsMessage

LOOP2:

	for i := idx; i < (len(msgs) - 1); i++ {

		response, err := mds.NextNoWait()
		assert.NoError(t, err)
		lri, err := mds.LastReadIndex()
		assert.NoError(t, err)
		log.Infof("Read %d should match LastRead %d\n", i+1, lri)
		m, ok := response.(reconws.WsMessage)
		assert.True(t, ok)
		if string(m.Data) != "{\"cmd\":\"hb\"}" {
			responseFiltered = m
			break LOOP2
		}
	}

	// should just be a pocket.RangeQuery result when it is a success
	// unmarshal to get the command info

	var cr pocket.CustomResult

	err = json.Unmarshal(responseFiltered.Data, &cr)

	assert.NoError(t, err)

	expectedError := "Error. No existing calibration. Please calibrate with rc command"

	assert.Equal(t, expectedError, cr.Message)

	// Test rangeCal

	rq = pocket.RangeQuery{
		Command:         pocket.Command{Command: "rc", ID: "good"},
		Range:           pocket.Range{Start: 200000, End: 5000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		// no need for Select because cal routine handles it
	}

	message, err = json.Marshal(rq)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	log.Debugf("SENT:" + string(ws.Data) + "\n")

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	idx, err = mds.LastReadIndex()
	assert.NoError(t, err)

	<-time.After(5 * time.Second)

	msgs = mds.All()

	responseFiltered = reconws.WsMessage{}

	if !mds.IsEmpty() {

		log.Infof("i:%d,len(msgs):%d\n", idx+1, len(msgs))
	LOOP3:
		for i := idx; i < (len(msgs) - 1); i++ {

			response, err := mds.NextNoWait()
			assert.NoError(t, err)
			lri, err := mds.LastReadIndex()
			assert.NoError(t, err)
			log.Infof("Read %d should match LastRead %d\n", i+1, lri)
			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)
			if string(m.Data) != "{\"cmd\":\"hb\"}" { //hb is heartbeat
				responseFiltered = m
				break LOOP3
			}
		}

		// check for RF switch error

		cr = pocket.CustomResult{}

		err = json.Unmarshal(responseFiltered.Data, &cr)

		assert.NoError(t, err)

		if err == nil && cr.Message != "" {
			if strings.Contains(cr.Message, "Error setting RF switch") {
				t.Errorf("Hardware issue - repeat test! %s", cr.Message)
			} else {
				t.Error(cr.Message)
			}
		} else {

			// should just be a pocket.RangeQuery result when it is a success
			// unmarshal to get the command info

			err = json.Unmarshal(responseFiltered.Data, &rq)

			assert.NoError(t, err)

			assert.Equal(t, "rc", rq.Command.Command)
			assert.Equal(t, "good", rq.Command.ID)

			// check we got some results back
			assert.Equal(t, 2, len(rq.Result))

			//avoid panic if results are unexpectedly empty, and still fail the test
			if len(rq.Result) == 2 {
				assert.Equal(t, 200000, int(rq.Result[0].Freq))
				assert.Equal(t, 5000000, int(rq.Result[1].Freq))
			} else {
				t.Error("Wrong length array, could not check values")
			}
		}
	} else {
		t.Error("No new messages for Test3")
	}

	// Test calibratedRangeQuery
	crq = pocket.CalibratedRangeQuery{
		Command: pocket.Command{Command: "crq"},
		What:    "dut1",
		Avg:     1,
		Select:  pocket.SParamSelect{S11: true, S12: false, S21: false, S22: false},
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

	idx, err = mds.LastReadIndex()
	assert.NoError(t, err)

	<-time.After(2 * time.Second)

	msgs = mds.All()

	responseFiltered = reconws.WsMessage{}

	if !mds.IsEmpty() {
		log.Infof("i:%d,len(msgs):%d\n", idx+1, len(msgs))
	LOOP4:
		for i := idx; i < (len(msgs) - 1); i++ {

			response, err := mds.NextNoWait()
			assert.NoError(t, err)
			lri, err := mds.LastReadIndex()
			assert.NoError(t, err)
			log.Infof("Read %d should match LastRead %d\n", i+1, lri)
			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)
			if string(m.Data) != "{\"cmd\":\"hb\"}" {
				responseFiltered = m
				break LOOP4
			}
		}

		cr = pocket.CustomResult{}
		err = json.Unmarshal(responseFiltered.Data, &cr)

		assert.NoError(t, err)

		if err == nil && cr.Message != "" {
			if strings.Contains(cr.Message, "Error setting RF switch") {
				t.Errorf("Hardware issue - repeat test! %s", cr.Message)
			} else {
				t.Error(cr.Message)
			}
		} else {

			// should just be a pocket.RangeQuery result when it is a success
			// unmarshal to get the command info

			crq = pocket.CalibratedRangeQuery{}

			err = json.Unmarshal(responseFiltered.Data, &crq)

			assert.NoError(t, err)

			assert.Equal(t, "crq", crq.Command.Command)

			//avoid panic if results are unexpectedly empty, and still fail the test
			if len(rq.Result) == 2 {
				assert.Equal(t, 200000, int(rq.Result[0].Freq))
				assert.Equal(t, 5000000, int(rq.Result[1].Freq))
			} else {
				t.Error("Wrong length array, could not check values")
			}
		}
	} else {
		t.Error("No messages received for Test4")
	}

	if verbose {
		for i, msg := range mds.All() {

			log.Infof("%d: %s\n", i, ((msg.(reconws.WsMessage)).Data))
		}
	}

}
*/
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
