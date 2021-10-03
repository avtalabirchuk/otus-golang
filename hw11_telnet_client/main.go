package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var durTime = 10

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	if len(os.Args) < 3 {
		log.Fatal("Please use: telnet host port [--timeout=2s]")
	}
	timeout := flag.Duration("timeout", time.Duration(durTime), "use it to specify dial timeout")
	flag.Parse()
	log.Println(os.Args[1:])
	host := os.Args[2]
	port := os.Args[3]
	address := net.JoinHostPort(host, port)
	log.Println(address)
	client := NewTelnetClient(
		address,
		*timeout,
		os.Stdin,
		os.Stdout)

	defer client.Close()
	if err := client.Connect(); err != nil {
		panic(err)
	}
	log.Printf("Connect to %v...", address)
	sigintChannel := make(chan os.Signal, 1)
	doneCh := make(chan struct{})

	signal.Notify(sigintChannel, syscall.SIGINT)

	go func() {
		<-sigintChannel
		fmt.Println("Got SIGINT")
		doneCh <- struct{}{}
	}()
	go func() {
		log.Println("Start receiving")
		for {
			err := client.Receive()
			if err != nil {
				log.Println("Error during receive:", err)
				break
			}
			log.Println("Data received")
		}
		doneCh <- struct{}{}
	}()
	go func() {
		defer client.Close()
		log.Println("Start sending")
		for {
			err := client.Send()
			if err != nil {
				log.Println("Error during send:", err)
				break
			}
			log.Println("Data send")
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
}
