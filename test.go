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
var id int64

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
func giveDNS(){//funcion que designa el dns cada vez que se agregue un admin
	designio:= id%3 //verifica que designio al azar le corresponde
	 //guarda la direccion que requiere el admin
	idadm=append(idadm,int64(randomizer[designio])) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
}
func giveDNSReset(idn int){//funcion que designa el dns cada vez que se agregue un admin
	designio:= idn%3 //verifica que designio al azar le corresponde
	 //guarda la direccion que requiere el admin
	idadm[idn]=int64(randomizer[designio]) //se guarda una lista opcional que registra donde se conectará cada admin
	if designio==2{//si es el ultimo de los 3 admins por grupo se reinicia el arreglo que designa al azar
		upRandomizer()
	}
}

func resetDNS(){
	upRandomizer()
	for i:=0; i<int(id);i++{
		giveDNSReset(i)
	}
}
func DetectUpdate(w string)bool{ //funcion que detecta si es cambio de ip o cambio de dominio
	str:= strings.Split(w, ".")
	if len(str)==4{ //se asume que la ip vendrá en formato x.x.x.x
		return false //es cambio de dominio
	}else{
		return true //es cambio de ip
	}
}

func main() {
	fmt.Println(DetectUpdate("asdsd"))
	fmt.Println(DetectUpdate("1.1.1.1"))
	fmt.Println(DetectUpdate("sad.c"))
}