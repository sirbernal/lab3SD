package main

import (
	"context"
	"fmt"
	//"time"
	"strings"
//	"bufio"
	"os"
	"log"
	"net"

	//pb "github.com/sirbernal/lab3SD/proto/client_service"
	pb3 "github.com/sirbernal/lab3SD/proto/dns_service"
	"google.golang.org/grpc"

)

type server struct {
}

var clock = []int64{0,0,0}

var dominios []string //[com,cl,ez]
var registro [][]string //[[agregar,borrar,wea],[],[],]
var pags [][][]string //[[[algo.com,direccion],[xd.com, direccion]],[lel.cl],[]]
var clocks [][]int64 

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


func ReceiveOp(op []string)(){ //operacion,valores
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
	ActReg(pos)
}
func (s *server) GetClock(ctx context.Context, msg *pb3.GetClockRequest) (*pb3.GetClockResponse, error) {
	return &pb3.GetClockResponse{Clock: clock}  , nil
}


func main() {
	ReceiveOp([]string{"append","google.cl aquiIP"})
	ReceiveOp([]string{"append","google.com Ipqlia"})
	ReceiveOp([]string{"append","asd.cl asdhj"})
	ReceiveOp([]string{"append","lel.zz sadkjasdh"})
	ReceiveOp([]string{"delete","google.cl"})
	ReceiveOp([]string{"update","google.com nueva Ip"})
	ReceiveOp([]string{"append","lul.zz ñaña"})
	ReceiveOp([]string{"update","lel.zz holi"})
	ReceiveOp([]string{"update","lel.zz asd"})
	fmt.Println(pags)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()

	pb3.RegisterDNSServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}