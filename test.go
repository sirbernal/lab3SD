package main

import (

	"fmt"
	"time"
)

func changer(){
	for{
		time.Sleep(time.Duration(3)*time.Second)
		variable[0] ++ 
		variable[1] ++ 
		variable[2] ++ 
		fmt.Println(variable)
	}
	
}

var variable = []int{1,2,3}

func main() {
	go changer()
	time.Sleep(time.Duration(3)*time.Second)
	go changer()
	time.Sleep(time.Duration(10000)*time.Second)

}