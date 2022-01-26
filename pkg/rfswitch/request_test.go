package rfswitch

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

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
//	// Convert http://127.0.0.1 to ws://127.0.0.
//	u := "ws" + strings.TrimPrefix(s.URL, "http")
//
//	rf := New(u, ctx)
//
//	mt := int(websocket.TextMessage)
//
//	/* Test ReasonableFrequencyRange */
//	message := []byte("{\"cmd\":\"rr\",\"id\":\"xyz123\"}")
//
//	toClient <- reconws.WsMessage{
//		Data: message,
//		Type: mt,
//	}
//
//	select {
//	case reply := <-fromClient:
//
//		rr := pocket.ReasonableFrequencyRange{}
//
//		err := json.Unmarshal(reply.Data, &rr)
//
//		if err != nil {
//			t.Error("Cannot marshal response to rr command")
//		}
//
//		assert.Equal(t, "xyz123", rr.ID)
//		// weak test - with real kit attached, we should get non-zero numbers
//		assert.True(t, rr.Result.Start > 0)
//		assert.True(t, rr.Result.End > rr.Result.Start)
//
//	case <-time.After(timeout):
//		t.Error("timeout waiting for reply to rr command")
//	}
//
//}
//
func TestPipeInterfaceToWs(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeInterfaceToWs(chanInterface, chanWs, ctx)

	/* Test Request */

	chanInterface <- Command{
		Set: "port",
		To:  "short",
	}

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanWs:

		expected := "{\"set\":\"port\",\"to\":\"short\"}"

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

	/* Test Report */
	message := []byte("{\"report\":\"port\",\"is\":\"short\"}")

	ws := reconws.WsMessage{
		Data: message,
		Type: mt,
	}

	chanWs <- ws

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanInterface:
		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(Report{}))
		report := reply.(Report)
		assert.Equal(t, "short", report.Is)
	}

}

// Reply with expected response if port command is correctly formed
/*
func switchMock(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	mt := int(websocket.TextMessage)

	message := []byte("{\"report\":\"error\",\"is\":\"unrecognised command\"}")

	for {

		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		var r Report

		err := json.Unmarshal([]byte(msg.Data), &r)

		if r.Report == "port" {
			switch r.Is {
			case "short":
				message = []byte("{\"report\":\"port\",\"is\":\"short\"}")
			case "open":
				message = []byte("{\"report\":\"port\",\"is\":\"open\"}")
			case "load":
				message = []byte("{\"report\":\"port\",\"is\":\"load\"}")
			case "dut":
				message = []byte("{\"report\":\"port\",\"is\":\"dut\"}")
			default:
				message = []byte("{\"report\":\"error\",\"is\":\"unrecognised port\"}")
			}
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
*/

func channelHandler(toClient, fromClient chan reconws.WsMessage, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {

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
				return
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
