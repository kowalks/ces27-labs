package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// global variables
const (
	RELEASED = 0
	WANTED   = 1
	HELD     = 2
)
var err string
var myClock int = 0                      // clock
var sharedResourcePort string = ":10001" // shared resource port
var myPort string                        // my server port
var myId string                          // my process id
var nProcess int                         // number of other process servers
var nServers int                         // number of other process servers + shared resource server
var myState int = RELEASED               // state of process
var ClientConn []*net.UDPConn            // array of connections with processes
var ServerConn *net.UDPConn              // shared Resource connection
var repliesCounter int = 0               // if repliesCounter is equal to nProcess -1 so process held the CS
var queueId []int                        // queue of ids that wait to CS
var queueTime []int                      // queue of clock of ids that wait to CS
var mutex sync.Mutex                     // mutex

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func incrementClock() {
	mutex.Lock()
	myClock = myClock + 1
	fmt.Println("  clock:", myClock)
	mutex.Unlock()
}

func checkClocksAndIncrement(clock1 int, clock2 int) {
	mutex.Lock()

	if clock1 > clock2 {
		myClock = clock1 + 1
	} else {
		myClock = clock2 + 1
	}
	fmt.Println("  clock:", myClock)

	mutex.Unlock()
}

func didProcessAlreadyRequestCS() bool {
	if myState != RELEASED {
		return true
	}

	return false
}

func enterOnCS() {
	fmt.Println("\nEntering the CS...")
	sendMessageToSharedResource("Entrei na CS")

	// sleeping and wake up after 10 seconds
	time.Sleep(10 * time.Second)

	myState = RELEASED
	
	sendMessageToSharedResource("Saí da CS")
	fmt.Println("\nExiting the CS...")

	// send message to process in queue
	queueLength := len(queueId)
	for p := 0; p < queueLength; p++ {
		id := queueId[p] - 1
		//clock := queueTime[p]
		sendMessageToAnotherServer(id, "OK" + ",ID:" + myId + ",CLOCK:" + strconv.Itoa(myClock))
	}

	// clean arrays
	queueId = []int{}
	queueTime = []int{}

	incrementClock()
}

func doServerJob() {

	buf := make([]byte, 1024)

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		entireMessage := string(buf[0:n])
		entireMessage = strings.Trim(entireMessage, "\n")
		msg := strings.Split(entireMessage, ",")[0]
		id := strings.Split(strings.Split(entireMessage, ":")[1], ",")[0]
		clock := strings.Split(entireMessage, ":")[2]
		fmt.Println("\n\nReceived", msg, "[from ID: " + id + "; CLOCK: " + clock + "]")

		if err != nil {
			fmt.Println("Error: ", err)
		}

		if msg == "x" {

			// time.Sleep(5 * time.Second)

			intId, _ := strconv.Atoi(id)
			intClock, _ := strconv.Atoi(clock)

			checkClocksAndIncrement(myClock, intClock)

			isReceiving := true
			ricartAgrawala(msg, isReceiving, intId, intClock)

		} else if msg == "OK" {

			intId, _ := strconv.Atoi(id)       // will be not used
			intClock, _ := strconv.Atoi(clock)

			checkClocksAndIncrement(myClock, intClock)

			isReceiving := true
			ricartAgrawala(msg, isReceiving, intId, intClock)

		}
	}
}

func sendMessageToAnotherServer(processId int, msg string) {

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
		msg = strings.Trim(msg, "\n");
		CheckError(err)

		if msg == "x" {

			incrementClock()

			// time.Sleep(5 * time.Second)

			if didProcessAlreadyRequestCS() == false {
				myState = WANTED
				isReceiving := false // because is sending
				myIntId, _ := strconv.Atoi(myId)
				ricartAgrawala(msg + ",ID:" + myId + ",CLOCK:" + strconv.Itoa(myClock), isReceiving, myIntId, myClock)
			} else {
				fmt.Println("x ignorado")
			}

		} else if msg == myId {

			incrementClock()

		} else {

			sendMessageToSharedResource(msg)
			incrementClock()

		}

	}
}

func ricartAgrawala(msg string, isReceiving bool, processId int, processClock int) {

	if isReceiving == true { // receiving message from another process

		if myState == RELEASED {

			sendMessageToAnotherServer(processId-1, "OK" + ",ID:" + myId + ",CLOCK:" + strconv.Itoa(myClock))

		} else if myState == WANTED && msg == "x" {

			myIntId, _ := strconv.Atoi(myId)

			if myClock > processClock && myIntId < processId {
				// fmt.Println("*** current clock", myClock, "msg clock", processClock)
				queueId = append(queueId, processId)
				queueTime = append(queueTime, processClock)

			} else {

				sendMessageToAnotherServer(processId-1, "OK" + ",ID:" + myId + ",CLOCK:" + strconv.Itoa(myClock))

			}


		} else if myState == WANTED && msg == "OK" {

			repliesCounter = repliesCounter + 1
			if repliesCounter == 2 {
				repliesCounter = 0
				myState = HELD
				go enterOnCS()
			}

		} else if myState == HELD {

			queueId = append(queueId, processId)
			queueTime = append(queueTime, processClock)

		}

	} else { // sending messages to other processes to request CS
		
		for j := 0; j < nServers-1; j++ {
			if strconv.Itoa(j+1) != myId {
				sendMessageToAnotherServer(j, msg)
			}
		}

	}

}

func sendMessageToSharedResource(msg string) {
	sendMessageToAnotherServer(nServers-1, msg + " [from ID: " + myId + "; CLOCK: " + strconv.Itoa(myClock) + "]")
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
