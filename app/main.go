package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "127.0.0.1:"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func main() {
	var db = make(map[string]string)

	fmt.Println("Server is Open on: " + CONN_HOST + CONN_PORT)

	// Listen for incoming connections.
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+CONN_PORT)

	defer listener.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close the listener when the application closes.
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine for multiples acess.
		go handleRequest(conn, db)
	}

}

func handleRequest(conn net.Conn, db map[string]string) {
	// Close the connection when you're done with it.
	conn.Write([]byte("+PONG\r\n"))

	// Receive message
	for {
		r := bufio.NewReader(conn)

		n, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Command received: ", n)

		comand, err := converToComand(n)

		if err != nil {
			conn.Write([]byte(err.Error() + "\r\n"))
			continue
		}

		result := comandReader(comand, db, conn)

		if len(n) > 0 {
			conn.Write([]byte(result))
		}
	}
}

func converToComand(command string) ([]string, error) {
	// Parse command removin separetor "/" to array
	arrCommand := strings.Split(command, "/")

	fmt.Println("arrCommand: ", arrCommand)

	// Valid command
	if arrCommand[0] == "SET" && len(arrCommand) < 2 {
		return nil, fmt.Errorf("Invalid parameters to command SET")
	}

	if arrCommand[0] == "GET" && len(arrCommand) < 1 {
		return nil, fmt.Errorf("Invalid parameters to command GET")
	}

	return arrCommand, nil
}

func comandReader(command []string, db map[string]string, conn net.Conn) string {
	// Switch to command
	// TODO: add more commands

	switch command[0] {
	case "PING":
		return "PONG"
	case "SET":
		db[command[1]] = command[2]
		fmt.Println("db: ", db)
		return "OK\r\n"
	case "DEL":
		delete(db, command[1])
		return "OK\r\n"
	case "QUIT":
		conn.Close()
		return "BYE\r\n"
	case "GET":
		if val, ok := db[command[1]]; ok {
			return val + "\r\n"
		}
		return "ERROR TO FIND\r\n"
	}

	return "COMAND ERROR\r\n"
}
