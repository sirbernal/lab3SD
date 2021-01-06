package main

import (
	"context"
	"fmt"
	"time"
	"strings"
//	"bufio"
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

var idDNS = 0
var dominios []string //[com,cl,ez]
var registro [][]string //[[agregar,borrar,...],[],[],]
var pags [][][]string //[[[algo.com,direccion],[xd.com, direccion]],[lel.cl],[]]
var clocks [][]int64  //[[0,0,1],[0,0,3]]
var timeout = time.Duration(1)*time.Second
var dns = []string{"localhost:50052","localhost:50053","localhost:50054"}
var mergedns [][]string  // se guardan los dominios de los dns para hacer los merges
var mergereg [][][]string
var brokerip = "localhost:50051"
func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func DivideData(data string)[]string{
	return strings.Split(data," ")
}
func DetectDomain(comm string)string{
	str:= DivideData(comm)
	dom:=strings.Split(str[0],".")
	return dom[1]
}
func SearchDomain(dom string)int{
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
func ActReg(pos int){
	fmt.Println("Actualizando Registro ZF dominio: "+dominios[pos])
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
func UpdateFiles(){
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
func Merge(){
	time.Sleep(time.Duration(60)*time.Second)
	// Avisar a los demas dns que se hara un merge, por lo que ellos enviaran los dominios que ellos posean en registro
	for i,dire:= range dns{
		
		if i == idDNS{
			mergedns = append(mergedns, dominios)
			continue
		}

		conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el broker
		if err != nil {
			fmt.Println("Problemas al hacer conexion")
		}
		defer conn.Close()

		client := pb3.NewDNSServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		msg:= &pb3.SendChangesRequest{Soli: "Merge"} 

		resp, err := client.SendChanges(ctx, msg)
		if err != nil {
			fmt.Println("Error, no esta el server conectado ")
		}

		mergedns = append(mergedns, resp.GetDominios())
	}
	for i,_:= range mergedns { //[]dns []dominio []registros
		mergereg=append(mergereg,[][]string{})
		for j,_:= range mergedns[i]{//var mergereg [][][]string .. dns0[[registros],[registros]]
			if i==idDNS{
				mergereg[i]=append(mergereg[i],registro[j])
				continue
			}
			conn, err := grpc.Dial(dns[i], grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				fmt.Println("Problemas al hacer conexion")
			}
			defer conn.Close()
			client := pb3.NewDNSServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			msg:= &pb3.SendChangesRequest{Soli: strconv.Itoa(j)} 
			resp, err := client.SendChanges(ctx, msg)
			if err != nil {
				fmt.Println("Error, no esta el server conectado ")
			}
			mergereg[i]=append(mergereg[i],resp.GetDominios())			
		}
	}
	RealMerge()
	mergedns=[][]string{}
	mergereg=[][][]string{}
	CleanLocalMerge()
	for pag,i:=range pags{
		for index,j:= range i{
			for k,dire:=range dns{
				if k == idDNS{
					continue
				}
		
				conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el broker
				if err != nil {
					fmt.Println("Problemas al hacer conexion")
				}
				defer conn.Close()
		
				client := pb3.NewDNSServiceClient(conn)
		
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				tipo:=int64(1)
				if pag== 0 && index == 0{
					tipo=0
				}
				if pag==len(pags)-1 && index==len(i)-1{
					tipo=2
				}
				msg:= &pb3.ReceiveChangesRequest{Operations: []string{"append",j[0]+" "+j[1]}, Type: tipo} 
				_, err = client.ReceiveChanges(ctx, msg)
				if err != nil {
					fmt.Println("Error, no esta el server conectado ")
				}
			}
		}
	}
	for _,i:=range clocks{
		for k,dire:=range dns{
			if k == idDNS{
				continue
			}
	
			conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				fmt.Println("Problemas al hacer conexion")
			}
			defer conn.Close()
	
			client := pb3.NewDNSServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			tipo:=int64(3)
			msg:= &pb3.ReceiveChangesRequest{Operations: TranslateClock(i), Type: tipo} 
			_, err = client.ReceiveChanges(ctx, msg)
			if err != nil {
				fmt.Println("Error, no esta el server conectado ")
			}
		}
	}
	// Notificamos a Broker que se hizo un merge, para que pueda designar nuevamente ip's nuevas a cada admin
	// Ya que como todos los dns tendran las mismos archivos, no necesariamente para mantener el Read your Writes debe
	// solamente un admin quedarse fijo por siempre a un DNS
	conn, err := grpc.Dial(brokerip, grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				fmt.Println("Problemas al hacer conexion")
			}
			defer conn.Close()
	
			client := pb3.NewDNSServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			msg:= &pb3.NotifyBrokerRequest{Notify: "Merge realizado"} 
			_, err = client.NotifyBroker(ctx, msg)
			if err != nil {
				fmt.Println("Error, broker no esta conectado ")
}
}

func TranslateClock (clk []int64)[]string{
	a:=[]string{}
	for _,j:= range clk{
		b:=strconv.FormatInt(j,10)
		a=append(a,b)
	}
	return a
}
func RealMerge(){
	set0:= make(map[string]bool)//arreglo que guarda direcciones reservadas por prioridad del dns0
	set1:= make(map[string]bool) //arreglo que guarda direcciones reservadas por prioridad del dns1
	for _,i:=range mergereg[0]{ //se guardan las direcciones editadas por dns0
		for _,j:= range i{
			value:=DivideData(j)
			set0[value[1]]=true
		}
	}
	//Se hara revision del dns 1
	for _,i:=range mergereg[1]{
		for _,j:= range i{
			values:=DivideData(j)
			if set0[values[1]]{
				domain:=DetectDomain(values[1])
				pos:=SearchDomain(domain)
				clocks[pos][1]++
				continue
			}
			set1[values[1]]=true
			ReceiveOp(DetectCommand(j),1)
		}
	}
	for _,i:=range mergereg[2]{
		for _,j:= range i{
			values:=DivideData(j)
			if set0[values[1]]||set1[values[1]]{
				domain:=DetectDomain(values[1])
				pos:=SearchDomain(domain)
				clocks[pos][2]++
				continue
			}
			ReceiveOp(DetectCommand(j),2)
		}
	}
}
func CleanLocalMerge(){
	registro= [][]string{}
	for _,_= range dominios{
		registro=append(registro,[]string{})
	}
	for _,domain:= range dominios{
		os.Remove("RegistroZF"+domain+".txt")
		os.Remove("LOG"+domain+".txt")
	}
	UpdateFiles()
}
func ReceiveOp(op []string, dnsid int)(){ //operacion,valores
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
	clocks[pos][dnsid]++
	fmt.Println(clocks)
	ActReg(pos)
}
func (s *server) GetIPBroker(ctx context.Context, msg *pb3.GetIPBrokerRequest) (*pb3.GetIPBrokerResponse, error) {
	pag:=msg.GetDireccion()
	pos:=SearchDomain(DetectDomain(pag)) //por si no encuentra la pag
	for _,j:=range pags[pos]{
		if j[0]==pag{
			return &pb3.GetIPBrokerResponse{Clock: clocks[pos],Ip: j[1]}  , nil
		}
	}
	return &pb3.GetIPBrokerResponse{Clock: []int64{-1},Ip:""}  , nil
}

func (s *server) DnsCommand(ctx context.Context, msg *pb2.DnsCommandRequest) (*pb2.DnsCommandResponse, error) {
	
	fmt.Println("llego: ", msg.GetCommand() )
	fmt.Println( msg.GetCommand()[0] )
	fmt.Println( msg.GetCommand()[1] )
	fmt.Println( len(msg.GetCommand()))
	ReceiveOp(msg.GetCommand(),0)
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

func (s *server) ReceiveChanges(ctx context.Context, msg *pb3.ReceiveChangesRequest) (*pb3.ReceiveChangesResponse, error) {
	
	return &pb3.ReceiveChangesResponse{Status: "listo"}  , nil
	
}
func (s *server) NotifyBroker(ctx context.Context, msg *pb3.NotifyBrokerRequest) (*pb3.NotifyBrokerResponse, error) {
	
	return &pb3.NotifyBrokerResponse{Resp: "Done!"}  , nil
	
}


func main() {
	ReceiveOp([]string{"append","google.cl 193.168.1.3"},0)
	ReceiveOp([]string{"update","google.cl 1.2.3.4"},0)
	ReceiveOp([]string{"update","google.cl gooogle.cl"},0)
	ReceiveOp([]string{"update","google.cl 1.2.3.5"},0)
	ReceiveOp([]string{"delete","google.es"},0)
	go func(){
		for{
			Merge()
		}
	}()
	lis, err := net.Listen("tcp", ":50052")
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