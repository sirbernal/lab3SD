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

var dires = []string{"localhost:50051"}

var timeout = time.Duration(1)*time.Second

var id_admin int64

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
/*
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
} */
/*
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
} */
/*
func Delete(url string){
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
*/

func RegAdmin(){
	conn, err := grpc.Dial(dires[0], grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.RegAdmRequest{Reg: "registro"} 

	resp, err := client.RegAdm(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}
	fmt.Println(resp.GetId())
	id_admin = resp.GetId()


}

func CommandtoDNS(action []string,ipdns string){
	conn, err := grpc.Dial(ipdns, grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.DnsCommandRequest{Command: action} 

	_, err = client.DnsCommand(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

}

func Broker()string{
	conn, err := grpc.Dial(dires[0], grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.BrokerRequest{AdmId: id_admin} //envia la consulta por medio de la palabra "status"

	msg2, err := client.Broker(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}

	return msg2.GetIp()

}
func main() {
	RegAdmin()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hola, puede escribir su opcion")
	for{
		// Se descompone el string creado para hacer los envios correspondientes
		text, _ := reader.ReadString('\n')
		comandos := DetectCommand(text)
		str := strings.Split(comandos[1], "\n")
		new_comandos := []string{comandos[0],str[0]}
		
		// primero se solicita la ip del dns al broker 
		ipDNS:= Broker()


		// Segundo, se envia los comandos al dns designado x broker
		fmt.Println(new_comandos[0])
		fmt.Println(new_comandos[1])
		CommandtoDNS(new_comandos,ipDNS)
		comandos = []string{}
		/*
		
		switch comandos[0]{
			case "create","append":
				//fmt.Println("Escribio crear")
				//str := strings.Split(comandos[1], " ")
				Broker("create")
				//Create(str[0], str[1])

			case "update":
				//fmt.Println("Escribio subida")
				//str := strings.Split(comandos[1], " ")
				//option := " "

				/*if strings.Count(str[1], ".") == 3{
					option = "ip"
				}else{
					option = "dire"
				} */

				//Broker("update")
				//Update(str[0], str[1], option)
				
			//case "delete":
				//Broker("delete")
				//fmt.Println("Escribio borrar")
				//Delete(comandos[1])


		
		comandos = []string{}
	}

	fmt.Println(DetectCommand("Hola, que sucede mi pana"))
	fmt.Println("Hola")
	//SolicitarIP()
}