package main

import (
	"context"
	"fmt"
	"time"
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

var idDNS = 0 //id del dns correspondiente 
var dominios []string //arreglo que guarda los dominios en formato ["cl","com",...]
var registro [][]string //arreglo que mantiene los registros por dominio[[registros cl],[registros com],...] => [registros cl]=["append a.cl 1.2.3.4", "update u.cl 2.2.2.2"...]
var pags [][][]string //arreglo que guarda las paginas e ips por dominio[[paginas.cl],[paginas.com],...] =>[paginas.cl]=[["a.cl","1.2.3.4"],["u.cl","2.2.2.2"]]
var clocks [][]int64  //relojes asociados por dominio [[cl],[com],...] => [cl]=[x,y,z]
var timeout = time.Duration(1)*time.Second //timeout para conexiones
//var dns = []string{"localhost:50052","localhost:50053","localhost:50054"} 
var dns = []string{"10.10.28.82:50052","10.10.28.83:50053","10.10.28.84:50054"} //ips de los dns
var mergedns [][]string  // auxiliar donde guardan los dominios de los dns para hacer los merges
var mergereg [][][]string //auxiliar donde se guardan los registros de los dns para merge 
//var brokerip = "localhost:50051"
var brokerip = "10.10.28.81:50051" //ip del broker
func DetectCommand(comm string)[]string{ //funcion que separa el comando (append... etc) del resto de datos y lo estandariza en minuscula en un arreglo [comando, contenido]
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func DivideData(data string)[]string{ //funcion que divide un string por espacio(solo acorta la funcion para escribir menos y ahorrar tiempo)
	return strings.Split(data," ")
}
func DetectDomain(comm string)string{ //funcion que detecta el dominio de una pagina ingresada, por ejemplo abc.com => com
	str:= DivideData(comm)
	dom:=strings.Split(str[0],".")
	return dom[1]
}
func SearchDomain(dom string)int{ //funcion que busca la posicion del dominio en el arreglo guardado
	pos:= -1 //retorna -1 para identificar que el dominio aun no ha sido ingresado
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
func ActReg(pos int){ //funcion que actualiza tanto el registro zf del dominio señalado (con el index) como el log de dicho dominio
	file,err:= os.OpenFile("RegistroZF"+dominios[pos]+".txt",os.O_CREATE|os.O_WRONLY,0777) 
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
	file2,err:= os.OpenFile("LOG"+dominios[pos]+".txt",os.O_CREATE|os.O_WRONLY,0777) 
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
func Merge(){ //funcion que realiza el merge
	time.Sleep(time.Duration(300)*time.Second)
	// Avisar a los demas dns que se hara un merge, por lo que ellos enviaran los dominios que ellos posean en registro
	for i,dire:= range dns{
		if i == idDNS{
			mergedns = append(mergedns, dominios)
			continue
		}

		conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el broker
		if err != nil {
			fmt.Println("Merge Error...Problemas al hacer conexion con dns "+strconv.Itoa(i)+"Omitiendo Merge para nodo offline")
			continue
		}
		defer conn.Close()

		client := pb3.NewDNSServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		msg:= &pb3.SendChangesRequest{Soli: "Merge"} 

		resp, err := client.SendChanges(ctx, msg)
		if err != nil {
			fmt.Println("Merge Error...Problemas al hacer conexion con dns "+strconv.Itoa(i)+"Omitiendo Merge para nodo offline")
			continue
		}

		mergedns = append(mergedns, resp.GetDominios())
	}
	//ya con los dominios de todos comienza a solicitar los registros asociados a estos a todos los dominios
	for i,_:= range mergedns { 
		mergereg=append(mergereg,[][]string{})
		for j,_:= range mergedns[i]{
			if i==idDNS{
				mergereg[i]=append(mergereg[i],registro[j])
				continue
			}
			conn, err := grpc.Dial(dns[i], grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				continue
			}
			defer conn.Close()
			client := pb3.NewDNSServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			msg:= &pb3.SendChangesRequest{Soli: strconv.Itoa(j)} 
			resp, err := client.SendChanges(ctx, msg)
			if err != nil {
				continue
			}
			mergereg[i]=append(mergereg[i],resp.GetDominios())			
		}
	}
	RealMerge() //realiza el merge de forma local en el dns 0
	mergedns=[][]string{} //vacía los auxiliares para el procesp
	mergereg=[][][]string{}
	CleanLocalMerge() //realiza la limpieza de archivos y variables que se requieren post-merge
	for pag,i:=range pags{ //Procede a enviar las paginas que quedaron luego de realizado el proceso simulando ser un administrador con "append pagina ip" al resto de dns
		for index,j:= range i{
			for k,dire:=range dns{
				if k == idDNS{
					continue
				}
				conn, err := grpc.Dial(dire, grpc.WithInsecure())
				if err != nil {
					continue
				}
				defer conn.Close()
		
				client := pb3.NewDNSServiceClient(conn)
		
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				tipo:=int64(1)
				if pag== 0 && index == 0{
					tipo=0 //si es la primera pagina avisa para el primer clean a cada dns
				}
				if pag==len(pags)-1 && index==len(i)-1{
					tipo=2 //si es la ultima pagina avisa para el ultimo clean de cada dns
				}
				msg:= &pb3.ReceiveChangesRequest{Operations: []string{"append",j[0]+" "+j[1]}, Type: tipo} 
				_, err = client.ReceiveChanges(ctx, msg)
				if err != nil {
					continue
				}
			}
		}
	}
	for _,i:=range clocks{//procede a enviar los clocks para mantener la consistencia
		for k,dire:=range dns{
			if k == idDNS{
				continue
			}
	
			conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				continue
			}
			defer conn.Close()
	
			client := pb3.NewDNSServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			tipo:=int64(3)
			msg:= &pb3.ReceiveChangesRequest{Operations: TranslateClock(i), Type: tipo} 
			_, err = client.ReceiveChanges(ctx, msg)
			if err != nil {
				continue
			}
		}
	}
	// Notificamos a Broker que se hizo un merge, para que pueda designar nuevamente ip's nuevas a cada admin
	// Ya que como todos los dns tendran las mismos archivos, no necesariamente para mantener el Read your Writes debe
	// solamente un admin quedarse fijo por siempre a un DNS
	conn, err := grpc.Dial(brokerip, grpc.WithInsecure()) //genera la conexion con el broker
			if err != nil {
				fmt.Println("Problemas al hacer conexion con broker para notificar Merge")
			}
			defer conn.Close()
	
			client := pb3.NewDNSServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			msg:= &pb3.NotifyBrokerRequest{Notify: "Merge realizado"} 
			_, err = client.NotifyBroker(ctx, msg)
			if err != nil {
				fmt.Println("Problemas al hacer conexion con broker para notificar Merge")
}
}

func TranslateClock (clk []int64)[]string{ //traduce el clock de int64 a string para merge
	a:=[]string{}
	for _,j:= range clk{
		b:=strconv.FormatInt(j,10)
		a=append(a,b)
	}
	return a
}
func RealMerge(){ //funcion que realiza el merge interno
	set0:= make(map[string]bool)//arreglo que guarda direcciones reservadas por prioridad del dns0
	set1:= make(map[string]bool) //arreglo que guarda direcciones reservadas por prioridad del dns1
	for _,i:=range mergereg[0]{ //se guardan las direcciones editadas por dns0
		for _,j:= range i{
			value:=DivideData(j)
			set0[value[1]]=true
		}
	}
	//Se hace revision del dns 1
	for _,i:=range mergereg[1]{
		for _,j:= range i{
			values:=DivideData(j)
			if set0[values[1]]{//si es una direccion reservada en dns0 solo actualiza el clock
				domain:=DetectDomain(values[1])
				pos:=SearchDomain(domain)
				clocks[pos][1]++
				continue
			}
			set1[values[1]]=true //reserva la pagina
			ReceiveOp(DetectCommand(j),1) //realiza la accion en dns0
		}
	}
	for _,i:=range mergereg[2]{
		for _,j:= range i{
			values:=DivideData(j)
			if set0[values[1]]||set1[values[1]]{//si es una direccion reservada en dns0 o dns1 solo actualiza el clock
				domain:=DetectDomain(values[1])
				pos:=SearchDomain(domain)
				clocks[pos][2]++
				continue
			}
			ReceiveOp(DetectCommand(j),2) //realiza la accion en dns0
		}
	}
}
func CleanLocalMerge(){//realiza la limpieza de variables y archivos post merge y actualiza los archivos
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
func ReceiveOp(op []string, dnsid int)(){ //funcion que recibe el comando de admin (id 0) u otro dns en el merge (id del dns de origen)
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
	ActReg(pos)
}
func (s *server) GetIPBroker(ctx context.Context, msg *pb3.GetIPBrokerRequest) (*pb3.GetIPBrokerResponse, error) { //retorna la ip de la pagina solitada por el cliente al broker
	pag:=msg.GetDireccion()
	pos:=SearchDomain(DetectDomain(pag)) //por si no encuentra la pag
	if pos==-1{
		return &pb3.GetIPBrokerResponse{Clock: []int64{-1},Ip:""}  , nil
	}
	for _,j:=range pags[pos]{
		if j[0]==pag{
			return &pb3.GetIPBrokerResponse{Clock: clocks[pos],Ip: j[1]}  , nil
		}
	}
	return &pb3.GetIPBrokerResponse{Clock: []int64{-1},Ip:""}  , nil //retorna nulo si no lo encuentra
}

func (s *server) DnsCommand(ctx context.Context, msg *pb2.DnsCommandRequest) (*pb2.DnsCommandResponse, error) {//funcion que recibe el comando del admin
	ReceiveOp(msg.GetCommand(),0) 
	return &pb2.DnsCommandResponse{Clock: []int64{} }, nil
}

func (s *server) Broker(ctx context.Context, msg *pb2.BrokerRequest) (*pb2.BrokerResponse, error) {
	return &pb2.BrokerResponse{Ip: "192.168.0.1"}, nil
}
func (s *server) RegAdm(ctx context.Context, msg *pb2.RegAdmRequest) (*pb2.RegAdmResponse, error) {
	return &pb2.RegAdmResponse{Id: 0 }, nil
}

func (s *server) SendChanges(ctx context.Context, msg *pb3.SendChangesRequest) (*pb3.SendChangesResponse, error) { //canal de envio usado en merge
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
	go func(){
		for{
			Merge()
		}
	}()
	fmt.Println("DNS 0 en línea")
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