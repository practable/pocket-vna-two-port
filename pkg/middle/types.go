package middle

import (
	"github.com/timdrysdale/go-pocketvna/pkg/calibration"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/rfswitch"
	"github.com/timdrysdale/go-pocketvna/pkg/stream"
)

type Middle struct {
	Calibration calibration.Calibration
	Stream      stream.Stream
	Switch      rfswitch.Switch
	VNA         pocket.VNAService
}
