package calibration

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
//	// Convert http://127.0.0.1 to ws://127.0.0.
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
//func TestPipeInterfaceToWs(t *testing.T) {
//	timeout := 100 * time.Millisecond
//
//	chanWs := make(chan reconws.WsMessage)
//	chanInterface := make(chan interface{})
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	go PipeInterfaceToWs(chanInterface, chanWs, ctx)
//
//	/* Test Request */
//
//	chanInterface <- Command{
//		Set: "port",
//		To:  "short",
//	}
//
//	select {
//
//	case <-time.After(timeout):
//		t.Error("timeout awaiting response")
//	case reply := <-chanWs:
//
//		expected := "{\"set\":\"port\",\"to\":\"short\"}"
//
//		assert.Equal(t, expected, string(reply.Data))
//	}
//
//}
//
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
//		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(Report{}))
//		report := reply.(Report)
//		assert.Equal(t, "short", report.Is)
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
