package calibration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

func init() {

	log.SetLevel(log.WarnLevel)

}

var upgrader = websocket.Upgrader{}

//func TestNew(t *testing.T) {
//
//	timeout := 100 * time.Millisecond
//
//	toClient := make(chan reconws.WsMessage)
//	fromClient := make(chan reconws.WsMessage)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	// Create test server with the channel handler.
//	s := httptest.NewServer(http.HandlerFunc(channelHandler(toClient, fromClient, ctx)))
//
//	defer s.Close()
//
//	ctx_mock, cancel_mock := context.WithCancel(context.Background())
//
//	go switchMock(fromClient, toClient, ctx_mock)
//
//	// Convert http://127.0.0.1 to ws://127.0.0.1
//	u := "ws" + strings.TrimPrefix(s.URL, "http")
//
//	c := New(u, ctx)
//
//	ports := []string{"short", "open", "load", "dut"}
//
//	for _, port := range ports {
//
//		c := Command{
//			Set: "port",
//			To:  port,
//		}
//
//		rf.Request <- c
//
//		select {
//		case report := <-rf.Response:
//
//			v, ok := report.(Report)
//
//			assert.True(t, ok)
//
//			assert.Equal(t, "port", v.Report)
//			assert.Equal(t, port, v.Is)
//
//		case <-time.After(timeout):
//			t.Error("timeout waiting for reply to Set port")
//		}
//	}
//
//	// bad port
//	c := Command{
//		Set: "port",
//		To:  "foo",
//	}
//
//	rf.Request <- c
//
//	select {
//	case report := <-rf.Response:
//
//		v, ok := report.(Report)
//
//		assert.True(t, ok)
//
//		assert.Equal(t, "error", v.Report)
//		assert.Equal(t, "unrecognised port", v.Is)
//
//	case <-time.After(timeout):
//		t.Error("timeout waiting for reply to Set port")
//	}
//
//	// bad command
//	c = Command{
//		Set: "bar",
//		To:  "foo",
//	}
//
//	rf.Request <- c
//
//	select {
//	case report := <-rf.Response:
//
//		v, ok := report.(Report)
//
//		assert.True(t, ok)
//
//		assert.Equal(t, "error", v.Report)
//		assert.Equal(t, "unrecognised command", v.Is)
//
//	case <-time.After(timeout):
//		t.Error("timeout waiting for reply to Set port")
//	}
//
//	// not even a  command
//	r := Report{
//		Report: "bar",
//		Is:     "foo",
//	}
//
//	rf.Request <- r
//
//	select {
//	case report := <-rf.Response:
//
//		v, ok := report.(Report)
//
//		assert.True(t, ok)
//
//		assert.Equal(t, "error", v.Report)
//		assert.Equal(t, "unrecognised command", v.Is)
//
//	case <-time.After(timeout):
//		t.Error("timeout waiting for reply to Set port")
//	}
//
//	// check that the commands complete ok when getting the expected response
//
//	err := rf.SetShort()
//
//	assert.NoError(t, err)
//
//	err = rf.SetOpen()
//
//	assert.NoError(t, err)
//
//	err = rf.SetLoad()
//
//	assert.NoError(t, err)
//
//	err = rf.SetDUT()
//
//	assert.NoError(t, err)
//
//	cancel_mock() //check commands are asking for the right port
//
//	go rf.SetShort()
//
//	msg := <-fromClient
//	err = json.Unmarshal([]byte(msg.Data), &c)
//	assert.NoError(t, err)
//	assert.Equal(t, "port", c.Set)
//	assert.Equal(t, "short", c.To)
//
//	go rf.SetOpen()
//
//	msg = <-fromClient
//	err = json.Unmarshal([]byte(msg.Data), &c)
//	assert.NoError(t, err)
//	assert.Equal(t, "port", c.Set)
//	assert.Equal(t, "open", c.To)
//
//	go rf.SetLoad()
//
//	msg = <-fromClient
//	err = json.Unmarshal([]byte(msg.Data), &c)
//	assert.NoError(t, err)
//	assert.Equal(t, "port", c.Set)
//	assert.Equal(t, "load", c.To)
//
//	go rf.SetDUT()
//
//	msg = <-fromClient
//	err = json.Unmarshal([]byte(msg.Data), &c)
//	assert.NoError(t, err)
//	assert.Equal(t, "port", c.Set)
//	assert.Equal(t, "dut", c.To)
//
//}
//

func makeCommand() (Command, error) {

	freq := []uint64{1000, 2000, 3000}
	short_real := []float64{0.31, 0.32, -0.33}
	short_imag := []float64{0.41, -0.42, 0.45}
	open_real := []float64{1.31, 1.32, -1.33}
	open_imag := []float64{1.41, -1.42, 1.45}
	load_real := []float64{2.31, 2.32, -2.33}
	load_imag := []float64{2.41, -2.42, 2.45}
	dut_real := []float64{3.31, 3.32, -3.33}
	dut_imag := []float64{3.41, -3.42, 3.45}

	pshort, err := makeSParam(freq, short_real, short_imag)
	if err != nil {
		return Command{}, err
	}

	popen, err := makeSParam(freq, open_real, open_imag)
	if err != nil {
		return Command{}, err
	}

	pload, err := makeSParam(freq, load_real, load_imag)
	if err != nil {
		return Command{}, err
	}

	pdut, err := makeSParam(freq, dut_real, dut_imag)
	if err != nil {
		return Command{}, err
	}

	return MakeOnePort(pshort, popen, pload, pdut)

}

func TestPipeInterfaceToWs(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeInterfaceToWs(chanInterface, chanWs, ctx)

	/* Test Request */

	c, err := makeCommand()

	assert.NoError(t, err)

	chanInterface <- c

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanWs:

		var cc Command

		err := json.Unmarshal([]byte(reply.Data), &cc)

		assert.NoError(t, err)

		assert.True(t, cmp.Equal(c, cc))
	}

}

//func TestPipeWsToInterface(t *testing.T) {
//	timeout := 100 * time.Millisecond
//
//	chanWs := make(chan reconws.WsMessage)
//	chanInterface := make(chan interface{})
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	go PipeWsToInterface(chanWs, chanInterface, ctx)
//
//	mt := int(websocket.TextMessage)
//
//	/* Test Report */
//	message := []byte("{\"report\":\"port\",\"is\":\"short\"}")
//
//	ws := reconws.WsMessage{
//		Data: message,
//		Type: mt,
//	}
//
//	chanWs <- ws
//
//	select {
//
//	case <-time.After(timeout):
//		t.Error("timeout awaiting response")
//	case reply := <-chanInterface:
//		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(Result{}))
//		report := reply.(Result)
//		assert.Equal(t, "short", Result.Is)
//	}
//
//}

//
//// Reply with expected response if port command is correctly formed
//
//func CalMock(request, response chan reconws.WsMessage, ctx context.Context) {
//
//	mt := int(websocket.TextMessage)
//
//	message := []byte("{\"report\":\"error\",\"is\":\"unrecognised command\"}")
//
//	timeout := 1 * time.Millisecond
//
//	for {
//
//		select {
//
//		case <-ctx.Done():
//			return
//
//		case msg := <-request:
//
//			//fmt.Println(string(msg.Data))
//
//			var c Command
//
//			err := json.Unmarshal([]byte(msg.Data), &c)
//
//			message = []byte("{\"report\":\"error\",\"is\":\"unrecognised command\"}")
//
//			if err == nil {
//
//				if c.Set == "port" {
//					switch c.To {
//					case "short":
//						message = []byte("{\"report\":\"port\",\"is\":\"short\"}")
//					case "open":
//						message = []byte("{\"report\":\"port\",\"is\":\"open\"}")
//					case "load":
//						message = []byte("{\"report\":\"port\",\"is\":\"load\"}")
//					case "dut":
//						message = []byte("{\"report\":\"port\",\"is\":\"dut\"}")
//					default:
//						message = []byte("{\"report\":\"error\",\"is\":\"unrecognised port\"}")
//					}
//				}
//			}
//
//			select {
//
//			case <-ctx.Done():
//				return
//
//			case <-time.After(timeout):
//				//carry on
//
//			case response <- reconws.WsMessage{Data: message, Type: mt}:
//				//fmt.Println(string(message))
//				// carry on
//			}
//
//		}
//	}
//}
//
//func channelHandler(toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		timeout := 1 * time.Millisecond
//
//		c, err := upgrader.Upgrade(w, r, nil)
//		if err != nil {
//			fmt.Println("Cannot upgrade")
//			return
//		}
//		defer c.Close()
//
//		for { //normal order - receive message over websocket, and reply
//
//			// read from the Client's websocket connection
//			mt, message, err := c.ReadMessage()
//			if err != nil {
//				break
//			}
//
//			// timeout if we don't manage to write to the fromClient channel
//			select {
//			case <-time.After(timeout):
//				return
//			case fromClient <- reconws.WsMessage{Data: message, Type: mt}:
//			}
//
//			select {
//
//			case <-ctx.Done():
//				return
//
//			case <-time.After(timeout):
//				//carry on
//
//			case msg := <-toClient:
//
//				err = c.WriteMessage(msg.Type, msg.Data)
//				if err != nil {
//					break
//				}
//
//			} //select
//
//		} // for
//
//	} //anon func
//
//}
//
