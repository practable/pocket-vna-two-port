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
	Switch *rfusb.Switch
	VNA    *pocket.VNA
}
type Mock struct {
	Switch  *rfusb.Switch
	VNA     *pocket.VNA
	Results map[string][]pocket.SParam
}

func NewHardware(v *pocket.VNA, s *rfusb.Switch) *Hardware {

	return &Hardware{
		Switch: s,
		VNA:    v,
	}
}

func NewMock(v *pocket.VNA, s *rfusb.Switch) *Mock {

	return &Mock{
		Switch:  s,
		VNA:     v,
		Results: make(map[string][]pocket.SParam),
	}
}

func (h *Hardware) Measure(rq *pocket.RangeQuery) error {

	if rq == nil {
		return errors.New("nil command")
	}
	err := (*h.Switch).SetPort(rq.What)

	if err != nil {
		return fmt.Errorf("error setting switch to %s because %s", rq.What, err.Error())
	}

	return (*h.VNA).RangeQuery(rq)

}

func (m *Mock) Measure(rq *pocket.RangeQuery) error {
	if rq == nil {
		return errors.New("nil command")
	}
	if _, ok := m.Results[rq.What]; !ok {
		return fmt.Errorf("no mock result for %s", rq.What)
	}
	rq.Result = m.Results[rq.What]
	return nil

}
