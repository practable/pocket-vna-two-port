package calibrate

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/practable/pocket-vna-two-port/pkg/pb"
	"google.golang.org/grpc"
)

func TestCalibrate(t *testing.T) {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:9001")

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalibrateClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CalibrateTwoPort(ctx, &pb.CalibrateTwoPortRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

}
