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

func main() {
	port = 2424
	ip = "localhost"

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
