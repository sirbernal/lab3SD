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

var dires = []string{"10.10.28.81:50051"} // direccion del broker

var timeout = time.Duration(1)*time.Second // timeout para los envios

var id_admin int64 // Id de admin que se asignara una vez iniciado

func DetectCommand(comm string)[]string{ // Funcion que sirve para detectar especificamente si se hace un append, update o delete
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}


func RegAdmin(){ // funcion que contacta a broker para que este le asigne una id
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
	id_admin = resp.GetId() // se asigna la id recibida


}
// Funcion que hace el envio a un dns especifico con alguna accion o comando especifico
func CommandtoDNS(action []string,ipdns string){ 
	conn, err := grpc.Dial(ipdns, grpc.WithInsecure()) //genera la conexion con el dns
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

func Broker()string{ // Con esta funcion se contacta a broker para que este le designe una Ip de un DNS al cual debe conectarse
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
	RegAdmin() // Nos registramos con broker y obtenemos una id
	reader := bufio.NewReader(os.Stdin) 
	fmt.Print("Menú Admin\nIngrese comando deseado: ")
	Menu: 
		for{
			// Se lee lo que se escribe en pantalla
			text, _ := reader.ReadString('\n')
			comandos := DetectCommand(text) //descomponemos el string que recibimos
			str := strings.Split(comandos[1], "\n") //quitamos el salto de linea que se nos genera
			new_comandos := []string{comandos[0],str[0]} // generamos el string final que enviaremos
			switch new_comandos[0]{           // dependiendo del tipo de comando, recordemos que tienen mas parametros que otros por lo que
			case "append","update","create":  // es necesario separar en casos especificos
				if len(strings.Split(str[0]," "))!=2{
					fmt.Print("Ingreso de comando no válido, pruebe nuevamente: ")
					continue Menu
				}
			case "delete":
				if len(strings.Split(str[0]," "))!=1{
					fmt.Print("Ingreso de comando no válido, pruebe nuevamente: ")
					continue Menu
				}
			default:
				fmt.Print("Ingreso de comando no válido, pruebe nuevamente: ")
				continue Menu
			}
			// se solicita la ip del dns al broker 
			ipDNS:= Broker()

			// Segundo, se envia los comandos al dns designado x broker
			CommandtoDNS(new_comandos,ipDNS)
			comandos = []string{}
			fmt.Print("\nIngrese comando deseado:")
		}

}