package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type BroadcastCommand struct {
	sourceUsername string
	command        string
}

func main() {
	broadcast := make(chan BroadcastCommand)
	listener, err := net.Listen("tcp", ":3117")
	defer listener.Close()

	if err != nil {
		fmt.Println(fmt.Errorf("%f", err))
		os.Exit(1)
	}

	for {
		fmt.Println("Waiting for a connection...")
		conn, err := listener.Accept()
		fmt.Println("\tConnection made...")

		if err != nil {
			fmt.Println(fmt.Errorf("Error after accepting connection :: %f", err))
			break
		}

		go handleConnection(conn, broadcast)
	}
}

func handleConnection(conn net.Conn, broadcast chan BroadcastCommand) {
	var buffer = make([]byte, 1024)
	var sockReader = bufio.NewReader(conn)
	var sockWriter = bufio.NewWriter(conn)

ConnectionLoop:
	for {
		length, err := sockReader.Read(buffer)
		payload := string(buffer[:length])

		if err == io.EOF {
			break ConnectionLoop
		}

		if strings.HasPrefix(payload, ":login") {
			cmdParts := strings.Fields(payload)
			username := cmdParts[1]
			go playGame(username, sockReader, sockWriter, conn, broadcast)
			break ConnectionLoop
		}
	}
}

func playGame(username string, sockReader *bufio.Reader, sockWriter *bufio.Writer, conn net.Conn, broadcast chan BroadcastCommand) {
	defer func() {
		fmt.Println("Closing game...")
		conn.Close()
	}()

	fmt.Printf("\tPlayer %s just logged in... let's play!\n", username)

	go func() {
		for {
			command := <-broadcast

			if command.command == "close" {
				break
			}

			if command.sourceUsername != username {
				sockWriter.WriteString(fmt.Sprintln(command.command))
				sockWriter.Flush()
			}
		}
	}()

GameConnectionLoop:
	for {
		fmt.Println("Got here")
		payloadBytes, err := sockReader.ReadBytes('\n')
		payload := string(payloadBytes)

		if err == io.EOF {
			fmt.Println("EOF, quitting")
			break GameConnectionLoop
		}

		if err != nil {
			fmt.Println(err)
			break GameConnectionLoop
		}
		fmt.Println(fmt.Sprintf("\t\tMessage: %s", payload))
		if payload == ":logout" {
			fmt.Println("Logging out...")
			break GameConnectionLoop
		}

		if payload == "walk to city hall\n" {
			_, err = sockWriter.WriteString(fmt.Sprintln("You enter City Hall"))
			broadcast <- BroadcastCommand{sourceUsername: username, command: fmt.Sprintf("%s has entered City Hall and is near you", username)}
			if err != nil {
				fmt.Println("Errrrrrrror!")
			}
			sockWriter.Flush()
			fmt.Println("Sent message")
		}
	}
}
