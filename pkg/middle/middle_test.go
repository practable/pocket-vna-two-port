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
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/go-pocketvna/pkg/drain"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
	"github.com/timdrysdale/go-pocketvna/pkg/stream"
)

var verbose bool
var debug bool

func TestMain(m *testing.M) {
	// Setup  logging
	debug = false
	verbose = false

	if debug {
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
	return stream.New(u, ctx)
}

// This test demonstrates draining the fromClient channel
// if you don't, the handler throws a test error

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

		assert.Equal(t, "rr", v.Command.Command)
	}

	// Send something back to avoid mock Handler stalling on readmessage
	select {
	case stream.Response <- pocket.Command{ID: "0"}:
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

	time.Sleep(100 * time.Millisecond)
}

// this test confirms the stream package is working ok (TODO remove once tests developed)
func TestStream(t *testing.T) {

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

	//stream := stream.New(u, ctx) //This works
	stream := fakeMiddle(u, ctx) //this works

	mt := int(websocket.TextMessage)

	/* Test ReasonableFrequencyRange */

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

		assert.Equal(t, "rr", v.Command.Command)
	}

	// Send something back to avoid mock Handler stalling on readmessage
	select {
	case stream.Response <- pocket.Command{ID: "0"}:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send response")
	}

	// outgoing pipe does not depend on type...
	// note there is a 1ms timeout in the handler which could be
	// fragile to slow-running test environments
	// you will not get messages that were not available to be taken
	// from the channel within the timelimit (it does not block on that write)
	// as it is only a testing feature, it does not affect actual operation
	// there is no message loss in the actual code - just the testing mock
	// since not all messages need to be sent/received in the tests
	msg := <-fromClient
	expected := "{\"id\":\"0\",\"t\":0,\"cmd\":\"\"}"
	assert.Equal(t, expected, string(msg.Data))

	/* Test rangeQuery */
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
	case toClient <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case request := <-stream.Request:

		v, ok := request.(pocket.RangeQuery)

		assert.True(t, ok)

		assert.Equal(t, "rq", v.Command.Command)

	}
	// Send something back to avoid mock Handler stalling on readmessage
	select {
	case stream.Response <- pocket.Command{ID: "1"}:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send response")
	}

	// outgoing pipe does not depend on type...
	msg = <-fromClient
	expected = "{\"id\":\"1\",\"t\":0,\"cmd\":\"\"}"
	assert.Equal(t, expected, string(msg.Data))

	/* Test calibratedRangeQuery */
	crq := pocket.CalibratedRangeQuery{
		Command: pocket.Command{Command: "crq"},
		Avg:     1,
		Select:  pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
	}

	message, err = json.Marshal(crq)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
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

		v, ok := request.(pocket.CalibratedRangeQuery)

		assert.True(t, ok)

		assert.Equal(t, "crq", v.Command.Command)

	}

	// Send something back to avoid mock Handler stalling on readmessage
	select {
	case stream.Response <- pocket.Command{ID: "2"}:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send response")
	}

	// outgoing pipe does not depend on type...
	msg = <-fromClient
	expected = "{\"id\":\"2\",\"t\":0,\"cmd\":\"\"}"
	assert.Equal(t, expected, string(msg.Data))

}

func TestMiddle(t *testing.T) {

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
	uc := "ws://localhost:8888/ws/calibration"
	ur := "ws://localhost:8888/ws/rfswitch"
	us := "ws" + strings.TrimPrefix(ss.URL, "http")

	New(uc, ur, us, ctx)

	/* Test ReasonableFrequencyRange */

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
	}

	/* Test rangeQuery */

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

		assert.NoError(t, err)

		assert.Equal(t, "rq", rq.Command.Command)

		// TODO check message contents are ok

		// cast to int to make human readable in assert error message
		assert.Equal(t, 100000, int(rq.Result[0].Freq))
		assert.Equal(t, 4000000, int(rq.Result[1].Freq))

	}

	/* Test rangeCal */

	// Should throw an error because S21 is also true, but we only support 1port cal...

	rq = pocket.RangeQuery{
		Command:         pocket.Command{Command: "rc", ID: "bad"},
		Range:           pocket.Range{Start: 100000, End: 4000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
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

	// manage variable scope
	var responseFiltered reconws.WsMessage

FILTERBAD:

	for i := 0; i < 5; i++ {
		if debug {
			fmt.Printf("FILTERBAD: iteration %d\n", i)
		}
		if i > 4 {
			t.Error("timeout awaiting response")
		}
		select {

		case <-time.After(timeout):
			//silently wait - could be a slow scan
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)

			if debug {
				idx, err := mds.LastReadIndex()
				assert.NoError(t, err)
				fmt.Printf("BAD-UNFILTERED: %d->%s:\n", idx, m.Data)
			}

			if string(m.Data) != "{\"cmd\":\"hb\"}" {
				responseFiltered = m
				if debug {
					fmt.Printf("BAD-FILTERED: %s\n", responseFiltered.Data)
				}
				break FILTERBAD
			}
		}
	}

	var cr pocket.CustomResult

	err = json.Unmarshal(responseFiltered.Data, &cr)

	assert.NoError(t, err)

	expectedError := "Error: calibration is only supported on Port1 (S11). Resend the command with only S11 selected (true). You had S11:true, S12:false, S21:true, S22:false"

	assert.Equal(t, expectedError, cr.Message)

	// re unmarshal to get the command info

	err = json.Unmarshal(responseFiltered.Data, &rq)

	assert.NoError(t, err)

	assert.Equal(t, "rc", rq.Command.Command)
	assert.Equal(t, "bad", rq.Command.ID)

	// Check the command was sent back to us so we can check it
	assert.Equal(t, 100000, int(rq.Range.Start))
	assert.Equal(t, 4000000, int(rq.Range.End))

	// Check the results are empty (as they should be on an error)
	assert.Equal(t, 0, len(rq.Result))

	/* Test calibratedRangeQuery - will fail because there is no cal in place*/

	crq := pocket.CalibratedRangeQuery{
		Command: pocket.Command{Command: "crq"},
		What:    "dut",
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

	responseFiltered = reconws.WsMessage{}

FILTERCRQNOCAL:
	for i := 0; i < 20; i++ {
		if debug {
			fmt.Printf("FILTERCRQGOOD: iteration %d\n", i)
		}
		if i > 19 {
			t.Error("timeout awaiting response")
		}
		select {

		case <-time.After(timeout):
			//silently wait - could be a slow scan
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)
			if debug {
				idx, err := mds.LastReadIndex()
				assert.NoError(t, err)
				fmt.Printf("CRQGOOD-UNFILTERED: %d->%s:\n", idx, m.Data)
			}
			if string(m.Data) != "{\"cmd\":\"hb\"}" &&
				string(m.Data) != "" {

				responseFiltered = m
				if debug {
					fmt.Printf("CRGOOD-FILTERED: %s\n", responseFiltered.Data)
				}
				break FILTERCRQNOCAL
			}
		}
	}

	if debug {
		fmt.Printf("RECV:" + string(responseFiltered.Data) + "\n")
	}

	// should just be a pocket.RangeQuery result when it is a success
	// unmarshal to get the command info

	err = json.Unmarshal(responseFiltered.Data, &cr)

	assert.NoError(t, err)

	expectedError = "Error. No existing calibration. Please calibrate with rc command"

	assert.Equal(t, expectedError, cr.Message)

	/* Test rangeCal with correct S11 setting */
	rq = pocket.RangeQuery{
		Command:         pocket.Command{Command: "rc", ID: "good"},
		Range:           pocket.Range{Start: 200000, End: 5000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: false, S22: false},
	}

	message, err = json.Marshal(rq)

	assert.NoError(t, err)

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	if debug {
		fmt.Printf("SENT:" + string(ws.Data) + "\n")
	}

	select {
	case streamWrite <- ws:
	case <-time.After(timeout):
		t.Error(t, "timeout awaiting send message")
	}

	responseFiltered = reconws.WsMessage{}

FILTERGOOD:
	for i := 0; i < 20; i++ {
		if debug {
			fmt.Printf("FILTERGOOD: iteration %d\n", i)
		}
		if i > 19 {
			t.Error("timeout awaiting response")
		}
		select {

		case <-time.After(timeout):
			//silently wait - could be a slow scan
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)
			if debug {
				idx, err := mds.LastReadIndex()
				assert.NoError(t, err)
				fmt.Printf("GOOD-UNFILTERED: %d->%s:\n", idx, m.Data)
			}
			if string(m.Data) != "{\"cmd\":\"hb\"}" &&
				string(m.Data) != "" {

				responseFiltered = m
				if debug {
					fmt.Printf("GOOD-FILTERED: %s\n", responseFiltered.Data)
				}
				break FILTERGOOD
			}
		}
	}

	if debug {
		fmt.Printf("RECV:" + string(responseFiltered.Data) + "\n")
	}

	// should just be a pocket.RangeQuery result when it is a success
	// unmarshal to get the command info

	err = json.Unmarshal(responseFiltered.Data, &rq)

	assert.NoError(t, err)

	if debug {
		fmt.Printf("CHECK RECV AGAIN:" + string(responseFiltered.Data) + "\n")
		fmt.Printf("RQ:%+v\n", rq)
	}

	assert.Equal(t, "rc", rq.Command.Command)
	assert.Equal(t, "good", rq.Command.ID)

	if debug {
		fmt.Printf("MARSHALLED-rq: %+v\n", rq)
	}

	// check we got some results back
	assert.Equal(t, 2, len(rq.Result))

	//avoid panic if results are unexpectedly empty, and still fail the test
	if len(rq.Result) == 2 {
		assert.Equal(t, 200000, int(rq.Result[0].Freq))
		assert.Equal(t, 5000000, int(rq.Result[1].Freq))
	} else {
		t.Error("Wrong length array, could not check values")
	}

	if debug {
		for i, msg := range mds.All() {

			fmt.Printf("%d: %s\n", i, ((msg.(reconws.WsMessage)).Data))
		}
	}

	/* Test calibratedRangeQuery */
	crq = pocket.CalibratedRangeQuery{
		Command: pocket.Command{Command: "crq"},
		What:    "dut",
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

	responseFiltered = reconws.WsMessage{}

FILTERCRQGOOD:
	for i := 0; i < 20; i++ {
		if debug {
			fmt.Printf("FILTERCRQGOOD: iteration %d\n", i)
		}
		if i > 19 {
			t.Error("timeout awaiting response")
		}
		select {

		case <-time.After(timeout):
			//silently wait - could be a slow scan
		case response := <-mds.Next():

			m, ok := response.(reconws.WsMessage)
			assert.True(t, ok)
			if debug {
				idx, err := mds.LastReadIndex()
				assert.NoError(t, err)
				fmt.Printf("CRQGOOD-UNFILTERED: %d->%s:\n", idx, m.Data)
			}
			if string(m.Data) != "{\"cmd\":\"hb\"}" &&
				string(m.Data) != "" {

				responseFiltered = m
				if debug {
					fmt.Printf("CRGOOD-FILTERED: %s\n", responseFiltered.Data)
				}
				break FILTERCRQGOOD
			}
		}
	}

	if debug {
		fmt.Printf("RECV:" + string(responseFiltered.Data) + "\n")
	}

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
					fmt.Printf("userChannelHandler: writing message %s\n", string(msg.Data))
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
