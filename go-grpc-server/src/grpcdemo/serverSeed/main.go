package main

import (
	"grpcdemo/src/grpcdemo/pb"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	return nil, nil
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	return nil
}
