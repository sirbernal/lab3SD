package main

import (

	"fmt"
	"strings"
	"math/rand"
	"time"
)
var idadm []int64
var dns = []string{"localhost:50052","localhost:50053","localhost:50054"}
var randomizer []int

func DetectCommand(comm string)[]string{
	str:= strings.Split(comm, " ")
	var resp []string
	resp=append(resp,strings.ToLower(str[0]))
	resp=append(resp,strings.Join(str[1:]," "))
	return resp
}
func DetectDomain(comm string)string{
	str:= strings.Split(comm, " ")
	dom:=strings.Split(str[0],".")
	return dom[1]
}
func RemoveIndex(s [][][]string, index int) [][][]string {
	return append(s[:index], s[index+1:]...)
}
func upRandomizer(){ //actualiza la designacion al azar para los tres admins que se incorporen
	rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
	randomizer=rand.Perm(3)//genera una permutacion al azar de 0 a 2
}
func giveDNS(admin int64)string{//funcion que designa el dns que le corresponde al admin
	designio:= admin%3 //verifica que designio al azar le corresponde
	dnsadm:=dns[randomizer[designio]] //guarda la direccion que requiere el admin
	idadm=append(idadm,int64(randomizer[designio])) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
	return dnsadm//retorna la direccion al admin
}

func main() {
	upRandomizer()
	for i:=0;i<15;i++{
		fmt.Println(giveDNS(int64(i)))
	}
	fmt.Println()
	for _,j:=range idadm{
		fmt.Println(dns[j])
	} 
}