package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"strconv"
)

/**
I've loaded this up with some some example code. This is just for reference .
**/

var port int = 4242
var ranReader io.Reader
var privKey *rsa.PrivateKey

func getIP() string {
	var ip string
	fmt.Printf("Set Target IP: ")
	fmt.Scanln(&ip)
	return ip
}

func genPrivKey(ranReader io.Reader) *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(ranReader, 4096)
	if err != nil {
		panic(err)
	}
	privKey = privateKey
	return privateKey
}

func encrypt(rand io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand, pub, msg)
}

func decrypt(rand io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand, priv, ciphertext)
}

func connectionHandler(con net.Conn) { // one of these connection handlers is made for every connection made. It'll be in charge of the handshake
	fmt.Println("connection receaved")
	con.Write([]byte("Message received."))
	fmt.Println(con.RemoteAddr())

}

func startServer() {
	fmt.Println("Starting Server")
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
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
	targetIP := getIP()
	fmt.Println("Target IP has been set to", targetIP)

	fmt.Println("Creating Entropy")
	ranReader = rand.Reader

	genPrivKey(ranReader)

	fmt.Println("Generated Key Pair")

	go startServer()

	fmt.Scanln()

}
