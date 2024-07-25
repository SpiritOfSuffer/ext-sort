package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	fileSize = 2 * 1024 * 1024 * 1024 // 2 GB
)

func main() {
	file, err := os.Create("data.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	var totalBytes int64

	for totalBytes < fileSize {
		timestamp := rand.Int63() // generate random int64 timestamp
		line := fmt.Sprintf("%d\n", timestamp)
		n, err := file.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		totalBytes += int64(n)
	}

	fmt.Println("File generation completed")
}
