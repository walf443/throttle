package main

import (
	"flag"
	"bufio"
	"os"
	"fmt"
	"log"
)


func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
