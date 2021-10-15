/*
Package stream connects to a websocket server and transfers JSON messages corresponding to the
types in pkg/pocket i.e. commands and results for pocketVNA

*/

package stream

import (
	"context"

	"github.com/timdrysdale/go-pocketvna/pocket"
)

func Run(c chan pocket.Pocket, ctx context.Context) {

}
