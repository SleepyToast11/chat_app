package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	ServerHost = "localhost"
	ServerPort = "8086"
)

func makeServerConnection() (net.Conn, error) {
	listener, err := net.Listen("tcp", ServerHost+":"+ServerPort) //blocking
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return nil, err
	}

	fmt.Println("Server listening on", ServerHost+":"+ServerPort)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}
	return conn, nil
}

func makeClientConnection() (net.Conn, error) {
	address := fmt.Sprintf("%s:%s", ServerHost, 8080)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("failed to create TCP client")
	}
	return tcpConn, nil
}

func main() {

	var conn net.Conn
	var err error
	if len(os.Args) < 1 {
		println("serverMode")
		conn, err = makeClientConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		println("client mode")
		conn, err = makeServerConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	println("Connected to server")

	killChan := make(chan bool)
	strChanIn := make(chan string)
	strChanOut := make(chan string)

	go handleTerminal(strChanOut, killChan)
	go handleConnection(conn, killChan, strChanIn)

	var incomeMessages string
	for {
		select {
		case <-killChan:
			os.Exit(0)
		case str := <-strChanIn:
			incomeMessages += str + "\n"
		case str := <-strChanOut:
			if str == "!read" {
				if incomeMessages == "" {

				} else {
					println(incomeMessages)
					incomeMessages = ""
				}
			}
			_, err := conn.Write([]byte(str + "\n"))
			if err != nil {
				println("Error writing to connection:", err)
				return
			}

		}
	}
}

func handleTerminal(strChan chan<- string, killChan chan bool) {
	println("hello!")
	defer func() {
		killChan <- true
		close(strChan) // Close the channel when done reading
		close(killChan)
		return
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {

		case <-killChan:
			return
		default:
			if !scanner.Scan() {
				err := scanner.Err()
				fmt.Println("Error reading input:", err)
				return
			}
			input := scanner.Text()
			if input == "!quit" {
				fmt.Println("Bye Bye")
				return
			}
			strChan <- input
		}
	}

}

func handleConnection(conn net.Conn, kill chan bool, strChan chan<- string) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			println("Error closing connection:", err)
		}
	}(conn)
	for {
		select {
		case <-kill:
			fmt.Println("Received kill signal. Exiting handleConnection.")
			return
		default:
			buffer := make([]byte, 256)
			n, err := conn.Read(buffer)
			strChan <- string(buffer[:n])
			print("|")
			if err != nil {
				fmt.Println("Error reading data:", err)
				kill <- true
				return
			}
		}
	}
}
