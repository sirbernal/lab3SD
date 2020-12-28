package main

import (
	"net"
	"fmt"
	"context"
	"log"
	"time"
	"math/rand"
	pb "github.com/sirbernal/lab3SD/proto/client_service"
	pb2 "github.com/sirbernal/lab3SD/proto/admin_service"
	//pb3 "github.com/sirbernal/lab3SD/proto/dns_service"
	"google.golang.org/grpc"
)
var timeout = time.Duration(1)*time.Second
var id = int64(0)
var idadm []int64
var dns = []string{"localhost:50052","localhost:50053","localhost:50054"}
var randomizer []int
type server struct {
}

func (s *server) GetIP(ctx context.Context, msg *pb.GetIPRequest) (*pb.GetIPResponse, error) {
	fmt.Println(msg.GetDireccion())
	return &pb.GetIPResponse{Ip: "192.168.0.1",Clock: []int64{} }, nil
}

func (s *server) Broker(ctx context.Context, msg *pb2.BrokerRequest) (*pb2.BrokerResponse, error) {
	/*
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure()) 
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb3.NewDNSServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg2 := &pb3.GetClockRequest{Soli: " "} 
	clock, err := client.GetClock(ctx, msg2)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}
	*/
	return &pb2.BrokerResponse{Ip: sendDNSRoute(msg.GetAdmId())}, nil
}

func (s *server) DnsCommand(ctx context.Context, msg *pb2.DnsCommandRequest) (*pb2.DnsCommandResponse, error) {
	return &pb2.DnsCommandResponse{Clock: []int64{} }, nil
}

func upRandomizer(){ //actualiza la designacion al azar para los tres admins que se incorporen
	rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
	randomizer=rand.Perm(3)//genera una permutacion al azar de 0 a 2
}
func giveDNS(){//funcion que designa el dns cada vez que se agregue un admin
	designio:= id%3 //verifica que designio al azar le corresponde
	 //guarda la direccion que requiere el admin
	idadm=append(idadm,int64(randomizer[designio])) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
}
func resetDNS(){
	upRandomizer()
	idadm=[]int64{}
	for i:=0; i<=int(id);i++{
		giveDNS()
	}
}
func sendDNSRoute(admin int64)string{
	return dns[idadm[admin]]
}

/*
func (s *server) Append(ctx context.Context, msg *pb2.AppendRequest) (*pb2.AppendResponse, error) {
	return &pb2.AppendResponse{Status : "recibido" }, nil
}*/

/* func (s *server) Update(ctx context.Context, msg *pb2.UpdateRequest) (*pb2.UpdateResponse, error) {
	
	if msg.GetOption() == "ip"{
		return &pb2.UpdateResponse{Status : "Ip cambiada" }, nil
	}else{
		return &pb2.UpdateResponse{Status : "direccion cambiada" }, nil
	}
	return &pb2.UpdateResponse{Status : "recibido" }, nil
} */

/* func (s *server) Delete(ctx context.Context, msg *pb2.DeleteRequest) (*pb2.DeleteResponse, error) {
	return &pb2.DeleteResponse{Status : "recibido" }, nil
}*/


func (s *server) RegAdm(ctx context.Context, msg *pb2.RegAdmRequest) (*pb2.RegAdmResponse, error) {
	id_temp := id
	giveDNS()
	fmt.Println(idadm)
	id++
	return &pb2.RegAdmResponse{Id: id_temp }, nil

}


func main() {
	upRandomizer()
	fmt.Println("Broker corriendo")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterClientServiceServer(s, &server{})
	pb2.RegisterAdminServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


}