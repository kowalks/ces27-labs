package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

// global variables
var err string
var myClock int = 1                         // clock
var sharedResourcePort string = ":10001" // shared resource port
var myPort string                        // my server port
var myId string                          // my process id
var nProcess int                         // number of other process servers
var nServers int                         // number of other process servers + shared resource server
var ClientConn []*net.UDPConn            // array of connections with processes
var ServerConn *net.UDPConn              // shared Resource connection

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func doServerJob() {

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

func doClientJob(server int) {
	// for {

	// 	// reading from keyboard
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("Enter a message: ")
	// 	msg, err := reader.ReadString('\n')
	// 	CheckError(err)

	// 	// sending the typed message to 'server'
	// 	buf := []byte(msg)
	// 	_, err = ClientConn[server].Write(buf)
	// 	CheckError(err)

	// 	if err != nil {
	// 		fmt.Println(msg, err)
	// 	}
	// }
}

func beClientOfSharedResource() {
	for {

		// reading from keyboard
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a message: ")
		msg, err := reader.ReadString('\n')
		msg = msg + " [from ID: " + myId + "; CLOCK: " + strconv.Itoa(myClock) + "]"
		CheckError(err)

		// sending the typed message to 'server'
		buf := []byte(msg)
		_, err = ClientConn[nServers-1].Write(buf)
		CheckError(err)

		if err != nil {
			fmt.Println(msg, err)
		}
	}
}

func initConnections() {

	// STARTING CONNECTIONS WITH OTHER PROCESSES

	// getting the second argument (id) from './Process $id :myport :otherport :otherport :otherport ...'
	myId = os.Args[1]

	// getting the process port
	myIntId, err := strconv.Atoi(myId)
	CheckError(err)
	portPosition := myIntId + 1
	myPort = os.Args[portPosition]

	nProcess = len(os.Args) - 3
	nServers = nProcess + 1
	ClientConn = make([]*net.UDPConn, nServers)

	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+myPort)
	CheckError(err)
	ServerConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	for processServer := 0; processServer < nProcess; processServer++ {

		ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+os.Args[3+processServer])
		CheckError(err)

		Conn, err := net.DialUDP("udp", nil, ServerAddr)
		ClientConn[processServer] = Conn
		CheckError(err)
	}

	// START CONNECTION WITH SHARED RESOURCE

	ServerAddr, err = net.ResolveUDPAddr("udp", "127.0.0.1"+sharedResourcePort)
	CheckError(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	ClientConn[nServers-1] = Conn
	CheckError(err)
}

func main() {

	initConnections()

	defer ServerConn.Close()
	for i := 0; i < nServers; i++ {
		defer ClientConn[i].Close()
	}

	// listening messages from other processes
	go doServerJob()

	for j := 0; j < nServers-1; j++ {
		// sending messages to other processes
		go doClientJob(j)
	}

	// sending messages shared resource
	go beClientOfSharedResource()

	// infinite loop
	for {
	}
}
