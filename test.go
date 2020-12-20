package main

import (

	"fmt"
	"strings"
	//"time"
)

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

func main() {
	var asd [][][]string
	asd=append(asd, [][]string{})
	asd=append(asd, [][]string{})
	asd=append(asd, [][]string{})
	fmt.Println(len(asd))
	asd=RemoveIndex(asd,1)
	fmt.Println(len(asd))
}