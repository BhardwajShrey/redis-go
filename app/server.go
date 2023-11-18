package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
    PONG = "+PONG\r\n"
    INVALID_COMMAND = "+INVALID_COMMAND\r\n"
)

func createBulkString(reqArray []string) string {
    return strings.Join(reqArray, "\r\n")
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
    }

    fmt.Println("Listening on port 6379...")

    for {
        connection, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
            os.Exit(1)
        }

        go handleConnection(connection)
    }
}

func handleConnection(connection net.Conn) {
    defer connection.Close()

    buffer := make([]byte, 1024)
    _, err := connection.Read(buffer)
    if err != nil {
        fmt.Println("Error reading buffer from connection: ", err.Error())
        os.Exit(1)
    }

    reqArray := strings.Split(string(buffer), "\r\n")

    fmt.Printf("Request received: %v\n", reqArray)

    var res string
    switch reqArray[2] {
    case "ping":
        res = PONG
    case "ECHO":
        res = createBulkString(reqArray[3 :])
    default:
        res = INVALID_COMMAND
    }

    _, err = connection.Write([]byte(res))
    if err != nil {
        fmt.Println("Error writing response to connection", err.Error())
        os.Exit(1)
    }
}
