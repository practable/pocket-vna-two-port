package calibration

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

func TestNew(t *testing.T) {

	timeout := 100 * time.Millisecond

	toClient := make(chan reconws.WsMessage)
	fromClient := make(chan reconws.WsMessage)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create test server with the channel handler.
	s := httptest.NewServer(http.HandlerFunc(channelHandler(toClient, fromClient, ctx)))

	defer s.Close()

	ctx_mock, cancel_mock := context.WithCancel(context.Background())

	go calMock(fromClient, toClient, ctx_mock)

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	cal := New(u, ctx)

	defer cancel_mock()

	// "use" timeout
	select {
	case <-time.After(timeout):
	default:
	}

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
	assert.NoError(t, err)

	popen, err := makeSParam(freq, open_real, open_imag)
	assert.NoError(t, err)

	pload, err := makeSParam(freq, load_real, load_imag)
	assert.NoError(t, err)

	pdut, err := makeSParam(freq, dut_real, dut_imag)
	assert.NoError(t, err)

	err = cal.SetShortParam(pshort)
	assert.NoError(t, err)

	err = cal.SetLoadParam(pload)
	assert.NoError(t, err)

	err = cal.SetOpenParam(popen)
	assert.NoError(t, err)

	err = cal.SetDUTParam(pdut)
	assert.NoError(t, err)

	// check data is correctly assigned
	assert.Equal(t, freq, cal.Command.Freq)
	assert.Equal(t, short_real, cal.Command.Short.Real)
	assert.Equal(t, short_imag, cal.Command.Short.Imag)
	assert.Equal(t, open_real, cal.Command.Open.Real)
	assert.Equal(t, open_imag, cal.Command.Open.Imag)
	assert.Equal(t, load_real, cal.Command.Load.Real)
	assert.Equal(t, load_imag, cal.Command.Load.Imag)
	assert.Equal(t, dut_real, cal.Command.DUT.Real)
	assert.Equal(t, dut_imag, cal.Command.DUT.Imag)

	// send data to service

	result, err := cal.Apply()

	assert.NoError(t, err)

	// the mock returns the original DUT data as the supposedly calibrated data
	assert.Equal(t, dut_real, result.S11.Real)
	assert.Equal(t, dut_imag, result.S11.Imag)

}

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

func TestPipeWsToInterface(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeWsToInterface(chanWs, chanInterface, ctx)

	mt := int(websocket.TextMessage)

	/* Test Report */
	freq := []uint64{1, 2}
	real := []float64{0.1, 0.2}
	imag := []float64{0.3, 0.4}

	r := Result{
		Freq: freq,
		S11: ComplexArray{
			Real: real,
			Imag: imag,
		},
	}

	payload, err := json.Marshal(r)

	assert.NoError(t, err)

	ws := reconws.WsMessage{
		Data: payload,
		Type: mt,
	}

	chanWs <- ws

	select {

	case <-time.After(timeout):
		t.Error("timeout awaiting response")
	case reply := <-chanInterface:
		assert.Equal(t, reflect.TypeOf(reply), reflect.TypeOf(Result{}))
		result := reply.(Result)
		assert.Equal(t, freq, result.Freq)
		assert.Equal(t, real, result.S11.Real)
		assert.Equal(t, imag, result.S11.Imag)
	}

}

//
//// Reply with expected response if port command is correctly formed
//
func calMock(request, response chan reconws.WsMessage, ctx context.Context) {

	mt := int(websocket.TextMessage)

	message := []byte("{\"report\":\"error\",\"is\":\"unrecognised command\"}")

	timeout := 1 * time.Millisecond

	for {

		select {

		case <-ctx.Done():
			return

		case msg := <-request:

			var c Command

			err := json.Unmarshal([]byte(msg.Data), &c)

			message = []byte("{\"report\":\"error\",\"is\":\"unrecognised command\"}")

			if err == nil {

				// return the DUT result uncalibrated for convenience

				r := Result{
					Freq: c.Freq,
					S11:  c.DUT,
				}

				msg, err := json.Marshal(r)

				if err != nil {
					message = []byte("{\"report\":\"error\",\"is\":\"cannot create result\"}")
				} else {
					message = msg
				}

			}

			select {

			case <-ctx.Done():
				return

			case <-time.After(timeout):
				//carry on

			case response <- reconws.WsMessage{Data: message, Type: mt}:
				//fmt.Println(string(message))
				// carry on
			}

		}
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
