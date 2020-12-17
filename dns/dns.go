package main

import (
	"context"
	"fmt"
	//"time"
	"strings"
	"bufio"
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
var regmem [][]string
func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func Exists(name string) bool { //función sacada de https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}
func InitReg(){
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
}
func ActReg(){
	fmt.Println("Actualizando Registro")
	file,err:= os.OpenFile("RegistroZF.txt",os.O_CREATE|os.O_WRONLY,0777) //abre o genera el archivo de registro
	defer file.Close()
	if err !=nil{
		os.Exit(1)
	}
	for _,j:= range regmem{
		word:= j[0]+" "+j[1]
		_, err := file.WriteString(word + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	file.Close()
}

func (s *server) GetClock(ctx context.Context, msg *pb3.GetClockRequest) (*pb3.GetClockResponse, error) {
	return &pb3.GetClockResponse{Clock: clock}  , nil
}


func main() {
	InitReg()
	fmt.Println(regmem)
	regmem=append(regmem,[]string{"asd","1"})
	regmem=append(regmem,[]string{"asd","2"})
	regmem=append(regmem,[]string{"asd","31"})
	regmem=append(regmem,[]string{"asd","4"})
	ActReg()
	fmt.Println(regmem)


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