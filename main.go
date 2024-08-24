package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	ServerHost = "localhost"
	ServerPort = "8090"
)

func main() {
	killChan := make(chan bool)
	strChanIn := make(chan string)
	strChanOut := make(chan string)

	defer close(killChan)
	defer close(strChanOut)
	defer close(killChan)

	go handleTerminal(strChanOut, killChan)
	conn, err := connectNonBlocking(killChan)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println("Connected to server")

	go handleConnection(conn, killChan, strChanIn)

	if err := controller(killChan, strChanIn, strChanOut, conn); err != nil {
		fmt.Println(err)
	}
}

func connectNonBlocking(killChan <-chan bool) (net.Conn, error) {
	//allow to kill program before initialization
	for {
		select {
		case <-killChan:
			os.Exit(0)
		default:
			if len(os.Args) < 2 {
				println("client Mode")
				conn, err := MakeClientConnection(ServerHost, ServerPort)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return conn, nil
			} else {
				println("server mode")
				conn, err := MakeServerConnection(ServerHost, ServerPort)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return conn, nil
			}
		}
	}

}

func controller(killChan chan bool, strChanIn chan string, strChanOut chan string, conn net.Conn) error {
	var incomeMessages string
	for {
		select {
		case <-killChan:
			return nil
		case str := <-strChanIn:
			incomeMessages += str + "\n"
		case str := <-strChanOut:
			if str == "!read" {
				if incomeMessages == "" {
				} else {
					println(incomeMessages)
					incomeMessages = ""
				}
			} else {
				_, err := conn.Write([]byte(str + "\n"))
				if err != nil {
					println("Error writing to connection:", err)
					return err
				}
			}
		}
	}
}

func handleTerminal(strChan chan<- string, killChan chan bool) {
	defer func() {
		killChan <- true
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
