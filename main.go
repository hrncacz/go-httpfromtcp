package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	filepath, err := filepath.Abs("./message.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for {
		readBuffer := make([]byte, 8)
		_, err := file.Read(readBuffer)
		if err != nil {
			os.Exit(0)
		}
		fmt.Printf("read: %s\n", readBuffer)
	}
}
