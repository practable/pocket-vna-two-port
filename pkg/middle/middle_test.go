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

func TestMain(m *testing.M) {
	// Setup  logging
	debug := false
	verbose = false

	if debug {
		log.SetLevel(log.TraceLevel)
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

	timeout := time.Millisecond * 100

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	calWrite := make(chan reconws.WsMessage)
	calRead := make(chan reconws.WsMessage)

	switchWrite := make(chan reconws.WsMessage)
	switchRead := make(chan reconws.WsMessage)

	streamWrite := make(chan reconws.WsMessage)
	streamRead := make(chan reconws.WsMessage)

	mdc := drain.NewWs(calRead, ctx)
	mdr := drain.NewWs(switchRead, ctx)
	mds := drain.NewWs(streamRead, ctx)

	// Create test server with the channel handler.
	sc := httptest.NewServer(http.HandlerFunc(serviceChannelHandler(t, calWrite, calRead, ctx)))
	defer sc.Close()

	// Create test server with the channel handler.
	sr := httptest.NewServer(http.HandlerFunc(serviceChannelHandler(t, switchWrite, switchRead, ctx)))
	defer sr.Close()

	// Create test server with the channel handler.
	ss := httptest.NewServer(http.HandlerFunc(userChannelHandler(t, streamWrite, streamRead, ctx)))
	defer ss.Close()

	// URL only assigned after starting
	// Convert http://127.0.0.1 to ws://127.0.0.
	uc := "ws" + strings.TrimPrefix(sc.URL, "http")
	ur := "ws" + strings.TrimPrefix(sr.URL, "http")
	us := "ws" + strings.TrimPrefix(ss.URL, "http")

	New(uc, ur, us, ctx)

	//fmt.Printf("%v\n", middle)

	//time.Sleep(time.Second)
	//t.Logf("counts: %d/%d/%d\n", mdc.Count(), mdr.Count(), mds.Count())

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
	case request := <-mds.Next():

		fmt.Printf(string((request.(reconws.WsMessage)).Data))

		m, ok := request.(reconws.WsMessage)

		assert.True(t, ok)

		var rr pocket.ReasonableFrequencyRange

		err := json.Unmarshal(m.Data, &rr)

		assert.NoError(t, err)

		assert.True(t, ok)

		assert.Equal(t, "rr", rr.Command.Command)
	}

	// to avoid compiler errors for not using them (yet)
	mdc.Count()
	mdr.Count()

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
