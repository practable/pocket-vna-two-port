package calibration

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

func init() {

	log.SetLevel(log.WarnLevel)

}

var upgrader = websocket.Upgrader{}

type testDataPocket struct {
	s1 []pocket.SParam
	s2 []pocket.SParam
	s  []pocket.SParam
	o1 []pocket.SParam
	o2 []pocket.SParam
	o  []pocket.SParam
	l1 []pocket.SParam
	l2 []pocket.SParam
	l  []pocket.SParam
	t  []pocket.SParam
	d  []pocket.SParam
}

type testDataResult struct {
	s1 Result
	s2 Result
	o1 Result
	o2 Result
	l1 Result
	l2 Result
	t  Result
	d  Result
	z  Result
}

func makeData(t *testing.T) (testDataPocket, testDataResult) {

	// Make some results in a convenient form (lists)
	// make each value unique to check for
	// typos in assignments of values in structs
	freq := []uint64{1000, 2000, 3000}
	_a := []float64{0.1, 0.2, 0.3}
	_b := []float64{0.4, 0.5, 0.6}
	_c := []float64{0.7, 0.8, 0.9}
	_d := []float64{1.0, 1.1, 1.2}
	_e := []float64{1.3, 1.4, 1.5}
	_f := []float64{1.6, 1.7, 1.8}
	_g := []float64{1.9, 2.0, 2.1}
	_h := []float64{2.2, 2.3, 2.4}
	_i := []float64{2.5, 2.6, 2.7}
	_j := []float64{2.8, 2.9, 3.0}
	_k := []float64{3.1, 3.2, 3.3}
	_l := []float64{3.4, 3.5, 3.6}
	_m := []float64{3.8, 3.9, 4.0}
	_n := []float64{4.1, 4.2, 4.3}
	_o := []float64{4.4, 4.5, 4.6}
	_p := []float64{4.7, 4.8, 4.9}
	_q := []float64{5.0, 5.1, 5.2}
	_r := []float64{5.3, 5.4, 5.5}
	_s := []float64{5.6, 5.7, 5.8}
	_t := []float64{5.9, 6.0, 6.1}
	_u := []float64{6.2, 6.3, 6.4}
	_v := []float64{6.5, 6.6, 6.7}
	_w := []float64{6.8, 6.9, 7.0}
	_x := []float64{7.1, 7.2, 7.3}
	_y := []float64{7.4, 7.5, 7.6}
	_z := []float64{7.7, 7.8, 7.9}
	_aa := []float64{8.0, 8.1, 8.2}
	_ab := []float64{8.3, 8.4, 8.5}

	z := make([]float64, 3)
	caz := ComplexArray{
		Imag: z,
		Real: z,
	}

	rz := Result{
		Freq: freq,
		S11:  caz,
		S12:  caz,
		S21:  caz,
		S22:  caz,
	}

	rs1 := rz
	rs1.S11.Imag = _a
	rs1.S11.Real = _b

	rs2 := rz
	rs2.S22.Imag = _c
	rs2.S22.Real = _d

	ro1 := rz
	ro1.S11.Imag = _e
	ro1.S11.Real = _f

	ro2 := rz
	ro2.S22.Imag = _g
	ro2.S22.Real = _h

	rl1 := rz
	rl1.S11.Imag = _i
	rl1.S11.Real = _j

	rl2 := rz
	rl2.S22.Imag = _k
	rl2.S22.Real = _l

	rt := Result{
		Freq: freq,
		S11: ComplexArray{
			Imag: _m,
			Real: _n,
		},
		S12: ComplexArray{
			Imag: _o,
			Real: _p,
		},
		S21: ComplexArray{
			Imag: _q,
			Real: _r,
		},
		S22: ComplexArray{
			Imag: _s,
			Real: _t,
		},
	}

	rd := Result{
		Freq: freq,
		S11: ComplexArray{
			Imag: _u,
			Real: _v,
		},
		S12: ComplexArray{
			Imag: _w,
			Real: _x,
		},
		S21: ComplexArray{
			Imag: _y,
			Real: _z,
		},
		S22: ComplexArray{
			Imag: _aa,
			Real: _ab,
		},
	}

	// convert into pocket form
	pshort1, err := makePocketSParam(rs1)
	assert.NoError(t, err)
	pshort2, err := makePocketSParam(rs2)
	assert.NoError(t, err)
	popen1, err := makePocketSParam(ro1)
	assert.NoError(t, err)
	popen2, err := makePocketSParam(ro2)
	assert.NoError(t, err)
	pload1, err := makePocketSParam(rl1)
	assert.NoError(t, err)
	pload2, err := makePocketSParam(rl2)
	assert.NoError(t, err)
	pthru, err := makePocketSParam(rt)
	assert.NoError(t, err)
	pdut, err := makePocketSParam(rd)
	assert.NoError(t, err)
	// Test cal functions

	pshort := CombineTwoReflectiveStandardsPocket(pshort1, pshort2)
	popen := CombineTwoReflectiveStandardsPocket(popen1, popen2)
	pload := CombineTwoReflectiveStandardsPocket(pload1, pload2)

	tdp := testDataPocket{
		s1: pshort1,
		s2: pshort2,
		s:  pshort,
		o1: popen1,
		o2: popen2,
		o:  popen,
		l1: pload1,
		l2: pload2,
		l:  pload,
		t:  pthru,
		d:  pdut,
	}

	tdr := testDataResult{
		s1: rs1,
		s2: rs2,
		o1: ro1,
		o2: ro2,
		l1: rl1,
		l2: rl2,
		t:  rt,
		d:  rd,
		z:  rz,
	}

	return tdp, tdr
}

func makeCommand(t *testing.T) (Command, error) {

	tdp, _ := makeData(t)

	return MakeTwoPort(tdp.s1, tdp.s2, tdp.o1, tdp.o2, tdp.l1, tdp.l2, tdp.t, tdp.d)

}

func makeSParamS11(freq []uint64, real, imag []float64) ([]pocket.SParam, error) {

	if len(freq) != len(real) || len(freq) != len(imag) {
		return []pocket.SParam{}, errors.New("Freq/real/imag inconsistent length")
	}

	pa := []pocket.SParam{}

	for i := range freq {
		p := pocket.SParam{
			Freq: freq[i],
			S11: pocket.Complex{
				Imag: imag[i],
				Real: real[i],
			},
		}
		pa = append(pa, p)
	}

	return pa, nil
}
func makePocketSParam(r Result) ([]pocket.SParam, error) {

	pa := []pocket.SParam{}

	for i := range r.Freq {
		p := pocket.SParam{
			Freq: r.Freq[i],
			S11: pocket.Complex{
				Imag: r.S11.Imag[i],
				Real: r.S11.Real[i],
			},
			S12: pocket.Complex{
				Imag: r.S12.Imag[i],
				Real: r.S12.Real[i],
			},
			S21: pocket.Complex{
				Imag: r.S21.Imag[i],
				Real: r.S21.Real[i],
			},
			S22: pocket.Complex{
				Imag: r.S22.Imag[i],
				Real: r.S22.Real[i],
			},
		}
		pa = append(pa, p)
	}

	return pa, nil
}

func TestMakeData(t *testing.T) {

	freq := []uint64{1000, 2000, 3000}
	_a := []float64{0.1, 0.2, 0.3}
	_b := []float64{0.4, 0.5, 0.6}
	/*
		_c := []float64{0.7, 0.8, 0.9}
		_d := []float64{1.0, 1.1, 1.2}
		_e := []float64{1.3, 1.4, 1.5}
		_f := []float64{1.6, 1.7, 1.8}
		_g := []float64{1.9, 2.0, 2.1}
		_h := []float64{2.2, 2.3, 2.4}
		_i := []float64{2.5, 2.6, 2.7}
		_j := []float64{2.8, 2.9, 3.0}
		_k := []float64{3.1, 3.2, 3.3}
		_l := []float64{3.4, 3.5, 3.6}
		_m := []float64{3.8, 3.9, 4.0}
		_n := []float64{4.1, 4.2, 4.3}
		_o := []float64{4.4, 4.5, 4.6}
		_p := []float64{4.7, 4.8, 4.9}
		_q := []float64{5.0, 5.1, 5.2}
		_r := []float64{5.3, 5.4, 5.5}
		_s := []float64{5.6, 5.7, 5.8}
		_t := []float64{5.9, 6.0, 6.1}
		_u := []float64{6.2, 6.3, 6.4}
		_v := []float64{6.5, 6.6, 6.7}
		_w := []float64{6.8, 6.9, 7.0}
		_x := []float64{7.1, 7.2, 7.3}
		_y := []float64{7.4, 7.5, 7.6}
		_z := []float64{7.7, 7.8, 7.9}
		_aa := []float64{8.0, 8.1, 8.2}
		_ab := []float64{8.3, 8.4, 8.5}
	*/
	tdp, _ := makeData(t)

	// Check frequencies are correct in sparam results
	assert.Equal(t, freq[0], tdp.s[0].Freq)
	assert.Equal(t, freq[1], tdp.s[1].Freq)
	assert.Equal(t, freq[2], tdp.s[2].Freq)
	assert.Equal(t, freq[0], tdp.o[0].Freq)
	assert.Equal(t, freq[1], tdp.o[1].Freq)
	assert.Equal(t, freq[2], tdp.o[2].Freq)
	assert.Equal(t, freq[0], tdp.l[0].Freq)
	assert.Equal(t, freq[1], tdp.l[1].Freq)
	assert.Equal(t, freq[2], tdp.l[2].Freq)
	assert.Equal(t, freq[0], tdp.t[0].Freq)
	assert.Equal(t, freq[1], tdp.t[1].Freq)
	assert.Equal(t, freq[2], tdp.t[2].Freq)
	assert.Equal(t, freq[0], tdp.d[0].Freq)
	assert.Equal(t, freq[1], tdp.d[1].Freq)
	assert.Equal(t, freq[2], tdp.d[2].Freq)

	// Check selected results
	assert.Equal(t, _a[0], tdp.s[0].S11.Imag)
	assert.Equal(t, _a[1], tdp.s[1].S11.Imag)
	assert.Equal(t, _a[2], tdp.s[2].S11.Imag)
	assert.Equal(t, _b[0], tdp.s[0].S11.Real)
	assert.Equal(t, _b[1], tdp.s[1].S11.Real)
	assert.Equal(t, _b[2], tdp.s[2].S11.Real)

}

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

	// Get some data

	tdp, tdr := makeData(t)

	err := cal.SetShortParam(tdp.s)
	assert.NoError(t, err)

	err = cal.SetLoadParam(tdp.l)
	assert.NoError(t, err)

	err = cal.SetOpenParam(tdp.o)
	assert.NoError(t, err)

	err = cal.SetThruParam(tdp.t)
	assert.NoError(t, err)

	err = cal.SetDUTParam(tdp.d)
	assert.NoError(t, err)

	// check data is correctly assigned
	assert.Equal(t, tdr.s1.Freq, cal.Command.Freq)
	assert.Equal(t, tdr.s1.S11, cal.Command.Short.S11)
	assert.Equal(t, tdr.z.S12, cal.Command.Short.S12)
	assert.Equal(t, tdr.z.S21, cal.Command.Short.S21)
	assert.Equal(t, tdr.s2.S22, cal.Command.Short.S22)

	assert.Equal(t, tdr.o1.S11, cal.Command.Open.S11)
	assert.Equal(t, tdr.z.S12, cal.Command.Open.S12)
	assert.Equal(t, tdr.z.S21, cal.Command.Open.S21)
	assert.Equal(t, tdr.o2.S22, cal.Command.Open.S22)

	assert.Equal(t, tdr.l1.S11, cal.Command.Load.S11)
	assert.Equal(t, tdr.z.S12, cal.Command.Load.S12)
	assert.Equal(t, tdr.z.S21, cal.Command.Load.S21)
	assert.Equal(t, tdr.l2.S22, cal.Command.Load.S22)

	assert.Equal(t, tdr.t, cal.Command.Thru)
	assert.Equal(t, tdr.d, cal.Command.DUT)

	// send data to service

	result, err := cal.Apply()

	assert.NoError(t, err)

	// the mock returns the original DUT data as the supposedly calibrated data
	assert.Equal(t, tdr.d, result)

}

func TestPipeInterfaceToWs(t *testing.T) {
	timeout := 100 * time.Millisecond

	chanWs := make(chan reconws.WsMessage)
	chanInterface := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go PipeInterfaceToWs(chanInterface, chanWs, ctx)

	/* Test Request */

	c, err := makeCommand(t)

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

func TestConvertPocketToCalibration(t *testing.T) {

	tdp, tdr := makeData(t)

	freq, sparam := PocketToCalibration(tdp.t)

	assert.Equal(t, tdr.t.Freq, freq)
	assert.Equal(t, tdr.t.S11, sparam.S11)
	assert.Equal(t, tdr.t.S12, sparam.S12)
	assert.Equal(t, tdr.t.S21, sparam.S21)
	assert.Equal(t, tdr.t.S22, sparam.S22)
}

func TestMakeTwoPort(t *testing.T) {

	tdp, tdr := makeData(t)

	// Check works with good data
	tp, err := MakeTwoPort(tdp.s1, tdp.s2, tdp.o1, tdp.o2, tdp.l1, tdp.l2, tdp.t, tdp.d)

	assert.NoError(t, err)

	assert.Equal(t, tdr.s1.Freq, tp.Freq)
	assert.Equal(t, tdr.s1.S11, tp.Short.S11)
	assert.Equal(t, tdr.z.S12, tp.Short.S12)
	assert.Equal(t, tdr.z.S21, tp.Short.S21)
	assert.Equal(t, tdr.s2.S22, tp.Short.S22)

	assert.Equal(t, tdr.o1.S11, tp.Open.S11)
	assert.Equal(t, tdr.z.S12, tp.Open.S12)
	assert.Equal(t, tdr.z.S21, tp.Open.S21)
	assert.Equal(t, tdr.o2.S22, tp.Open.S22)

	assert.Equal(t, tdr.l1.S11, tp.Load.S11)
	assert.Equal(t, tdr.z.S12, tp.Load.S12)
	assert.Equal(t, tdr.z.S21, tp.Load.S21)
	assert.Equal(t, tdr.l2.S22, tp.Load.S22)

	assert.Equal(t, tdr.t, tp.Thru)
	assert.Equal(t, tdr.d, tp.DUT)

	// Check throws error if data lengths are different
	_, err = MakeTwoPort([]pocket.SParam{}, tdp.s2, tdp.o1, tdp.o2, tdp.l1, tdp.l2, tdp.t, tdp.d)
	assert.Error(t, err)

}

func TestCalibrationToPocket(t *testing.T) {

	tdp, tdr := makeData(t)

	pt, err := CalibrationToPocket(tdr.t)

	assert.NoError(t, err)

	assert.Equal(t, tdp.t, pt)
	/*
		for i := range tdr.t.Freq {
			assert.Equal(t, tdr.t.Freq[i], pt[i].Freq)
			assert.Equal(t, tdr.t.S11[i], pt[i].S11)
			assert.Equal(t, tdr.t.S12[i], pt[i].S12)
			assert.Equal(t, tdr.t.S21[i], pt[i].S21)
			assert.Equal(t, tdr.t.S22[i], pt[i].S22)
		}*/

	badt := tdr.t

	badt.S12 = ComplexArray{}

	_, err = CalibrationToPocket(badt)

	assert.Error(t, err)
	assert.Equal(t, "Freq and S12 are different lengths", err.Error())

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
					S11:  c.DUT.S11,
					S12:  c.DUT.S12,
					S21:  c.DUT.S21,
					S22:  c.DUT.S22,
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
