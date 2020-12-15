package main

import (
	"net"
	"fmt"
	"context"
	"log"
	pb "github.com/sirbernal/lab3SD/proto/client_service"
	"google.golang.org/grpc"
)


type server struct {
}

func (s *server) GetIP(ctx context.Context, msg *pb.GetIPRequest) (*pb.GetIPResponse, error) {
	fmt.Println(msg.GetDireccion())
	return &pb.GetIPResponse{Ip: "192.168.0.1",Clock: []int64{} }, nil
}


func main() {
	fmt.Println("Broker corriendo")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterClientServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


}