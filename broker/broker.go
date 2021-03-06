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
	pb3 "github.com/sirbernal/lab3SD/proto/dns_service"
	"google.golang.org/grpc"
)
var timeout = time.Duration(1)*time.Second //definimos el timeout
var id = int64(0) // variable contador de ids para asignarles a los admins
var idadm []int64 // se guarda en ip del dns en el index correspondiente al id de algun admin 
//var dns = []string{"localhost:50052","localhost:50053","localhost:50054"} 
var dns = []string{"10.10.28.82:50052","10.10.28.83:50053","10.10.28.84:50054"} // direcciones de los dns
var randomizer []int
var randoclient []int
type server struct {
}
// Funcion donde se recibe la solicitud del cliente al broker
func (s *server) GetIP(ctx context.Context, msg *pb.GetIPRequest) (*pb.GetIPResponse, error) { 
	// Como se buscara la IP del sitio solicitado en cada dns, con la funcion de abajo 
	// haremos que el orden de conexion de los DNS seran aleatorios
	upRandomizerClient() 
	for _,dire:=range randoclient{ // Recorremos el arreglo generado por la funcion de arriba
		conn, err := grpc.Dial(dns[dire], grpc.WithInsecure()) //genera la conexion con el DNS correspondiente
		if err != nil {
			fmt.Println("Problemas al hacer conexion")
			continue
		}
		defer conn.Close()
		client := pb3.NewDNSServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		msg:= &pb3.GetIPBrokerRequest{Direccion: msg.GetDireccion()} // Le pasamos al dns la direccion correspondiente
		resp, err := client.GetIPBroker(ctx, msg) // aqui recibimos la respuesta del DNS
		if err != nil {
			fmt.Println("Error, dns no esta conectado ")
			continue
		}
		if resp.GetClock()[0]==-1{ // Este es el caso donde no se encuentra el sitio en el dns correspondiente, 
			continue               // por lo que continuamos buscando
		}
		//Al primer DNS donde encontremos el sitio, retornamos la ip del sitio y su clock
		return &pb.GetIPResponse{Ip: resp.GetIp() ,Clock: resp.GetClock() }, nil 
	}// Si no se logra encontrar en ningun DNS, retornamos esto que indicaria al cliente que no encontramos el sitio
	return &pb.GetIPResponse{Clock: []int64{-1},Ip: ""}, nil 
}
// Esta funcion retorna la ip del dns que se debe conectar el admin correspondiente
func (s *server) Broker(ctx context.Context, msg *pb2.BrokerRequest) (*pb2.BrokerResponse, error) { 
	
	return &pb2.BrokerResponse{Ip: sendDNSRoute(msg.GetAdmId())}, nil
}
// FUNCION QUE NO SE USA ACA, PERO SE DEFINE PARA QUE EL PROGRAMA CORRA
func (s *server) DnsCommand(ctx context.Context, msg *pb2.DnsCommandRequest) (*pb2.DnsCommandResponse, error) {
	return &pb2.DnsCommandResponse{Clock: []int64{} }, nil
}

func upRandomizer(){ //actualiza la designacion al azar para los tres admins que se incorporen
	rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
	randomizer=rand.Perm(3)//genera una permutacion al azar de 0 a 2
}

func upRandomizerClient(){ //actualiza la designacion al azar para los tres admins que se incorporen
	rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
	randoclient=rand.Perm(3)//genera una permutacion al azar de 0 a 2
}
func giveDNS(){//funcion que designa el dns cada vez que se agregue un admin
	designio:= id%3 //verifica que designio al azar le corresponde
	 //guarda la direccion que requiere el admin
	idadm=append(idadm,int64(randomizer[designio])) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
}
func giveDNSReset(idn int){//funcion que designa el dns cada vez que se agregue un admin
	designio:= idn%3 //verifica que designio al azar le corresponde
	 //guarda la direccion que requiere el admin
	idadm[idn]=int64(randomizer[designio]) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
}

// Como se explica en el README, cada vez que se hace un merge se volveran a asignar ips de dns nuevas a cada admin conectado
func resetDNS(){ 
	upRandomizer()
	for i:=0; i<int(id);i++{
		giveDNSReset(i)
	}
}
func sendDNSRoute(admin int64)string{ // retorna la ip del dns correspondiente a un admin en especifico (por id)
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

//Con esta funcion, registramos a un admin y le asignamos un id
func (s *server) RegAdm(ctx context.Context, msg *pb2.RegAdmRequest) (*pb2.RegAdmResponse, error) {
	id_temp := id
	giveDNS()
	fmt.Println(idadm)
	id++
	return &pb2.RegAdmResponse{Id: id_temp }, nil

}
//Con esta funcion el DNS notifica a broker que se realizo un merge
func (s *server) NotifyBroker(ctx context.Context, msg *pb3.NotifyBrokerRequest) (*pb3.NotifyBrokerResponse, error) {
	resetDNS()
	fmt.Println("Merge realizado: ", idadm)
	return &pb3.NotifyBrokerResponse{Resp: "Done!"}  , nil
	
}


////////////////// FUNCIONES DECLARADAS EN BROKER SOLO PARA QUE FUNCIONE EL PROGRAMA, ACA NO HACEN NADA///////////////////////////
func (s *server) ReceiveChanges(ctx context.Context, msg *pb3.ReceiveChangesRequest) (*pb3.ReceiveChangesResponse, error) { //////
																															//////
	return &pb3.ReceiveChangesResponse{Status: "listo"}  , nil																//////
																															//////
}																															//////

func (s *server) SendChanges(ctx context.Context, msg *pb3.SendChangesRequest) (*pb3.SendChangesResponse, error) {			//////

	return &pb3.SendChangesResponse{Dominios: []string{}}  , nil
}
func (s *server) GetIPBroker(ctx context.Context, msg *pb3.GetIPBrokerRequest) (*pb3.GetIPBrokerResponse, error) {
	return &pb3.GetIPBrokerResponse{Clock: []int64{},Ip: ""}  , nil
}

func (s *server) GetClock(ctx context.Context, msg *pb3.GetClockRequest) (*pb3.GetClockResponse, error) {					//////
	return &pb3.GetClockResponse{Clock: []int64{}}  , nil																	//////
}																															//////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
	pb3.RegisterDNSServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


}