package main

import (
	"fmt"
	"net"
	"os"
)

var myPort string = ":10001" // Shared Resource port
var ServerConn *net.UDPConn  // Shared Resource connection

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func main() {
	// preparing a address to listen at Shared Resource port
	Address, err := net.ResolveUDPAddr("udp", myPort)
	CheckError(err)

	// start listening
	ServerConn, err := net.ListenUDP("udp", Address)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}

		//Loop infinito para receber mensagem e escrever todo
		//conteúdo (processo que enviou, relógio recebido e texto)
		//na tela
		//FALTA FAZER
	}
}
