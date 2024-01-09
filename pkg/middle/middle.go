// Package middle coordinates the response to user requests that require the use of the rfswitch and calibration services
package middle

import (
	"context"
	"errors"
	"time"

	"github.com/practable/pocket-vna-two-port/pkg/measure"
	"github.com/practable/pocket-vna-two-port/pkg/pb"
	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	"github.com/practable/pocket-vna-two-port/pkg/rfusb"
	"github.com/practable/pocket-vna-two-port/pkg/stream"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Ready struct {
	Setup bool
	Short bool
	Open  bool
	Load  bool
	Thru  bool
}

// Middle holds config and service pointers
type Middle struct {
	c       *pb.CalibrateClient
	conn    *grpc.ClientConn // calibration
	ctx     context.Context
	h       *measure.Hardware // rf switch & VNA
	s       *stream.Stream    // data stream from user
	timeout time.Duration
	rq      *pocket.RangeQuery //current calibration
	short   []pocket.SParam
	open    []pocket.SParam
	load    []pocket.SParam
	thru    []pocket.SParam
	dut     []pocket.SParam
	dutcal  []pocket.SParam
	ctpr    *pb.CalibrateTwoPortRequest
	ready   Ready
}

// for the channel in Handle
type Response struct {
	Result interface{}
	Error  error
}

// func New returns a new middleware - do this way so in Run we can call Handle without passing parameters to it
// addr is the host:port of the local gRPC calibration service (unlikely to be remote due to difficulties in proxying HTTP/2)
// port is the usb port for the rf switch, e.g. `/dev/ttyUSB0`
// baud is usb port baud e.g. 57600
// timeoutUSB is the timeout for USB comms e.g. 2m TODO is this needed?
// topic is the address for the stream to connect to at the local `relay host` e.g. ws://localhost:8888/data (TODO check this address for correct format, e.g. does it need the ws://?)

func New(ctx context.Context, addr, port string, baud int, timeoutUSB, timeoutRequest time.Duration, topic string, v *pocket.VNA) Middle {

	// open the serial connection to the rf switch
	r := rfusb.NewRFUSB()
	r.Open(port, baud, timeoutUSB)
	// r.Close() is in Run()

	// create a new measure.Hardware using the rfswitch and VNA
	// note that vna has it's own context (same parent as this context though)
	h := measure.NewHardware(v, r)

	// open the gRPC connection to the calibration service
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to calibration gRPC service %s because %v", addr, err)
	}
	// conn.Close() is in Run()

	c := pb.NewCalibrateClient(conn) //this doesn't need closing, apparently.

	// open the command/data stream to the user (via relay etc)
	s := stream.New(ctx, topic)

	ctpr := &pb.CalibrateTwoPortRequest{}
	ctpr.Reset()

	return Middle{
		c:       &c,
		conn:    conn,
		ctpr:    ctpr,
		ctx:     ctx,
		h:       h,
		ready:   Ready{},
		s:       &s,
		timeout: timeoutRequest,
	}

}

func (m *Middle) Run() {

	defer m.h.Switch.Close()
	defer m.conn.Close()

	for {

		select {

		case request := <-m.s.Request:

			rctx, cancel := context.WithTimeout(m.ctx, m.timeout)

			var response interface{}

			response, err := m.Handle(rctx, request)

			if err != nil {
				response = pocket.CustomResult{
					Message: err.Error(),
					Command: request,
				}
			}

			m.s.Response <- response

			cancel()

		case <-m.ctx.Done():
			return
		}

	} //for

}

func (m *Middle) Handle(ctx context.Context, request interface{}) (response interface{}, err error) {

	r := make(chan Response)

	// now try the request
	// any calls that hang will result in a leakage of the associated goro
	// but hopefully small impact compared to whole system hanging
	go func() {

		// add any new commands to pkg/stream.PipeWsToInterface() or else they'll be silently dropped
		// before getting to this handler
		switch request.(type) {

		case pocket.ReasonableFrequencyRange:

			req := request.(pocket.ReasonableFrequencyRange)
			err := m.h.ReasonableFrequencyRange(&req)

			r <- Response{
				Result: req,
				Error:  err,
			}

		// contains request for raw range query OR to do calibration
		case pocket.RangeQuery:

			rq := request.(pocket.RangeQuery)

			switch rq.Command.Command {

			case "rq", "rangequery":

				req := request.(pocket.RangeQuery)
				err := m.h.MeasureRange(&req)
				r <- Response{
					Result: req,
					Error:  err,
				}

			case "rc", "rangecal":
				req := request.(pocket.RangeQuery)
				err := m.CalibrateRange(&req)
				r <- Response{
					Result: req,
					Error:  err,
				}

			case "sc", "setupcal":
				req := request.(pocket.RangeQuery)
				err := m.CalibrateSetup(&req)
				r <- Response{
					Result: pocket.CustomResult{
						Message: "ok",
						Command: req,
					},
					Error: err,
				}
			case "mc", "measurecal":
				req := request.(pocket.RangeQuery)
				err := m.CalibrateMeasure(&req)
				r <- Response{
					Result: pocket.CustomResult{
						Message: "ok",
						Command: req,
					},
					Error: err,
				}
			case "cc", "confirmcal":
				req := request.(pocket.RangeQuery)
				err := m.CalibrateConfirm(&req)
				r <- Response{
					Result: req,
					Error:  err,
				}
			default:
				log.Errorf("pkg/middle: Handle did not understand command %s in RangeQuery", rq.Command.Command)
			}

		case pocket.CalibratedRangeQuery:

			req := request.(pocket.CalibratedRangeQuery)

			err := m.MeasureRangeCalibrated(&req)
			r <- Response{
				Result: req,
				Error:  err,
			}
		default:
			log.Errorf("Uknown command type %s", request)
		}

	}()

	select {
	case response := <-r:
		return response.Result, response.Error
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}
}

// func MeasureRangeCalibrated measures and applies a calibration, returning calibrated results
func (m *Middle) MeasureRangeCalibrated(request *pocket.CalibratedRangeQuery) error {

	if m.rq == nil {
		return errors.New("not calibrated yet")
	}

	if !m.CalibrationOK() {
		return errors.New("not calibrated yet")
	}

	// measure dut set by user
	m.rq.What = request.What

	err := m.h.MeasureRange(m.rq)

	if err != nil {
		return err
	}

	m.dut = m.rq.Result

	//reuse the other parts of the protocol buffer that are already there from the cal
	m.ctpr.Dut = Meas2Cal(m.dut)

	r, err := (*m.c).CalibrateTwoPort(m.ctx, m.ctpr)
	if err != nil {
		log.Fatalf("could not calibrate: %v", err)
	}

	m.dutcal = Cal2Meas(r.GetFrequency(), r.GetResult())

	request.Result = m.dutcal

	return nil

}

// func CalibrateRange performs the calibration measurements
func (m *Middle) CalibrateRange(request *pocket.RangeQuery) error {

	// store frequency range, size, LogDistribution
	// Measure & save SOLT for all S-params
	// return Sparams for the calibrated item that was listed in the What?
	// Avg can be changed without invalidating the cal, so don't save it

	request.What = "thru" //we'll force the return of the thru results for simplicity

	rq := *request //make a local copy of the request to break the link to the original request
	// so it's not changed by future requests coming in
	m.rq = &rq

	m.ready.Setup = false //we've changed the stored frequency range so a step-by-step cal must be restarted
	m.ready.Short = false
	m.ready.Open = false
	m.ready.Load = false
	m.ready.Thru = false

	// we need to measure all Sparams, so ignore user's select settings
	m.rq.Select = pocket.SParamSelect{
		S11: true,
		S12: true,
		S21: true,
		S22: true,
	}

	// measure cal standards

	//short
	m.rq.What = "short"
	err := m.h.MeasureRange(m.rq)

	if err != nil {
		return err
	}

	m.short = m.rq.Result

	// open
	m.rq.What = "open"
	err = m.h.MeasureRange(m.rq)

	if err != nil {
		return err
	}

	m.open = m.rq.Result

	// load
	m.rq.What = "load"
	err = m.h.MeasureRange(m.rq)

	if err != nil {
		return err
	}

	m.load = m.rq.Result

	// thru
	m.rq.What = "thru"
	err = m.h.MeasureRange(m.rq)

	if err != nil {
		return err
	}

	m.thru = m.rq.Result

	// Use the thru for the DUT for the purpose of this cal
	m.dut = m.thru

	// Prepare the cal buffer...
	m.ctpr.Reset()

	m.ctpr.Frequency = Meas2Freq(m.short)

	m.ctpr.Short = Meas2Cal(m.short)
	m.ctpr.Open = Meas2Cal(m.open)
	m.ctpr.Load = Meas2Cal(m.load)
	m.ctpr.Thru = Meas2Cal(m.thru)
	m.ctpr.Dut = Meas2Cal(m.dut)

	r, err := (*m.c).CalibrateTwoPort(m.ctx, m.ctpr)
	if err != nil {
		log.Fatalf("could not calibrate: %v", err)
	}

	m.dutcal = Cal2Meas(r.GetFrequency(), r.GetResult())

	request.Result = m.dutcal

	// set calibration ok
	m.ready.Setup = true
	m.ready.Short = true
	m.ready.Open = true
	m.ready.Load = true
	m.ready.Thru = true

	return nil

}

func Meas2Freq(s []pocket.SParam) []float64 {
	freq := []float64{}

	for _, v := range s {
		freq = append(freq, float64(v.Freq))
	}

	return freq
}

func Meas2Cal(s []pocket.SParam) *pb.SParams {

	var s11, s12, s21, s22 []*pb.Complex

	for _, v := range s {
		s11 = append(s11, &pb.Complex{
			Real: v.S11.Real,
			Imag: v.S11.Imag,
		})
		s12 = append(s12, &pb.Complex{
			Real: v.S12.Real,
			Imag: v.S12.Imag,
		})
		s21 = append(s21, &pb.Complex{
			Real: v.S21.Real,
			Imag: v.S21.Imag,
		})
		s22 = append(s22, &pb.Complex{
			Real: v.S22.Real,
			Imag: v.S22.Imag,
		})

	}

	return &pb.SParams{
		S11: s11,
		S12: s12,
		S21: s21,
		S22: s22,
	}

}

func Cal2Meas(f []float64, s *pb.SParams) []pocket.SParam {

	var ps []pocket.SParam

	for i := range s.S11 {

		p := pocket.SParam{
			Freq: uint64(f[i]),
			S11: pocket.Complex{
				Real: s.S11[i].Real,
				Imag: s.S11[i].Imag,
			},
			S12: pocket.Complex{
				Real: s.S12[i].Real,
				Imag: s.S12[i].Imag,
			},
			S21: pocket.Complex{
				Real: s.S21[i].Real,
				Imag: s.S21[i].Imag,
			},
			S22: pocket.Complex{
				Real: s.S22[i].Real,
				Imag: s.S22[i].Imag,
			},
		}

		ps = append(ps, p)

	}

	return ps

}

// func CalibrateSetup stores the frequency range and distribution
func (m *Middle) CalibrateSetup(request *pocket.RangeQuery) error {

	// store frequency range, size, LogDistribution
	// Measure & save SOLT for all S-params
	// return Sparams for the calibrated item that was listed in the What?
	// Avg can be changed without invalidating the cal, so don't save it

	request.What = "thru" //we'll force the return of the thru results for simplicity

	rq := *request //make a local copy of the request to break the link to the original request
	// so it's not changed by future requests coming in
	m.rq = &rq

	// we need to measure all Sparams, so ignore user's select settings
	m.rq.Select = pocket.SParamSelect{
		S11: true,
		S12: true,
		S21: true,
		S22: true,
	}

	m.ready.Setup = true //this stays true once the first CalibrateSetup has been run
	// this is set false by any other command that would change the frequency range.
	m.ready.Short = false
	m.ready.Open = false
	m.ready.Load = false
	m.ready.Thru = false

	log.Info("pkg/middle: CalibrateSetup")

	return nil

}

// func CalibrateMeasure stores the frequency range and distribution
func (m *Middle) CalibrateMeasure(request *pocket.RangeQuery) error {

	if m.rq == nil {
		return errors.New("not calibrated yet")
	}

	if m.ready.Setup == false {
		return errors.New("calibration not yet setup (use sc or setupcal command)")
	}

	// measure cal standards

	switch request.What {

	case "short":
		m.rq.What = "short"
		err := m.h.MeasureRange(m.rq)

		if err != nil {
			return err
		}

		m.short = m.rq.Result
		m.ready.Short = true

	case "open":
		m.rq.What = "open"
		err := m.h.MeasureRange(m.rq)

		if err != nil {
			return err
		}

		m.open = m.rq.Result
		m.ready.Open = true

	case "load":
		m.rq.What = "load"
		err := m.h.MeasureRange(m.rq)

		if err != nil {
			return err
		}

		m.load = m.rq.Result
		m.ready.Load = true

	case "thru":
		m.rq.What = "thru"
		err := m.h.MeasureRange(m.rq)

		if err != nil {
			return err
		}

		m.thru = m.rq.Result
		m.ready.Thru = true
	}

	return nil

}

func (m *Middle) CalibrationOK() bool {
	return m.ready.Setup && m.ready.Short && m.ready.Open && m.ready.Load && m.ready.Thru
}

func (m *Middle) CalibrateConfirm(request *pocket.RangeQuery) error {

	if m.rq == nil {
		return errors.New("not calibrated yet")
	}

	if m.ready.Setup == false {
		return errors.New("calibration not setup yet")
	}

	if m.ready.Short == false {
		return errors.New("calibration not complete (missing short, maybe others)")
	}
	if m.ready.Open == false {
		return errors.New("calibration not complete (missing open, maybe others)")
	}
	if m.ready.Load == false {
		return errors.New("calibration not complete (missing load, maybe others)")
	}
	if m.ready.Thru == false {
		return errors.New("calibration not complete (missing thru)")
	}

	// Use the thru for the DUT for the purpose of this cal
	m.dut = m.thru

	// Prepare the cal buffer...
	m.ctpr.Reset()

	m.ctpr.Frequency = Meas2Freq(m.short)

	m.ctpr.Short = Meas2Cal(m.short)
	m.ctpr.Open = Meas2Cal(m.open)
	m.ctpr.Load = Meas2Cal(m.load)
	m.ctpr.Thru = Meas2Cal(m.thru)
	m.ctpr.Dut = Meas2Cal(m.dut)

	r, err := (*m.c).CalibrateTwoPort(m.ctx, m.ctpr)
	if err != nil {
		log.Fatalf("could not calibrate: %v", err)
	}

	m.dutcal = Cal2Meas(r.GetFrequency(), r.GetResult())

	request.What = "thru" //let the user know what result we're sending back
	request.Result = m.dutcal

	return nil

}
