package main

import (
	"fmt"
	"net"
)

func MakeServerConnection(host string, port string) (net.Conn, error) {
	listener, err := net.Listen("tcp", host+":"+port) //blocking don't forget
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return nil, err
	}
	fmt.Println("Server listening on", host+":"+host)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}
	return conn, nil
}

func MakeClientConnection(host string, port string) (net.Conn, error) {
	fmt.Println("Trying to reach", host+":"+port)
	address := fmt.Sprintf("%s:%s", host, port)
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
