/*

The calibration implementation in Python's scikit-rf meets our needs, so using it avoids us having to re-implement it.

It's possible to consider calling the python code directly e.g. https://poweruser.blog/embedding-python-in-go-338c0399f3d5

But it's more future-proof to try out gRPC because then the implementation of the calibration routine could be done with anything

*/
// package calibrate sends calibration requests over gRPC to a calibration server
package calibrate

/*
import (
	"google.golang.org/grpc"
	pb "github.com/practable/pocket-vna-two-port/pkg/pb"

)


var (
	addr = flag.String("addr", "localhost:9001", "the address to connect to")
)
*/
