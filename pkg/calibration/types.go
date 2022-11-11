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
	Scan     interface{}
}

/* Command object definition in python calibration service

   {
    "cmd":"twoportport",
    "freq": np.linspace(1e6,100e6,num=N),
    "short":{
        "s11": {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
         "s12" : {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
        "s21": {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
         "s22" : {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
      },
     "open":{ //as for short  },
     "load":{ // as for short },
     "load":{ // as for short },
     "dut":{ // as for short  }
    }
*/

// Command represents a twoport calibration command (request)
type Command struct {
	Command string   `json:"cmd"`
	Freq    []uint64 `json:"freq"`
	Short   SParam   `json:"short"`
	Open    SParam   `json:"open"`
	Load    SParam   `json:"load"`
	Through SParam   `json:"thru"`
	DUT     SParam   `json:"dut"`
}

/* 2-port Result object definition in python calibration service
{
   "freq": network.f,
   "s11": {
      "real": np.squeeze(network.s_re[:,0,0]),
      "imag": np.squeeze(network.s_im[:,0,0]),
           },
  "s12" : {
      "real": np.squeeze(network.s_re[:,0,1]),
      "imag": np.squeeze(network.s_im[:,0,1]),
     },
 "s21": {
      "real": np.squeeze(network.s_re[:,1,0]),
      "imag": np.squeeze(network.s_im[:,1,0]),
     },
  "s22" : {
      "real": np.squeeze(network.s_re[:,1,1]),
      "imag": np.squeeze(network.s_im[:,1,1]),
     },
   }
*/

type Result struct {
	Freq []uint64     `json:"freq"`
	S11  ComplexArray `json:"s11"`
	S12  ComplexArray `json:"s12"`
	S21  ComplexArray `json:"s21"`
	S22  ComplexArray `json:"s22"`
}

type SParam struct {
	S11 ComplexArray `json:"s11"`
	S12 ComplexArray `json:"s12"`
	S21 ComplexArray `json:"s21"`
	S22 ComplexArray `json:"s22"`
}

type ComplexArray struct {
	Real []float64 `json:"real"`
	Imag []float64 `json:"imag"`
}
