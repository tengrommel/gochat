package main

import (
	"flag"
	"strings"
	"net"
	"bufio"
	"os"
	"fmt"
	"log"
)

func main()  {
	op := flag.String("type","", "Server (s) or client (c) ?")
	address := flag.String("addr",":8000", "address? host:port")
	flag.Parse()
	switch strings.ToUpper(*op) {
	case "S":
		runServer(*address)
	case "C":
		runClient(*address)
	}
}

func runUDPClient(address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil{
		return err
	}
	defer conn.Close()
	
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("What message would you like to send?")
	for scanner.Scan() && err == nil {
		fmt.Println("Writing ", scanner.Text())
		conn.Write(scanner.Bytes())
		fmt.Println("What message would you like to send?")
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println(string(buffer))
	}
	return nil
}

func runUDPServer(address string) error {
	pc, err := net.ListenPacket("ucp", address)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("Listening......")
	defer pc.Close()
	buffer := make([]byte, 1024)
	for {
		_, addr,_ := pc.ReadFrom(buffer)
		fmt.Println("Received %s from address %s \n", string(buffer), addr)
		_, err := pc.WriteTo([]byte("message Received"), addr)
		if err != nil{
			log.Fatal("Could not write ", err)
		}
	}
}
