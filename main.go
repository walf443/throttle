package main

import (
	"flag"
	"bufio"
	"os"
	"log"
	"time"
	"fmt"
)

var interval *int = flag.Int("interval", 1000, "millisecond to wait output")

func main() {
	flag.Parse()
	graceful := make(chan int)
	done := make(chan int)
	out := make(chan string)
	go inputStream(out, done)
	go background(out, graceful)
	for {
		select {
		case <-graceful:
		case <-done:
			<-graceful
			return
		}
	}
}

func inputStream(out chan string, done chan int) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	done <- 1
}

func background(input chan string, graceful chan int) {
	buffer := make([]string, 0)
	timer := time.Tick(time.Duration(*interval) * time.Millisecond)
	for {
		select {
		case line := <-input:
			buffer = append(buffer, line)
		case <-timer:
			for _, buf := range(buffer) {
				fmt.Println(buf)
			}
			buffer = buffer[:0]
			graceful <- 1
		}
	}
}
