package main

import (
	"flag"
	"fmt"
	"grpcdemo/src/grpcdemo/pb"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "Command to Run")
	flag.Parse()

	fmt.Printf("Option is %v\n", *option)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("localhost"+port, opts...)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)

	switch *option {
	case 1:
		SendMetadata(client)
	case 2:
		GetByBadgeNumber(client)
	case 3:
		GetAll(client)
	case 4:
		AddPhoto(client)
	case 5:
		SaveAll(client)
	}
}

func SaveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		pb.Employee{
			BadgeNumber:         123,
			FirstName:           "John",
			LastName:            "Smith",
			VacationAccrualRate: 1.2,
			VacationAccrued:     0,
		},
		pb.Employee{
			BadgeNumber:         234,
			FirstName:           "Lisa",
			LastName:            "Wu",
			VacationAccrualRate: 1.7,
			VacationAccrued:     10,
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})

	// Receive server messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneCh <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-doneCh
}

func AddPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("Penguins.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	md := metadata.New(map[string]string{"badgenumber": "2080"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.AddPhoto(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)
}

func GetAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})

	if err != nil {
		log.Fatal(err)
	}

	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(emp.Employee)
	}
}

func GetByBadgeNumber(client pb.EmployeeServiceClient) {
	res, err := client.GetByBadgeNumber(context.Background(), &pb.GetByBadgeNumberRequest{BadgeNumber: 2080})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Employee)
}

func SendMetadata(client pb.EmployeeServiceClient) {
	md := metadata.MD{}
	md["user"] = []string{"diego"}
	md["password"] = []string{"password1"}

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, md)
	_, err := client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{})

	if err != nil {
		log.Fatal(err)
	}
}
