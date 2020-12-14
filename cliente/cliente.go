package main

import (
	"context"
	"fmt"
	"time"
	"strings"
	pb "github.com/sirbernal/lab3SD/proto/client_service"
	"google.golang.org/grpc"

)

var timeout = time.Duration(1)*time.Second

func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}

func SolicitarIP(){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.GetIPRequest{Direccion: "google.com" } //envia la consulta por medio de la palabra "status"

	resp, err := client.GetIP(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

	fmt.Println(resp.GetIp())
}


func main() {
	fmt.Println(DetectCommand("Hola, que sucede mi pana"))
	fmt.Println("Hola")
	//SolicitarIP()
}