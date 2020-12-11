package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"
	"os"
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

func serverStart() {

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

func main() {
	port = 2424
	ip = "localhost"

	fmt.Print("Start Server (Yes/No): ")
	var serverStatus string
	fmt.Scanln(&serverStatus)
	if serverStatus == "Yes" {
		go serverStart()
	}

	fmt.Print("Target IP: ")
	fmt.Scanln(&ip)

	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port)) // dials the server over tcp with the IP:Port 	Itoa converts ints to strings
	if err != nil {
		fmt.Println("Connection failed. Is port is use?")
		panic(err)
	}

	//get public key
	marshalisedPubKey := make([]byte, 270) // create a buffer of the right size
	conn.Read(marshalisedPubKey)

	pubKey, err := x509.ParsePKCS1PublicKey(marshalisedPubKey)
	if err != nil {
		panic(err)
	}

	for {
		msg := bufio.NewReader(os.Stdin)     //create input buffer
		plainText, _ := msg.ReadString('\n') //read  until new line

		//encrypt with public key
		cryptText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(plainText))
		if err != nil {
			panic(err)
		}

		// outgoing connections are formated as such. stringOfCryptTextLeng, \n(00), cryptText
		// this allows us to parse everything without having too much room on the end of the buffer or buffer overflow

		//find the length of cryptText so we can make a buffer of the right size on the server
		/**
		var i string = strconv.Itoa(len(cryptText)) + "\n" // we change this to a string because you can't cast a int as a []byte array

		conn.Write([]byte(i)) // send len of message as string with EOL
		if err != nil {
			panic(err)
		}
		**/

		conn.Write(cryptText) // send crypttext
		if err != nil {
			panic(err)
		}

	}
}
