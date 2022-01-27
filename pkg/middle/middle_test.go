package middle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
	"github.com/timdrysdale/go-pocketvna/pkg/stream"
)

func init() {
	// suppress info messages above closed connections
	// when test servers are stopped
	log.SetLevel(log.WarnLevel)

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

func testVNA(t *testing.T) {

	//timeout := time.Millisecond * 100

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	streamWrite := make(chan reconws.WsMessage)
	streamRead := make(chan reconws.WsMessage)

	calWrite := make(chan reconws.WsMessage)
	calRead := make(chan reconws.WsMessage)

	switchWrite := make(chan reconws.WsMessage)
	switchRead := make(chan reconws.WsMessage)

	// Create test server with the channel handler.
	sc := httptest.NewUnstartedServer(http.HandlerFunc(serviceChannelHandler(t, calWrite, calRead, ctx)))
	defer sc.Close()

	// Create test server with the channel handler.
	sr := httptest.NewUnstartedServer(http.HandlerFunc(serviceChannelHandler(t, switchWrite, switchRead, ctx)))
	defer sr.Close()

	// Create test server with the channel handler.
	ss := httptest.NewServer(http.HandlerFunc(userChannelHandler(t, streamWrite, streamRead, ctx)))
	defer ss.Close()

	// in case we are switching off parts of the middle ware during troubleshooting
	sc.Start()
	sr.Start()
	ss.Start()

	// URL only assigned after starting
	// Convert http://127.0.0.1 to ws://127.0.0.
	//uc := "ws" + strings.TrimPrefix(sc.URL, "http")
	//ur := "ws" + strings.TrimPrefix(sr.URL, "http")
	us := "ws" + strings.TrimPrefix(ss.URL, "http")

	stream := fakeMiddle(us, ctx)

	fmt.Printf("%+v", stream)

	//_ = New(uc, ur, us, ctx)

	time.Sleep(3 * time.Second)

	// Do GetReasonableFrequencyRange command

	//reasonable := Range{}
	//
	//id := "123xyz"
	//v.Request <- ReasonableFrequencyRange{Command: Command{ID: id}}
	//
	//select {
	//case <-time.After(timeout):
	//	t.Error("timeout")
	//case ri := <-v.Response:
	//
	//	if actual, ok := ri.(ReasonableFrequencyRange); !ok {
	//		t.Error("Wrong type returned")
	//	} else {
	//
	//		assert.Equal(t, actual.ID, id)
	//		// weak test - with real kit attached, we should get non-zero numbers
	//		assert.True(t, actual.Result.Start > 0)
	//		assert.True(t, actual.Result.End > actual.Result.Start)
	//		reasonable = actual.Result //save for RangeQuery
	//		if verbose {
	//			fmt.Println(actual.Result)
	//		}
	//	}
	//}
	//
	//// Do SingleQuery command
	//
	//id = "456abc"
	//v.Request <- SingleQuery{
	//	Command: Command{ID: id},
	//	Freq:    200000,
	//	Avg:     1,
	//	Select:  SParamSelect{true, true, true, true},
	//}
	//
	//select {
	//case <-time.After(timeout):
	//	t.Error("timeout")
	//case ri := <-v.Response:
	//
	//	if actual, ok := ri.(SingleQuery); !ok {
	//		t.Error("Wrong type returned")
	//	} else {
	//
	//		assert.Equal(t, actual.ID, id)
	//		// weak test - with real kit attached, we should get non-zero numbers
	//		assert.True(t, actual.Result.S11.Real != 0)
	//		if verbose {
	//			fmt.Println(actual.Result)
	//		}
	//	}
	//}
	//
	//// Do RangeQuery command
	//
	//id = "789def"
	//N := 7 // number of samples
	//v.Request <- RangeQuery{
	//	Command:         Command{ID: id},
	//	Range:           reasonable,
	//	Size:            N,
	//	Avg:             1,
	//	LogDistribution: true,
	//	Select:          SParamSelect{true, true, true, true},
	//}
	//
	//timeout = time.Second //need more time for this than a single query
	//
	//select {
	//case <-time.After(timeout):
	//	t.Error("timeout")
	//case ri := <-v.Response:
	//
	//	if actual, ok := ri.(RangeQuery); !ok {
	//		t.Error("Wrong type returned")
	//	} else {
	//
	//		assert.Equal(t, actual.ID, id)
	//		// weak test - with real kit attached, we should get non-zero numbers
	//		assert.Equal(t, len(actual.Result), N)
	//
	//		assert.Equal(t, reasonable.Start, actual.Result[0].Freq)
	//		assert.Equal(t, reasonable.End, actual.Result[N-1].Freq)
	//
	//		expectedFreq := LogFrequency(reasonable.Start, reasonable.End, N)
	//
	//		for i := 0; i < N; i++ {
	//			if verbose {
	//				fmt.Printf("%d: %d %d\n", i, int(expectedFreq[i]), int(actual.Result[i].Freq))
	//			}
	//			assert.Equal(t, int(expectedFreq[i]), int(actual.Result[i].Freq))
	//		}
	//
	//		if verbose {
	//			fmt.Println(actual.Result)
	//		}
	//	}
	//}
	//

}

func userChannelHandler(t *testing.T, toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		timeout := 1 * time.Millisecond

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
