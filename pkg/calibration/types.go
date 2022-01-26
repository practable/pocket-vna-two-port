package calibration

import (
	"context"
	"time"

	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

type Calibration struct {
	u        string
	R        *reconws.ReconWs
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
	Command  Command
}

/* Command object definition in python calibration service

   {
    "cmd":"oneport",
    "freq": np.linspace(1e6,100e6,num=N),
    "short":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },
     "open":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },
     "load":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },
     "dut":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            }
    }
*/

// this command only has enough fields to support a oneport calibration
// will need extending for two port, or other more exotic calibrations
type Command struct {
	Command string       `json:"cmd"`
	Freq    []uint64     `json:"freq"`
	Short   ComplexArray `json:"short"`
	Open    ComplexArray `json:"open"`
	Load    ComplexArray `json:"load"`
	DUT     ComplexArray `json:"dut"`
}

/* Result object definition in python calibration service
{
          "freq": network.f,
          "s11": {
                      "real": np.squeeze(network.s_re),
                      "imag": np.squeeze(network.s_im),
                  }
   }
*/

type Result struct {
	Freq []uint64     `json:"freq"`
	S11  ComplexArray `json:"s11"`
}

type ComplexArray struct {
	Real []float64 `json:"real"`
	Imag []float64 `json:"imag"`
}
