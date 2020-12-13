package internal

import (
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"errors"
	"io"
)

var Source string
var Target string
var openConnections int = 0
var MaxConnections int = 20

type Host struct {
	address net.IP
	port int
}

func (h *Host) getEndpoint() string {
	return h.address.String() + ":" + strconv.Itoa(h.port)
}

func getHost(endpoint string) (Host, error) {
	host := Host{}
	parts := strings.Split(endpoint, ":")
	if len(parts) != 2 {
		return host, errors.New("Malformed endpoint, it must be in the form of <bind address>:<port>, eg: 0.0.0.0:80")
	}
	host.address = net.ParseIP(parts[0])
	if host.address == nil {
		return host, errors.New("Malformed IP address" + parts[0])
	}
	i, err := strconv.Atoi(parts[1])
	if err != nil {
		return host, err
	}
	if i <0 || i > 65535 {
		return host, errors.New("Malformed port, it must be in [0, 65535]")
	}
	host.port = i
	return host, nil
}

func loop(sourceConn net.Conn, targetConn net.Conn) {
	for {
		_, err := io.Copy(sourceConn, targetConn)
		if err != nil {
			log.Fatal(err)
			log.Println("Closing channel due error in communication between source to target")
			break
		}
	}
	sourceConn.Close()
	log.Println("Channel closed")
}

func channel(sourceConn net.Conn, target Host) error {
	log.Println("Creating a channel...")
	targetConn, err := net.Dial("tcp", target.getEndpoint())
	if err != nil {
		return err
	}
	go loop(sourceConn, targetConn)
	go loop(targetConn, sourceConn)
	return nil
}

func Open() {
	log.Println("Starting server...")
	log.Println("Number of max connections: ", MaxConnections)
	source, err := getHost(Source)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	target, err2 := getHost(Target)
	if err2 != nil {
		log.Fatalln(err2)
		os.Exit(1)
	}
	listener, err3 := net.Listen("tcp", source.getEndpoint())
	if err3 != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		client, err4 := listener.Accept()
		if err4 != nil {
			log.Fatalln(err4)
		}
		if openConnections >= MaxConnections {
			log.Println("Maximum connections reached, a new incoming channel has been blocked")
			err4 = client.Close()
			if err4 != nil {
				log.Fatal(err4)
			}
		} else {
			log.Println("Creating a new channel...")
			channel(client, target)
			openConnections++
			log.Println("Open connections -> ", openConnections)
		}
	}
}