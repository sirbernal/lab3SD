package main

import (
	"context"
	"fmt"
	"time"
	"strings"
	pb "github.com/sirbernal/lab3SD/proto/admin_service"
	"google.golang.org/grpc"
	"bufio"
	"os"
)

var timeout = time.Duration(1)*time.Second

func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}

func DetectCommand2(comm string)[]string{
	str:= strings.Fields(comm)
	fmt.Println(str[1])
	return str
}

func Create(url string, ip string){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.AppendRequest{Url: url, Ip: ip } //envia la consulta por medio de la palabra "status"

	resp, err := client.Append(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

	fmt.Println(resp.GetStatus())
}

func Update(url string, ip string, option string){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.UpdateRequest{Url: url, Ip: ip, Option: option } //envia la consulta por medio de la palabra "status"

	resp, err := client.Update(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

	fmt.Println(resp.GetStatus())
}

func Delete(url string, ip string){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.DeleteRequest{Url: url} //envia la consulta por medio de la palabra "status"

	resp, err := client.Delete(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

	fmt.Println(resp.GetStatus())
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hola, puede escribir su opcion")
	for{
		text, _ := reader.ReadString('\n')
		comandos := DetectCommand(text)
		switch comandos[0]{
			case "create","append":
				fmt.Println("Escribio crear")

			case "update":
				fmt.Println("Escribio subida")
				
			case "delete":
				fmt.Println("Escribio borrar")	

		}	
		comandos = []string{}
	}

	fmt.Println(DetectCommand("Hola, que sucede mi pana"))
	fmt.Println("Hola")
	//SolicitarIP()
}