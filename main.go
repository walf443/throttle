package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var interval *time.Duration = flag.Duration("interval", 1000*time.Millisecond, "duration to wait output")
var debug *bool = flag.Bool("debug", false, "debug mode")

func main() {
	flag.Parse()
	args := flag.Args()
	execCommand := ""
	if len(args) == 1 {
		execCommand = args[0]
	}
	done := make(chan int)
	out := make(chan string)
	go inputStream(out, done)

	intervalCh := make(chan int)
	willShutdown := make(chan int)
	go background(out, execCommand, intervalCh, willShutdown)

	for {
		select {
		case <-intervalCh:
		case <-done:
			willShutdown <- 1
			<-intervalCh
			return
		}
	}
}

func inputStream(out chan<- string, done chan<- int) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	done <- 1
}

func background(input <-chan string, execCmd string, intervalCh chan<- int, willShutdown <-chan int) {
	buffer := make([]string, 0)
	timer := time.Tick(*interval)
	flush := func() {
		tmp := strings.Join(buffer, "\n")
		if tmp == "" {
			intervalCh <- 1
			return
		}

		if execCmd == "" {
			fmt.Println(tmp)
		} else {
			cmd := fmt.Sprintf(execCmd, tmp)
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
		intervalCh <- 1
	}
	for {
		select {
		case line := <-input:
			buffer = append(buffer, line)
		case <-willShutdown:
			flush()
		case <-timer:
			flush()
		}
	}
}
