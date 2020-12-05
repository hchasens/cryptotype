package main

import ("fmt"
  "net"
	"crypto/rsa"
	"crypto/rand"
)

/**
I've loaded this up with some some example code. This is just for reference .
**/

func getIP() string {
	var ip string
	fmt.Printf("Set Target IP: ")
	fmt.Scanln(&ip)
	return ip
}

func genPrivKey() *rsa.PrivateKey{ 
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	return privateKey;
}


func startServer() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		panic(err)
	}
  defer l.Close() // wait till the function returns then close the listner
}

func main() {
	targetIP := getIP();
	fmt.Println("Target IP has been set to", targetIP)


	privateKey := genPrivKey()
	publicKey := privateKey.PublicKey
	fmt.Println(publicKey)

}

