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



/*func Exists(name string) bool { //función sacada de https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}*/
/*func InitReg(){
	if !Exists("./RegistroZF.txt"){
		fmt.Println("Archivo: 'RegistroZF' no detectado... generando nuevo Registro")
		ActReg()
	}else{
		fmt.Println("Archivo: 'RegistroZF' detectado... actualizando memoria")
		file, err := os.Open("./RegistroZF.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} 
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line:=DetectCommand(scanner.Text())
			regmem = append(regmem, []string{line[0],line[1]})
		}
		file.Close()
	}
}*/
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
func Merge(){
	time.Sleep(time.Duration(15)*time.Second)
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
	fmt.Println(clocks)
	fmt.Println(pags)
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
		for i,j :=range pags[pos]{
			if j[0]==values[0]{
				pags[pos][i]=values
				fmt.Println(pags[pos][i][1])
				os.Remove("RegistroZF"+domain+".txt")
				break
			}
		}
	case "delete":
		for i,j :=range pags[pos]{
			if j[0]==values[0]{
				pags[pos]=RemoveIndex(pags[pos],i)
				os.Remove("RegistroZF"+domain+".txt")
				break
			}
		}
	}
	registro[pos]=append(registro[pos],op[0]+" "+op[1])
	clocks[pos][dnsid]++
	fmt.Println(clocks)
	ActReg(pos)
}
func (s *server) GetClock(ctx context.Context, msg *pb3.GetClockRequest) (*pb3.GetClockResponse, error) {
	return &pb3.GetClockResponse{Clock: []int64{}}  , nil
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


func main() {
	/*ReceiveOp([]string{"append","google.cl aquiIP"},0)
	ReceiveOp([]string{"append","google.com Ipqlia"},0)
	ReceiveOp([]string{"append","asd.cl asdhj"},0)
	ReceiveOp([]string{"append","lel.zz sadkjasdh"},0)
	ReceiveOp([]string{"delete","google.cl"},0)
	ReceiveOp([]string{"update","google.com nueva Ip"},0)
	ReceiveOp([]string{"append","lul.zz ñaña"},0)
	ReceiveOp([]string{"update","lel.zz holi"},0)
	ReceiveOp([]string{"update","lel.zz asd"},0)*/
	Merge()

	fmt.Println(registro)
	fmt.Println(dominios)
	fmt.Println(clocks)
	fmt.Println(pags)
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