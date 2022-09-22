package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

// global variables
const (
	RELEASED = 0
	WANTED   = 1
	HELD     = 2
)
var err string
var myClock int = 1                      // clock
var sharedResourcePort string = ":10001" // shared resource port
var myPort string                        // my server port
var myId string                          // my process id
var nProcess int                         // number of other process servers
var nServers int                         // number of other process servers + shared resource server
var myState int = RELEASED               // state of process
var ClientConn []*net.UDPConn            // array of connections with processes
var ServerConn *net.UDPConn              // shared Resource connection

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func incrementClock() {
	myClock = myClock + 1
}

func didProcessAlreadyRequestCS() bool {
	if myState != RELEASED {
		return true
	}

	return false
}

func doServerJob() {

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		msg := string(buf[0:n])
		fmt.Println("Received ", msg, " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}

		if msg == "x\n" {

			isReceiving := true
			ricartAgrawala(msg, isReceiving, 1)

		}
	}
}

func sendMessageToOtherProcesses(processId int, msg string) {

	buf := []byte(msg)
	_, err := ClientConn[processId].Write(buf)
	CheckError(err)

	if err != nil {
		fmt.Println(msg, err)
	}
}

func doClientJob() {
	for {

		// reading from keyboard
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a message: ")
		msg, err := reader.ReadString('\n')
		CheckError(err)

		if msg == "x\n" {

			if didProcessAlreadyRequestCS() == false {
				myState = WANTED
				isReceiving := false // because is sending
				myIntId, _ := strconv.Atoi(myId)
				ricartAgrawala(msg, isReceiving, myIntId)
			} else {
				fmt.Println("x ignorado")
			}

		} else if msg == myId + "\n" {

			incrementClock()

		} else {

			msg = msg + " [from ID: " + myId + "; CLOCK: " + strconv.Itoa(myClock) + "]"
	
			// sending the typed message to Shared Resource
			buf := []byte(msg)
			_, err = ClientConn[nServers-1].Write(buf)
			CheckError(err)
	
			if err != nil {
				fmt.Println(msg, err)
			}

		}

	}
}

func ricartAgrawala(msg string, isReceiving bool, processId int) {


	if isReceiving == true {

		if myState == RELEASED {
			sendMessageToOtherProcesses(processId-1, "pode entrar")
		}

	} else {
		// sending messages to other processes to request CS
		for j := 0; j < nServers-1; j++ {
			if strconv.Itoa(j+1) != myId {
				sendMessageToOtherProcesses(j, msg)
			}
		}

	}

}

func initConnections() {

	// STARTING CONNECTIONS WITH OTHER PROCESSES

	// getting the second argument (id) from './Process $id :port :port :port ...'
	myId = os.Args[1]

	// port indexes offset
	offset := 2

	// getting the process port
	myIntId, err := strconv.Atoi(myId)
	CheckError(err)
	portPosition := myIntId + 1
	myPort = os.Args[portPosition]

	nProcess = len(os.Args) - offset
	nServers = nProcess + 1
	ClientConn = make([]*net.UDPConn, nServers)

	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+myPort)
	CheckError(err)
	ServerConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	for processServer := 0; processServer < nProcess; processServer++ {

		//fmt.Println("add: ", "127.0.0.1"+os.Args[offset+processServer])

		ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+os.Args[offset+processServer])
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

	// sending message to another process and to SharedResource
	go doClientJob()

	// infinite loop
	for {
	}
}
