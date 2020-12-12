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

func handleChannel(client net.Conn, target Host) {
	targetConn, err := net.Dial("tcp", target.getEndpoint())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer targetConn.Close()
	var totalReceived int64 = 0
	for {
		received, err := io.Copy(client, targetConn)
		if err != nil {
			log.Fatal(err)
			log.Println("Closing channel")
			break
		}
		totalReceived += received
	}
	client.Close()
	openConnections--
}

func Open() {
	log.Println("Starting server...")
	log.Println("Number of max connections: ", MaxConnections)
	ss, err := getHost(Source)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	tt, err2 := getHost(Target)
	if err2 != nil {
		log.Fatalln(err2)
		os.Exit(1)
	}
	listener, err3 := net.Listen("tcp", ss.getEndpoint())
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
			go handleChannel(client, tt)
			openConnections++
			log.Println("Open connections -> ", openConnections)
		}
	}
}