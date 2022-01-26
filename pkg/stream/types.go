package stream

import (
	"context"
	"time"

	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

type Stream struct {
	u        string
	R        *reconws.ReconWs
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
}
