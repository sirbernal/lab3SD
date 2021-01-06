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
func SearchMonotonicRead(pag string,dir string,clk []int64){
	for i,j:= range direcciones{
		if pag==j[0]{
			verif:=VerifMonotonicRead(clocks[i],clk)
			if verif==0{
				direcciones[i]=[]string{pag,dir}
				clocks[i]=clk
				fmt.Println("Página actualizada a versión más reciente, dirección: "+dir)
				return
			}else if verif==1{
				fmt.Println("Versión desactualizada recibida, dirección local: "+j[1])
				return
			}else{
				fmt.Println("Versión al día en sistema, dirección: "+j[1])
				return
			}
		}
	}
	fmt.Println("Agregando página '"+pag+"' a memoria, dirección:"+dir)
	direcciones=append(direcciones,[]string{pag,dir})
	clocks=append(clocks,clk)
}
func SolicitarIP(direccion string){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //genera la conexion con el broker
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
	if resp.GetClock()[0]==-1{
		fmt.Println("Página no encontrada")
		return
	}
	SearchMonotonicRead(direccion,resp.GetIp(),resp.GetClock())
}


func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Hola, puede escribir su opcion: ")
	for{
		// Se descompone el string creado para hacer los envios correspondientes
		text, _ := reader.ReadString('\n')
		str := strings.Split(text, "\n")
		SolicitarIP(str[0])
		fmt.Print("Hola, puede escribir su opcion: ")
	}
}