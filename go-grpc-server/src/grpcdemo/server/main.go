package main

import (
	"errors"
	"fmt"
	"grpcdemo/src/grpcdemo/pb"
	"io"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	/*
		cred, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
		if err != nil {
			log.Fatal(err)
		}

		opts := []grpc.ServerOption{grpc.Creds(cred)}
		s := grpc.NewServer(opts...)
	*/
	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, new(employeeService))
	log.Println("Starting server on port " + port)
	s.Serve(lis)
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	for _, e := range employees {
		if e.BadgeNumber == req.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {

	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}

	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		employees = append(employees, *emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}
	fmt.Println(employees)
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if ok {
		fmt.Println(md)
		fmt.Printf("Receiving photo for badge number %v\n", md["badgenumber"][0])
	}

	imgData := []byte{}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received %v bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}
