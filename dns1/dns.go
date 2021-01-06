package main

import (
	"context"
	"fmt"
	"strings"
	"os"
	"log"
	"net"
	"strconv"
	//pb "github.com/sirbernal/lab3SD/proto/client_service"
	pb2 "github.com/sirbernal/lab3SD/proto/admin_service"
	pb3 "github.com/sirbernal/lab3SD/proto/dns_service"
	"google.golang.org/grpc"

)

type server struct {
}

var idDNS = 1
var dominios []string //arreglo que guarda los dominios en formato ["cl","com",...]
var registro [][]string //arreglo que mantiene los registros por dominio[[registros cl],[registros com],...] => [registros cl]=["append a.cl 1.2.3.4", "update u.cl 2.2.2.2"...]
var pags [][][]string //arreglo que guarda las paginas e ips por dominio[[paginas.cl],[paginas.com],...] =>[paginas.cl]=[["a.cl","1.2.3.4"],["u.cl","2.2.2.2"]]
var clocks [][]int64  //relojes asociados por dominio [[cl],[com],...] => [cl]=[x,y,z]

func DivideData(data string)[]string{//funcion que divide un string por espacio(solo acorta la funcion para escribir menos y ahorrar tiempo)
	return strings.Split(data," ")
}
func DetectDomain(comm string)string{ //funcion que detecta el dominio de una pagina ingresada, por ejemplo abc.com => com
	str:= DivideData(comm)
	dom:=strings.Split(str[0],".")
	return dom[1]
}
func SearchDomain(dom string)int{ //funcion que busca la posicion del dominio en el arreglo guardado
	pos:= -1
	for i,j :=range dominios{
		if j==dom{
			pos=i
			break
		}
	}
	return pos
}
func RemoveIndex(s [][]string, index int) [][]string { //función editada y sacada de https://www.golangprograms.com/how-to-delete-an-element-from-a-slice-in-golang.html
	return append(s[:index], s[index+1:]...)
}
func DetectUpdate(w string)bool{ //funcion que detecta si es cambio de ip o cambio de dominio
	str:= strings.Split(w, ".")
	if len(str)==4{ //se asume que la ip vendrá en formato x.x.x.x
		return false //es cambio de dominio
	}else{
		return true //es cambio de ip
	}
}
func ActReg(pos int){//funcion que actualiza tanto el registro zf del dominio señalado (con el index) como el log de dicho dominio
	file,err:= os.OpenFile("RegistroZF"+dominios[pos]+".txt",os.O_CREATE|os.O_WRONLY,0777) //abre o genera el archivo de registro
	defer file.Close()
	if err !=nil{
		os.Exit(1)
	}
	for _,j:= range pags[pos]{
		word:= j[0]+" IN A "+j[1]
		_, err := file.WriteString(word + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	file.Close()
	file2,err:= os.OpenFile("LOG"+dominios[pos]+".txt",os.O_CREATE|os.O_WRONLY,0777) //abre o genera el archivo de registro
	defer file2.Close()
	if err !=nil{
		os.Exit(1)
	}
	for _,j:= range registro[pos]{
		_, err := file2.WriteString(j + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	file2.Close()
}


func ReceiveOp(op []string)(){ //funcion que recibe el comando de admin o del dns 0 en el merge (solo appends)
	domain:=DetectDomain(op[1])
	values:=DivideData(op[1])
	pos:=SearchDomain(domain)
	switch op[0]{
	case "append":
		if pos==-1{
			dominios=append(dominios,domain)
			registro=append(registro,[]string{})
			pags=append(pags,[][]string{})
			pos=len(dominios)-1
			clocks=append(clocks,[]int64{0,0,0})
		}
		pags[pos]=append(pags[pos],values)
	case "update":
		if pos==-1{
			fmt.Println("Dominio y pagina inexistente, actualización no válida")
			return
		}
		flag:=false
		for i,j :=range pags[pos]{
			if j[0]==values[0]{
				if DetectUpdate(values[1]){
					pags[pos][i]=[]string{values[1],j[1]}
				}else{
					pags[pos][i]=values
				}
				os.Remove("RegistroZF"+domain+".txt")
				flag=true
				break
			}
		}
		if !flag{
			fmt.Println("Pagina a actualizar no existe")
			return
		}
	case "delete":
		flag:=false
		if pos==-1{
			fmt.Println("Dominio y pagina inexistente, eliminación no válida")
			return
		}
		for i,j :=range pags[pos]{
			if j[0]==values[0]{
				pags[pos]=RemoveIndex(pags[pos],i)
				os.Remove("RegistroZF"+domain+".txt")
				flag=true
				break
			}
		}
		if !flag{
			fmt.Println("Pagina a eliminar no existe")
			return
		}
	}
	registro[pos]=append(registro[pos],op[0]+" "+op[1])
	clocks[pos][idDNS]++
	ActReg(pos)
}
func UpdateFiles(){ //funcion que actualiza o reescribe todos los registros zf en todos los dominios despues de un merge
	for i,_:=range dominios{
		file,err:= os.OpenFile("RegistroZF"+dominios[i]+".txt",os.O_CREATE|os.O_WRONLY,0777) //abre o genera el archivo de registro
		defer file.Close()
		if err !=nil{
			os.Exit(1)
		}
		for _,j:= range pags[i]{
			word:= j[0]+" IN A "+j[1]
			_, err := file.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
		file.Close()
	}
}
func CleanLocalMerge(){//funcion que borra todos los archivos locales al inicio o final del merge
	for _,domain:= range dominios{
		os.Remove("RegistroZF"+domain+".txt")
		os.Remove("LOG"+domain+".txt")
	}
}
func InitialClean(){//funcion que borra lo necesario al iniciar el merge
	CleanLocalMerge()
	dominios= []string{}
	pags= [][][]string{}
}
func FinalClean(){//funcion que borra lo necesario al finalizar el merge
	registro= [][]string{}
	clocks= [][]int64{}
	CleanLocalMerge()
	UpdateFiles()
	for _,_= range dominios{
		registro=append(registro,[]string{})
	}
}
func (s *server) GetIPBroker(ctx context.Context, msg *pb3.GetIPBrokerRequest) (*pb3.GetIPBrokerResponse, error) {//retorna la ip de la pagina solitada por el cliente al broker
	pag:=msg.GetDireccion()
	pos:=SearchDomain(DetectDomain(pag)) 
	for _,j:=range pags[pos]{
		if j[0]==pag{
			return &pb3.GetIPBrokerResponse{Clock: clocks[pos],Ip: j[1]}  , nil
		}
	}
	return &pb3.GetIPBrokerResponse{Clock: []int64{-1},Ip: ""}  , nil
}
func (s *server) DnsCommand(ctx context.Context, msg *pb2.DnsCommandRequest) (*pb2.DnsCommandResponse, error) {//funcion que recibe el comando del admin
	ReceiveOp(msg.GetCommand())
	return &pb2.DnsCommandResponse{Clock: []int64{} }, nil
}

func (s *server) Broker(ctx context.Context, msg *pb2.BrokerRequest) (*pb2.BrokerResponse, error) {
	return &pb2.BrokerResponse{Ip: "192.168.0.1"}, nil
}
func (s *server) RegAdm(ctx context.Context, msg *pb2.RegAdmRequest) (*pb2.RegAdmResponse, error) {
	return &pb2.RegAdmResponse{Id: 0 }, nil
}

func (s *server) SendChanges(ctx context.Context, msg *pb3.SendChangesRequest) (*pb3.SendChangesResponse, error) {
	if msg.GetSoli()=="Merge"{
		return &pb3.SendChangesResponse{Dominios: dominios}  , nil
	}else{
		j,_:=strconv.Atoi(msg.GetSoli())
		return &pb3.SendChangesResponse{Dominios:registro[j]}  , nil
	}
	
}
func TranslateClock (clk []string)[]int64{ //recibe el clock en formato string y lo traduce a int64
	a:=[]int64{}
	for _,j:= range clk{
		b,_:=strconv.ParseInt(j,10,64)
		a=append(a,b)
	}
	return a
}

func (s *server) ReceiveChanges(ctx context.Context, msg *pb3.ReceiveChangesRequest) (*pb3.ReceiveChangesResponse, error) {//funcion de recepcion de datos post merge
	if msg.GetType()==0{//recepcion de pagina inicial
		InitialClean()
		ReceiveOp(msg.GetOperations())
	}
	if msg.GetType()==1{//recepcion de pagina normal
		ReceiveOp(msg.GetOperations())
	}
	if msg.GetType()==2{//recepcion de pagina final
		ReceiveOp(msg.GetOperations())
		FinalClean()
	}
	if msg.GetType()==3{//recepcion de clocks
		clocks=append(clocks,TranslateClock(msg.GetOperations()))
	}
	return &pb3.ReceiveChangesResponse{Status: "listo"}  , nil
}

func (s *server) NotifyBroker(ctx context.Context, msg *pb3.NotifyBrokerRequest) (*pb3.NotifyBrokerResponse, error) {
	
	return &pb3.NotifyBrokerResponse{Resp: "Done!"}  , nil
	
}


func main() {
	fmt.Println("DNS 1 en línea")
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()

	pb3.RegisterDNSServiceServer(s, &server{})
	pb2.RegisterAdminServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}