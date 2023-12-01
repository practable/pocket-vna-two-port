/* package measure gets a rangeQuery from the VNA for a given device, freq range, and S-params

The mock does not use the mock VNA mock switch because it simply returns
whichever set of results is requestd by the rq.What field.

An extension would select only the SParams specified, and help test single measurement calibrated returns from the middle ware.
*/

package measure

import (
	"errors"
	"fmt"

	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	"github.com/practable/pocket-vna-two-port/pkg/rfusb"
)

type Measure interface {
	Measure(rq *pocket.RangeQuery) error
}

type Hardware struct {
	Switch rfusb.Switch // expect user to supply a pointer to a Switch instance
	VNA    *pocket.VNA
}
type Mock struct {
	Switch                         rfusb.Switch // expect user to supply a pointer to a Switch instance
	VNA                            *pocket.VNA
	ResultRange                    map[string][]pocket.SParam //for range
	ResultSingle                   map[string]pocket.SParam   //for single
	ResultReasonableFrequencyRange pocket.Range
}

func NewHardware(v *pocket.VNA, s rfusb.Switch) *Hardware {

	return &Hardware{
		Switch: s,
		VNA:    v,
	}
}

func NewMock(v *pocket.VNA, s rfusb.Switch) *Mock {

	return &Mock{
		Switch:       s,
		VNA:          v,
		ResultRange:  make(map[string][]pocket.SParam),
		ResultSingle: make(map[string]pocket.SParam),
	}
}

func (h *Hardware) MeasureRange(rq *pocket.RangeQuery) error {

	if rq == nil {
		return errors.New("nil command")
	}
	err := h.Switch.SetPort(rq.What)

	if err != nil {
		return fmt.Errorf("error setting switch to %s because %s", rq.What, err.Error())
	}

	return (*h.VNA).RangeQuery(rq)

}

func (m *Mock) MeasureRange(rq *pocket.RangeQuery) error {
	if rq == nil {
		return errors.New("nil command")
	}
	if _, ok := m.ResultRange[rq.What]; !ok {
		return fmt.Errorf("no mock result for %s", rq.What)
	}
	rq.Result = m.ResultRange[rq.What]
	return nil

}

func (h *Hardware) MeasureSingle(sq *pocket.SingleQuery) error {

	if sq == nil {
		return errors.New("nil command")
	}
	err := h.Switch.SetPort(sq.What)

	if err != nil {
		return fmt.Errorf("error setting switch to %s because %s", sq.What, err.Error())
	}

	return (*h.VNA).SingleQuery(sq)

}

func (m *Mock) MeasureSingle(sq *pocket.SingleQuery) error {
	if sq == nil {
		return errors.New("nil command")
	}
	if _, ok := m.ResultSingle[sq.What]; !ok {
		return fmt.Errorf("no mock result for %s", sq.What)
	}
	sq.Result = m.ResultSingle[sq.What]
	return nil

}

func (h *Hardware) ReasonableFrequencyRange(rfr *pocket.ReasonableFrequencyRange) error {

	if rfr == nil {
		return errors.New("nil command")
	}

	return (*h.VNA).GetReasonableFrequencyRange(rfr)

}

func (m *Mock) ReasonableFrequencyRange(rfr *pocket.ReasonableFrequencyRange) error {
	if rfr == nil {
		return errors.New("nil command")
	}
	rfr.Result = m.ResultReasonableFrequencyRange
	return nil
}
