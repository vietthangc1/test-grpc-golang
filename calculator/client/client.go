package main

import (
	"context"
	"io"
	"log"
	"test-protobuf/calculator/calculatorpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	// CallSum(client)
	// CallPrimeNumberDecomposition(client)
	// CallAverage(client)
	CallFindMax(client)
}

func CallSum(c calculatorpb.CalculatorServiceClient) *calculatorpb.SumResponse {
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.GetResult())

	return resp
}

func CallPrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 100,
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("End")
			return
		}

		log.Println(resp.GetResult())
	}
}

func CallAverage(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	listReqs := []int32{1, 2, 3, 4, 2, 1, 7, 6}

	for _, req := range listReqs {
		err := stream.Send(&calculatorpb.AverageRequest{
			Number: req,
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("heere")

	resp, err := stream.CloseAndRecv()
	if err == io.EOF {
		log.Println("end")
	} else if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

func CallFindMax(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Normal

	// listReqs := []int32{1,2,3,4,2,1,7,6}
	// for _, req := range listReqs {
	// 	err := stream.Send(&calculatorpb.FindMaxRequest{
	// 		Number: req,
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	resp, err := stream.Recv()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(resp.GetResult())
	// 	time.Sleep(500 * time.Millisecond)
	// }

	// err = stream.CloseSend()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Channel

	waitc := make(chan struct{})

	go func() {
		listReqs := []int32{1, 2, 3, 4, 2, 1, 7, 6}
		for _, req := range listReqs {
			err := stream.Send(&calculatorpb.FindMaxRequest{
				Number: req,
			})
			if err != nil {
				log.Fatal(err)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("end")
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			log.Println(resp.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
