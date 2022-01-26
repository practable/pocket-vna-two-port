package rfswitch

import (
	"context"
	"time"

	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

type Command struct {
	Set string `json:"set"`
	To  string `json:"to"`
}

type Report struct {
	Report string `json:"report"`
	Is     string `json:"is"`
}

type Switch struct {
	u        string
	R        *reconws.ReconWs
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
}
