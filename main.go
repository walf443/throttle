package main

import (
	"flag"
	"bufio"
	"os"
	"log"
	"time"
	"fmt"
)


func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	ch := make(chan string)
	go background(ch)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func background(input chan string) {
	buffer := make([]string, 0)
	timer := time.Tick(200 * time.Millisecond)
	for {
		select {
		case line := <-input:
			buffer = append(buffer, line)
		case <-timer:
			for _, buf := range(buffer) {
				fmt.Println(buf)
			}
			buffer = buffer[:0]
		}
	}
}
