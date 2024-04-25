package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
)

const (
	DEFAULT_PORT = 8322
	NUM_CHANNELS = 2
	SAMPLE_RATE  = 44100
)

var (
	BYTE_ORDER = binary.BigEndian
	BUFFER     = make([]int32, 1024)
)

func main() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatalf("Failed to init portaudio! " + err.Error())
	}
	defer portaudio.Terminate()
	switch os.Args[1] {
	case "server":
		Server()
	case "client":
		Client(os.Args[2])
	default:
		fmt.Printf("usage: %s [client|server] host:port\n", os.Args[0])
	}
}

func Server() {
	stream, err := portaudio.OpenDefaultStream(NUM_CHANNELS, 0, SAMPLE_RATE, len(BUFFER), BUFFER)
	CheckError(err, "Failed to open input streams!")
	defer stream.Close()
	addr := fmt.Sprintf("0.0.0.0:%d", DEFAULT_PORT)
	fmt.Printf("Waiting for a client on %s...\n", addr)
	listener, err := net.Listen("tcp", addr)
	CheckError(err, "Cannot listen to this address!")
	defer listener.Close()
	conn, err := listener.Accept()
	CheckError(err, "Can't listen for clients!")
	fmt.Println("Client connected! Forwarding...")
	CheckError(stream.Start(), "Cannot start streaming!")
	defer stream.Stop()
	time.Sleep(time.Second / 5)
	for {
		CheckError(stream.Read(), "Failed to read streaming data!")
		CheckError(binary.Write(conn, BYTE_ORDER, BUFFER), "Impossibly to write data to the client!")
	}
}

func Client(server string) {
	stream, err := portaudio.OpenDefaultStream(0, NUM_CHANNELS, SAMPLE_RATE, len(BUFFER), &BUFFER)
	CheckError(err, "Cannot open output streams!")
	defer stream.Close()
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", server)
	CheckError(err, "Failed to connect!")
	fmt.Println("Connected!")
	defer stream.Stop()
	for {
		err := binary.Read(conn, BYTE_ORDER, BUFFER)
		if err == io.EOF {
			fmt.Println("Connection closed!")
			break
		}
		CheckError(err, "Impossibly to read from server!")
		err = stream.Write()
		if err == portaudio.StreamIsStopped {
			CheckError(stream.Start(), "Failed to start streaming!")
			stream.Write()
		} else if err != portaudio.OutputUnderflowed && err != nil {
			CheckError(err, "Cannot write to streamer!")
		}
	}
}

func CheckError(err error, message string) {
	if err != nil {
		fmt.Printf("%T %v", err, err)
		log.Fatalf("%s %s\n", message, err.Error())
	}
}
