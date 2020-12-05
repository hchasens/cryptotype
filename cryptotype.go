package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"strconv"
)

/**
I've loaded this up with some some example code. This is just for reference .
**/

var port int = 4242

func getIP() string {
	var ip string
	fmt.Printf("Set Target IP: ")
	fmt.Scanln(&ip)
	return ip
}

func genPrivKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func connectionHandler(con net.Conn) {
	fmt.Println("connection receaved")
}

func startServer() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	fmt.Println("Started Server")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go connectionHandler(conn)
	}
	//defer l.Close() // wait till the function returns then close the listner
}

func send(ip string, message string) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port)) // dials the server over tcp with the IP:Port 	Itoa converts ints to strings
	if err != nil {
		fmt.Println("Connection failed. Is port is use?")
		panic(err)
	}
	fmt.Fprintf(conn, message)
	//status, err := bufio.NewReader(conn).ReadString('\n')
}

func main() {
	//targetIP := getIP()
	//fmt.Println("Target IP has been set to", targetIP)

	go startServer()
	fmt.Scanln()
	//privateKey := genPrivKey()
	//publicKey := privateKey.PublicKey
	//fmt.Println(publicKey)

}
