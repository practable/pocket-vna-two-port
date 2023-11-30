package calibrate

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/practable/pocket-vna-two-port/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRequest(t *testing.T) {

	// Makes a request to the calibration server but does not check the quality of any data
	// the data quality has already been checked with the test in the python directory

	// This test requires python code to be running:
	// cd py
	// python3 server.py
	// TODO automate this for testing

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalibrateClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctpr := &pb.CalibrateTwoPortRequest{}

	ctpr.Reset()

	ctpr.Frequency = []float64{1e9}
	spFake := pb.SParams{
		S11: []*pb.Complex{&pb.Complex{Imag: 1, Real: 1}},
		S12: []*pb.Complex{&pb.Complex{Imag: 1, Real: 1}},
		S21: []*pb.Complex{&pb.Complex{Imag: 1, Real: 1}},
		S22: []*pb.Complex{&pb.Complex{Imag: 1, Real: 1}},
	}

	ctpr.Short = &spFake
	ctpr.Open = &spFake
	ctpr.Load = &spFake
	ctpr.Thru = &spFake
	ctpr.Dut = &spFake

	r, err := c.CalibrateTwoPort(ctx, ctpr)
	if err != nil {
		log.Fatalf("could not calibrate: %v", err)
	}
	_ = r.GetFrequency()
	_ = r.GetResult()

}
