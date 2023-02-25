package main

import (
	"context"
	"io"
	"log"
	"net"
	"test-protobuf/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum Calling")

	resp := &calculatorpb.SumResponse{
		Result: req.Num1 + req.Num2,
	}
	return resp, nil
}

func (s *server) PrimeNumberDecomposition(
	req *calculatorpb.PNDRequest, 
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer,
	) (error) {
	
	log.Println("PND Calling")
	k := int32(2)
	N := req.GetNumber()

	for N > 1 {
		if (N % k == 0) {
			N = N/k
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
		} else {
			k += 1
		}
	}
	return nil
}

func (s *server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average Calling")
	var sum float32
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("end")
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: float32(sum/float32(count)),
			})
		} else if err != nil {
			log.Fatal(err)
		}
		sum += float32(req.GetNumber())
		log.Println("Get: ", req.GetNumber())
		count += 1
	}
}

func (s *server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("FindMax Calling")
	var max int32
	var count int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("end")
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}

		currentNumber := req.GetNumber()
		
		count += 1
		if count == 1 {
			max = currentNumber
		} else {
			if currentNumber > max {
				max = currentNumber
			}
		}

		stream.Send(&calculatorpb.FindMaxResponse{
			Result: max,
		})

	}
}

func main() {
	log.Println("run")
	lis, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
