package main

import (
	"fmt"
	"lot/broker"

)


func main() {
    client:=broker.ConnectDevice()
    if client!=nil{
	fmt.Println("build client success")
    }else{
	fmt.Println("build client fail")
    }
	select{}
}
