package main

import (
	"flag"
	"bufio"
	"os"
	"os/exec"
	"log"
	"time"
	"fmt"
	"strings"
)

var interval *int = flag.Int("interval", 1000, "millisecond to wait output")
var debug *bool = flag.Bool("debug", false, "debug mode")

func main() {
	flag.Parse()
	args := flag.Args()
	execCommand := ""
	if len(args) == 1 {
		execCommand = args[0]
	}
	graceful := make(chan int)
	done := make(chan int)
	out := make(chan string)
	go inputStream(out, done)
	go background(out, graceful, execCommand)
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

func background(input chan string, graceful chan int, execCommand string) {
	buffer := make([]string, 0)
	timer := time.Tick(time.Duration(*interval) * time.Millisecond)
	for {
		select {
		case line := <-input:
			buffer = append(buffer, line)
		case <-timer:
			tmp := strings.Join(buffer, "\n")
			if tmp != "" {
				if execCommand == "" {
					fmt.Println(tmp)
				} else {
					cmd := fmt.Sprintf(execCommand, tmp)
					if *debug {
						fmt.Printf("[debug] execute \"%s\"\n", cmd)
					}
					out, err := exec.Command("/bin/bash", "-c", cmd).Output()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(string(out))
				}
				buffer = buffer[:0]
			}
			graceful <- 1
		}
	}
}
