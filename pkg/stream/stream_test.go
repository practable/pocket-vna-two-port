package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/pocket"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/reconws"
)

func init() {

	log.SetLevel(log.WarnLevel)

}

var upgrader = websocket.Upgrader{}

func TestHeartBeat(t *testing.T) {

	c := make(chan reconws.WsMessage)
	d := time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go HeartBeat(c, d, ctx)

	for i := 0; i < 10; i++ {
		select {
		case <-time.After(2 * d):
			t.Error("timeout on heartbeat")
		case msg := <-c:
			assert.Equal(t, []byte("{\"cmd\":\"hb\"}"), msg.Data)

		}
	}

}

func TestNew(t *testing.T) {

	timeout := 100 * time.Millisecond

	toClient := make(chan reconws.WsMessage)
	fromClient := make(chan reconws.WsMessage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create test server with the channel handler.
	s := httptest.NewServer(http.HandlerFunc(channelHandler(toClient, fromClient, ctx)))

	//s := httptest.NewServer(http.HandlerFunc(reasonableRange))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	stream := New(u, ctx)

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

func TestRunDirect(t *testing.T) {

	/* note this test will fail if the first heartbeat comes before the command-response tests
	are completed; could add conditional code to ignore heartbeat commands */

	timeout := 100 * time.Millisecond

	toClient := make(chan reconws.WsMessage)
	fromClient := make(chan reconws.WsMessage)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create test server with the channel handler.
	s := httptest.NewServer(http.HandlerFunc(channelHandler(toClient, fromClient, ctx)))

	//s := httptest.NewServer(http.HandlerFunc(reasonableRange))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	err := pocket.ForceUnlockDevices()

	if err != nil {
		t.Error("Can't unlock devices")
	}

	go RunDirect(u, ctx)

	mt := int(websocket.TextMessage)

	/* Test ReasonableFrequencyRange */
	message := []byte("{\"cmd\":\"rr\",\"id\":\"xyz123\"}")

	toClient <- reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case reply := <-fromClient:

		rr := pocket.ReasonableFrequencyRange{}

		err := json.Unmarshal(reply.Data, &rr)

		if err != nil {
			t.Error("Cannot marshal response to rr command")
		}

		assert.Equal(t, "xyz123", rr.ID)
		// weak test - with real kit attached, we should get non-zero numbers
		assert.True(t, rr.Result.Start > 0)
		assert.True(t, rr.Result.End > rr.Result.Start)

	case <-time.After(timeout):
		t.Error("timeout waiting for reply to rr command")
	}

	/* Test SingleQuery */
	message = []byte("{\"cmd\":\"sq\",\"id\":\"456abc\",\"freq\":200000,\"avg\":1,\"sparam\":{\"S11\":true,\"S12\":false,\"S21\":true,\"S22\":false}}")

	toClient <- reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case reply := <-fromClient:

		sq := pocket.SingleQuery{}

		err := json.Unmarshal(reply.Data, &sq)

		if err != nil {
			t.Error("Cannot marshal response to sq command")
		}

		assert.Equal(t, "456abc", sq.ID)
		// weak test - with real kit attached, we should get non-zero numbers
		assert.True(t, sq.Result.S11.Real != 0)
		assert.Equal(t, uint64(200000), sq.Result.Freq)

	case <-time.After(timeout):
		t.Error("timeout waiting for reply to sq command")
	}

	/* Test RangeQuery */
	message = []byte("{\"id\":\"def789\",\"t\":0,\"cmd\":\"rq\",\"range\":{\"Start\":100000,\"End\":4000000},\"size\":3,\"isLog\":true,\"avg\":1,\"sparam\":{\"S11\":true,\"S12\":false,\"S21\":true,\"S22\":false}}")

	toClient <- reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	select {
	case reply := <-fromClient:

		rq := pocket.RangeQuery{}

		err := json.Unmarshal(reply.Data, &rq)

		if err != nil {
			t.Error("Cannot marshal response to rq command")
		}

		assert.Equal(t, "def789", rq.ID)
		// weak test - with real kit attached, we should get non-zero numbers

		assert.True(t, rq.Result[0].S11.Real != 0)

		assert.Equal(t, len(rq.Result), 3)

		expectedFreq := pocket.LogFrequency(100000, 4000000, 3)

		for i := 0; i < 3; i++ {
			assert.Equal(t, int(expectedFreq[i]), int(rq.Result[i].Freq))
		}

	case <-time.After(10 * timeout):
		t.Error("timeout waiting for reply to rq command")
	}

	/* Test heartbeat */
	select {
	case reply := <-fromClient:

		assert.Equal(t, "{\"cmd\":\"hb\"}", string(reply.Data))

	case <-time.After(5 * time.Second):
		t.Error("No heartbeat")
	}

}

func TestPipeInterfaceToWs(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeInterfaceToWs(chanInterface, chanWs, ctx)

	/* Test ReasonableFrequencyRange */

	chanInterface <- pocket.ReasonableFrequencyRange{
		Command: pocket.Command{Command: "rr"}, Result: pocket.Range{Start: 100000, End: 4000000}}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanWs:

		expected := "{\"id\":\"\",\"t\":0,\"cmd\":\"rr\",\"range\":{\"start\":100000,\"end\":4000000}}"

		assert.Equal(t, expected, string(reply.Data))
	}

	/* Test SingleQuery */
	chanInterface <- pocket.SingleQuery{
		Command: pocket.Command{Command: "sq"},
		Freq:    100000,
		Avg:     1,
		Select:  pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
		Result: pocket.SParam{
			S11: pocket.Complex{Real: -1, Imag: 2},
			S21: pocket.Complex{Real: 0.34, Imag: 0.12},
		},
	}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanWs:

		expected := "{\"id\":\"\",\"t\":0,\"cmd\":\"sq\",\"freq\":100000,\"avg\":1,\"sparam\":{\"s11\":true,\"s12\":false,\"s21\":true,\"s22\":false},\"result\":{\"s11\":{\"real\":-1,\"imag\":2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.34,\"imag\":0.12},\"s22\":{\"real\":0,\"imag\":0},\"freq\":0}}"

		assert.Equal(t, expected, string(reply.Data))
	}

	/* Test RangeQuery */
	chanInterface <- pocket.RangeQuery{
		Command:         pocket.Command{Command: "rq"},
		Range:           pocket.Range{Start: 100000, End: 4000000},
		LogDistribution: true,
		Avg:             1,
		Size:            2,
		Select:          pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false},
		Result: []pocket.SParam{
			pocket.SParam{
				S11: pocket.Complex{Real: -1, Imag: 2},
				S21: pocket.Complex{Real: 0.34, Imag: 0.12},
			},
			pocket.SParam{
				S11: pocket.Complex{Real: -0.1, Imag: 0.2},
				S21: pocket.Complex{Real: 0.3, Imag: 0.4},
			},
		},
	}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanWs:

		expected := "{\"id\":\"\",\"t\":0,\"cmd\":\"rq\",\"range\":{\"start\":100000,\"end\":4000000},\"size\":2,\"islog\":true,\"avg\":1,\"sparam\":{\"s11\":true,\"s12\":false,\"s21\":true,\"s22\":false},\"result\":[{\"s11\":{\"real\":-1,\"imag\":2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.34,\"imag\":0.12},\"s22\":{\"real\":0,\"imag\":0},\"freq\":0},{\"s11\":{\"real\":-0.1,\"imag\":0.2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.3,\"imag\":0.4},\"s22\":{\"real\":0,\"imag\":0},\"freq\":0}]}"

		assert.Equal(t, expected, string(reply.Data))
	}

}

func TestPipeWsToInterface(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeWsToInterface(chanWs, chanInterface, ctx)

	mt := int(websocket.TextMessage)

	/* Test ReasonableFrequencyRange */
	message := []byte("{\"cmd\":\"rr\"}")

	ws := reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	chanWs <- ws

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanInterface:
		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(pocket.ReasonableFrequencyRange{}))
		rr := reply.(pocket.ReasonableFrequencyRange)
		assert.Equal(t, "rr", rr.Command.Command)
	}

	/* Test SingleQuery */
	message = []byte("{\"id\":\"\",\"t\":0,\"cmd\":\"sq\",\"freq\":100000,\"avg\":1,\"sparam\":{\"s11\":true,\"s12\":false,\"s21\":true,\"s22\":false},\"result\":{\"s11\":{\"real\":-1,\"imag\":2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.34,\"imag\":0.12},\"s22\":{\"real\":0,\"imag\":0}}}")

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	chanWs <- ws

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanInterface:
		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(pocket.SingleQuery{}))
		sq := reply.(pocket.SingleQuery)
		assert.Equal(t, "sq", sq.Command.Command)
		assert.Equal(t, uint64(100000), sq.Freq) //not testing Freq in result because this is just a piping test....
		assert.Equal(t, uint16(1), sq.Avg)
		assert.Equal(t, pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false}, sq.Select)
		// no need to check the Sparam results because we are not expecting to pass them in this direction
	}

	/* Test RangeQuery */
	message = []byte("{\"id\":\"\",\"t\":0,\"cmd\":\"rq\",\"range\":{\"start\":100000,\"end\":4000000},\"size\":2,\"islog\":true,\"avg\":1,\"sparam\":{\"s11\":true,\"s12\":false,\"s21\":true,\"s22\":false},\"result\":[{\"s11\":{\"real\":-1,\"imag\":2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.34,\"imag\":0.12},\"s22\":{\"real\":0,\"imag\":0}},{\"s11\":{\"real\":-0.1,\"imag\":0.2},\"s12\":{\"real\":0,\"imag\":0},\"s21\":{\"real\":0.3,\"imag\":0.4},\"s22\":{\"real\":0,\"imag\":0}}]}")

	ws = reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	chanWs <- ws

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanInterface:
		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(pocket.RangeQuery{}))
		rq := reply.(pocket.RangeQuery)
		assert.Equal(t, "rq", rq.Command.Command)
		assert.Equal(t, pocket.Range{Start: 100000, End: 4000000}, rq.Range)
		assert.Equal(t, uint16(1), rq.Avg)
		assert.Equal(t, pocket.SParamSelect{S11: true, S12: false, S21: true, S22: false}, rq.Select)
		assert.Equal(t, true, rq.LogDistribution)

		// no need to check the Sparam results because we are not expecting to pass them in this direction
	}

}

func reasonableRange(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	mt := int(websocket.TextMessage)

	message := []byte("{\"cmd\":\"rr\"}")

	for {

		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("Hello!")
		fmt.Println(message)

	}
}

func channelHandler(toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {

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
				return
			case fromClient <- reconws.WsMessage{Data: message, Type: mt}:
			}

		} // for

	} //anon func

}
