package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
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

	defer func() {
		fmt.Println("Clossing Connection From " + con.RemoteAddr().String())
	}()

	fmt.Println("Connection Receaved From " + con.RemoteAddr().String())

	//timeoutDur := 5 * time.Second
	//bufReader := bufio.NewReader(con)
	//bufWriter := bufio.NewWriter(con)

	bufWriter := bufio.NewWriter(con)
	enc := gob.NewEncoder(bufWriter)
	enc.Encode(privKey.PublicKey)
	bufWriter.Flush()

	//con.Write([]byte("test stff"))

	for {
		//con.SetReadDeadline(time.Now().Add(timeoutDur))
		var cryptBytes []byte
		//cryptBytes, err := bufReader.ReadBytes('\n')

		//i, err := bufReader.Read(cryptBytes)
		//if err != nil {
		//	panic(err)
		//}
		i, err := con.Read(cryptBytes)
		if err != nil {
			panic(err)
		}

		fmt.Printf("len is %d", i)

		//bufWriter.WriteString("test")

		//here we would decrypt
		fmt.Println("Decrypting Message")
		fmt.Println(cryptBytes)
		decryptedBytes, err := decrypt(ranReader, privKey, cryptBytes)
		//if err != nil {
		//	panic(err)
		//}
		fmt.Printf("%s", decryptedBytes)
	}

	//enc := gob.NewEncoder@con.
	//fmt.Println("Connection Receaved From " + con.RemoteAddr().String())
	//con.Write([]byte(privKey.PublicKey))

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
	//get the public key
	bufReader := bufio.NewReader(conn)
	dec := gob.NewDecoder(bufReader)

	var pubKey = rsa.PublicKey{}
	dec.Decode(&pubKey)

	fmt.Println("Public Key Receaved")
	cryptText, err := encrypt(ranReader, &pubKey, []byte(message))
	if err != nil {
		panic(err)
	}

	conn.Write(cryptText)
	//fmt.Fprintf(conn, message)
	//status, err := bufio.NewReader(conn).ReadString('\n')
}

func main() {
	targetIP := getIP()
	fmt.Println("Target IP has been set to", targetIP)

	fmt.Println("Creating Entropy")
	ranReader = rand.Reader

	genPrivKey(ranReader)
	fmt.Println("Generated Key Pair")

	var startServerBool string
	fmt.Printf("Start Server? (Yes/No): ")
	fmt.Scanln(&startServerBool)
	if startServerBool == "Yes" {
		go startServer()
	}

	//send("localhost", "this is a test :)")
	//fmt.Scanln()

	// Sending Portion

	conn, err := net.Dial("tcp", targetIP+":"+strconv.Itoa(port)) // dials the server over tcp with the IP:Port 	Itoa converts ints to strings
	if err != nil {
		fmt.Println("Connection failed. Is port is use?")
		panic(err)
	}

	//get the public key
	bufReader := bufio.NewReader(conn)
	dec := gob.NewDecoder(bufReader)

	var pubKey = rsa.PublicKey{}
	dec.Decode(&pubKey)

	fmt.Println("Public Key Receaved")

	var message string

	for {
		fmt.Scanln(&message)
		cryptText, err := encrypt(ranReader, &pubKey, []byte(message))
		if err != nil {
			panic(err)
		}
		fmt.Println("sending CryptText")
		fmt.Println(cryptText)
		conn.Write(cryptText)

	}
}
