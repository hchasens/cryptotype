package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"
	"strconv"
)

var ip string
var port int

func connectionHandler(conn net.Conn, privKey *rsa.PrivateKey) {
	defer conn.Close() // should the function finish exicuting the connection will close

	// send the public key
	marshaledPubKey := x509.MarshalPKCS1PublicKey(&privKey.PublicKey) // len is 270
	conn.Write(marshaledPubKey)

	//bufReader := bufio.NewReader(conn)

	//message handler
	for {
		/**
		is, err := bufReader.ReadBytes('\n') //i in string form
		if err != nil {
			panic(err)
		}

		bufLen, err := strconv.Atoi(string(is)) // is the length of the incomming message

		fmt.Println("len is " + string(is) + "test ")

		buffer := make([]byte, bufLen)
		**/
		buffer := make([]byte, 256)
		//buffer := make([]byte, 1000)

		conn.Read(buffer)
		plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, buffer)
		if err != nil {
			print(err)
		}

		fmt.Println(conn.RemoteAddr().String() + ": " + string(plaintext))
	}
}

func main() {
	port = 2424
	ip = "localhost"

	//generate Priv RSA Key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	//start listener
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	for { // for (ever) accept all incoming connections,
		// pass each connection (conn) to the connection handler in its own thread
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go connectionHandler(conn, privKey)
	}
}
