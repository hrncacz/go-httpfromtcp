package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ip, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ip)
	conn, err := net.DialUDP("udp", nil, ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("CHYBAAA")
		}
		conn.Write([]byte(userInput))
	}
}
