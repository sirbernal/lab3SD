package main

import (
	"context"
	"fmt"
	"time"
	"strings"
	pb "github.com/sirbernal/lab3SD/proto/client_service"
	"google.golang.org/grpc"
	"bufio"
	"os"
)

var timeout = time.Duration(1)*time.Second //variable para timeout en conexiones
var direcciones [][]string //[[google, ip]...] //guarda las direcciones que han sido solicitadas con su respectiva ip
var clocks [][]int64 //[clock de google, ...] //guarda los relojes asociados a la página

func DetectCommand(comm string)[]string{ // Funcion que sirve para detectar especificamente comandos
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func VerifMonotonicRead(oldclk []int64,newclk []int64)int{ //funcion que verifica el reloj actual de la página en caso de tener 
	igual:=0 
	for i:=0;i<3;i++{
		if oldclk[i]>newclk[i]{
			return 1 //1 llego version vieja
		}
		if oldclk[i]==newclk[i]{
			igual++
		}
	}
	if igual==3{
		return 2 //2 llego la misma version
	}
	return 0 //0 actualizar pag
}
// En el README explicamos con mas detalles como se respeta el Monotonic Read
func SearchMonotonicRead(pag string,dir string,clk []int64){
	for i,j:= range direcciones{ // Recorremos nuestro arreglo en memoria de direcciones consultadas
		if pag==j[0]{ // Si encontramos en nuestro arreglo la pag
			verif:=VerifMonotonicRead(clocks[i],clk) // Comparamos los relojes para saber si esta actualizado o no 
			if verif==0{ // Caso en el que se obtiene version actualizada
				direcciones[i]=[]string{pag,dir}
				clocks[i]=clk
				fmt.Println("Página actualizada a versión más reciente, dirección: "+dir)
				return
			}else if verif==1{ // caso en el que se obtiene version desactualizada
				fmt.Println("Versión desactualizada recibida, dirección local: "+j[1])
				return
			}else{ // ninguno de los dos casos anteriores, simplemente se recibio la misma version que tenemos en el sistema
				fmt.Println("Versión al día en sistema, dirección: "+j[1])
				return
			}
		}
	}
	fmt.Println("Agregando página '"+pag+"' a memoria, dirección:"+dir)
	direcciones=append(direcciones,[]string{pag,dir})
	clocks=append(clocks,clk)
}
func SolicitarIP(direccion string){ // con esta funcion contactamos a broker para que nos envie la ip del sitio solicitado
	//conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	conn, err := grpc.Dial("10.10.28.81:50051", grpc.WithInsecure()) //genera la conexion con el broker
	if err != nil {
		fmt.Println("Problemas al hacer conexion con broker, pruebe nuevamente...")
		return
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.GetIPRequest{Direccion: direccion } //envia la consulta por medio de la palabra "direccion"

	resp, err := client.GetIP(ctx, msg)
	if err != nil {
		fmt.Println("Problemas al hacer conexion con broker, pruebe nuevamente...")
		return
	}
	if resp.GetClock()[0]==-1{ // Si no se encontro en el sistema la pagina, printeamos que no lo encontramos
		fmt.Println("Página no encontrada")
		return
	}
	SearchMonotonicRead(direccion,resp.GetIp(),resp.GetClock()) // Si se encuentra, debemos comprobar que la direccion recibida
	// no es mas antigua a la ultima leida, para respetar la consistencia
}


func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Hola, puede escribir su opcion: ")
	Menu:
		for{
			// Se descompone el string creado para hacer los envios correspondientes
			text, _ := reader.ReadString('\n') //
			str := strings.Split(text, "\n")
			str2 := DetectCommand(str[0])
			switch str2[0]{
			case "connect","get":
				if len(strings.Split(str2[1],"."))==2{
					SolicitarIP(str2[1])
				}else{
					fmt.Print("Ingreso de comando no válido, pruebe nuevamente: ")
					continue Menu
				}
			default:
				fmt.Print("Ingreso de comando no válido, pruebe nuevamente: ")
				continue Menu
			}
			
			fmt.Print("Hola, puede escribir su opcion: ")
		}
}