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

// Middle holds config and service pointers
type Middle struct {
	c       *pb.CalibrateClient
	conn    *grpc.ClientConn // calibration
	ctx     context.Context
	h       *measure.Hardware // rf switch & VNA
	s       *stream.Stream    // data stream from user
	timeout time.Duration
	rq      pocket.RangeQuery //current calibration

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

	return Middle{
		c:       &c,
		conn:    conn,
		ctx:     ctx,
		h:       h,
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

				//TODO organise the calibration here

			}

		case pocket.CalibratedRangeQuery:

			//TODO measure and apply calibration

		}
	}()

	select {
	case response := <-r:
		return response.Result, response.Error
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}
}

/*

	case pocket.SingleQuery:

		return err = m.h.Measure(request)
		return <-v.Response

	case pocket.RangeQuery:

		// this type is used for different commands

		rq := request.(pocket.RangeQuery)

		switch rq.Command.Command {

		case "rq", "rangequery":

			log.WithFields(log.Fields{
				"request": rq,
			}).Infof("Middle.HandleRequest with ID: %s", rq.ID)

			return RangeQuery(rq, r, v)

		case "rc", "rangecal":

			log.WithFields(log.Fields{
				"request": rq,
			}).Infof("Middle.HandleRequest with ID: %s", rq.ID)

			return RangeCal(rq, c, r, v)

		default:
			return pocket.CustomResult{
				Message: "Unknown request",
				Command: request,
			}
		}

	case pocket.CalibratedRangeQuery:

		crq := request.(pocket.CalibratedRangeQuery)

		return CalibratedRangeQuery(crq, c, r, v)

	default:
		return pocket.CustomResult{
			Message: "Unknown request",
			Command: request,
		}
	}*/

/*
func (m *Middle) CalibratedRangeQuery(crq pocket.CalibratedRangeQuery) interface{} {

	sc, ok := (c.Scan).(pocket.RangeQuery)

	if !(ok && c.Ready) {
		return pocket.CustomResult{
			Message: "Error. No existing calibration. Please calibrate with rc command",
			Command: crq,
		}
	}

	var err error
	var name string

	switch {
	case crq.What == "short" || crq.What == "s":
		name = "short"
		err = r.SetShort()

	case crq.What == "open" || crq.What == "o":
		name = "open"
		err = r.SetOpen()

	case crq.What == "load" || crq.What == "l":
		name = "load"
		err = r.SetLoad()

	case crq.What == "thru" || crq.What == "t":
		name = "thru"
		err = r.SetThru()

	case crq.What == "dut1" || crq.What == "1":
		name = "dut1"
		err = r.SetDUT1()

	case crq.What == "dut2" || crq.What == "2":
		name = "dut2"
		err = r.SetDUT2()

	case crq.What == "dut3" || crq.What == "3":
		name = "dut3"
		err = r.SetDUT3()

	case crq.What == "dut4" || crq.What == "4":
		name = "dut4"
		err = r.SetDUT4()
	default:
		name = crq.What
		err = fmt.Errorf("unrecognised value of what: %s", name)
	}

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: crq,
		}
	}

	// modify the scan command to select only
	// the sparams specified by the user's
	// crq command. The cal scans had to do all four sparams,
	// the user might not always want all four
	// the calibration routine does not need four params in the dut
	// to work, according to testing in python of TwelveTerm (TDD Nov 2022)
	sc.Select = crq.Select

	v.Request <- sc

	log.Debugf("Scan request %v", sc)

	response := <-v.Response

	log.Debugf("Scan response %v", response)

	rrq, ok := response.(pocket.RangeQuery)

	log.Debugf("Scan response as range query %v", rrq)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result := rrq.Result

	if len(result) != sc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetDUTParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store as DUT: " + err.Error(),
			Command: result,
		}
	}

	// apply calibration to DUT data
	calibrated, err := c.Apply()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error applying calibration to measured data for " + name + ": " + err.Error(),
			// don't include result - not in correct format and will be nil anyway
		}
	}

	sparams, err := calibration.CalibrationToPocket(calibrated)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error converting calibrated data format for " + name + ": " + err.Error(),
			// don't include result - not in correct format and will be nil anyway
		}
	}

	crq.Result = sparams

	return crq

}

func (m *Midlle) RangeQuery(rq pocket.RangeQuery) interface{} {

	var err error
	var name string

	switch {
	case rq.What == "short" || rq.What == "s":
		name = "short"
		err = r.SetShort()

	case rq.What == "open" || rq.What == "o":
		name = "open"
		err = r.SetOpen()

	case rq.What == "load" || rq.What == "l":
		name = "load"
		err = r.SetLoad()

	case rq.What == "thru" || rq.What == "t":
		name = "thru"
		err = r.SetThru()

	case rq.What == "dut1" || rq.What == "1":
		name = "dut1"
		err = r.SetDUT1()

	case rq.What == "dut2" || rq.What == "2":
		name = "dut2"
		err = r.SetDUT2()

	case rq.What == "dut3" || rq.What == "3":
		name = "dut3"
		err = r.SetDUT3()

	case rq.What == "dut4" || rq.What == "4":
		name = "dut4"
		err = r.SetDUT4()
	}

	// throw no error if what is unrecognised, because it will be blank when rq is used by rangecal and calibratedrangequery
	// ideally we'd use this in the same way for all uses, but using rq externally only became necessary for troubleshooting
	// the two-port rig with 8-port switches, so we do it this way to minimise changes elsewhere for now.
	// but do throw error if the value is what is valid
	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rq,
		}
	}

	v.Request <- rq

	log.Debugf("Scan request %v", rq)

	response := <-v.Response

	log.Debugf("Scan response %v", response)

	rrq, ok := response.(pocket.RangeQuery)

	log.Debugf("Scan response as range query %v", rrq)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	return rrq

}

func (m *Middle) RangeCal(rc pocket.RangeQuery) interface{} {

	// clear previous cal
	c.Clear()

	// prepare the scanning command used to measure each standard
	scan := rc
	scan.Command.Command = "rq"

	//save it for the cqr to use later
	c.Scan = scan

	// SHORT

	name := "short"

	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}

	log.Debugf("Middle.RangeCal: setting rfswitch to %s", name)

	err := r.SetShort()

	if err != nil {
		log.Errorf("Middle.RangeCal error setting %s was %s", name, err.Error())
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	} else {
		log.Debug("Middle.RangeCal set short ok")
	}

	log.Debug("Middle.RangeCal requesting scan from VNA")

	v.Request <- scan

	log.Debug("Middle.RangeCal awaiting result from VNA")

	response := <-v.Response

	log.Debug("Middle.RangeCal checking result from VNA")

	log.Debugf("response: %s", response)

	rrq, ok := response.(pocket.RangeQuery)

	if !ok {
		log.Errorf("Middle.RangeCal error with scanning %s", name)
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result := rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetShortParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// OPEN

	name = "open"
	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}
	err = r.SetOpen()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response = <-v.Response

	rrq, ok = response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result = rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetOpenParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// LOAD

	name = "load"
	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}
	err = r.SetLoad()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response = <-v.Response

	rrq, ok = response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result = rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetLoadParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// THRU
	name = "thru"
	scan.Select = pocket.SParamSelect{
		S11: true,
		S12: true,
		S21: true,
		S22: true,
	}
	err = r.SetThru()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response = <-v.Response

	rrq, ok = response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result = rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetThruParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	c.Ready = true

	// send some results back so the success is confirmed with the presence of data
	rc.Result = rrq.Result

	return rc

	// don't return a custom result because with a command, with results, because
	// the json parser can't cope with this, causing failed tests. i.e. AVOID THIS:
	// return pocket.CustomResult{
	//	Message: "Success: SOL Calibration of Port 1 complete",
	//	Command: rc,
	// }

}
*/
