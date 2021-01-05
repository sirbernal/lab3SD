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
var direcciones [][]string //[[google, ip]...]
var clocks [][]int64 //[clock de google, ...]

func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func VerifMonotonicRead(oldclk []int64,newclk []int64)int{ //0 actualizar pag, 1 llego version vieja, 2 llego la misma version
	igual:=0
	for i:=0;i<3;i++{
		if oldclk[i]>newclk[i]{
			return 1
		}
		if oldclk[i]==newclk[i]{
			igual++
		}
	}
	if igual==3{
		return 2
	}
	return 0
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
		fmt.Println("Problemas al hacer conexion")
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	msg:= &pb.GetIPRequest{Direccion: direccion } //envia la consulta por medio de la palabra "direccion"

	resp, err := client.GetIP(ctx, msg)
	if err != nil {
		fmt.Println("Error, no esta el server conectado ")
	}
	if resp.GetClock()[0]==-1{
		fmt.Println("Página no encontrada")
		return
	}
	SearchMonotonicRead(direccion,resp.GetIp(),resp.GetClock())
}


func main() {

	SolicitarIP("google.cl")
	SolicitarIP("google.cl")
	SolicitarIP("google.cl")
	SolicitarIP("google.cl")
	SolicitarIP("google.cl")
	SolicitarIP("google.cl")
}