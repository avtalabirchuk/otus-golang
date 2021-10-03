package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
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
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Connect to %v...", address)
	sigintChannel := make(chan os.Signal, 1)

	signal.Notify(sigintChannel, syscall.SIGINT)

	go func() {
		<-sigintChannel
		fmt.Println("Got SIGINT")
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			err := client.Receive()
			if err != nil {
				log.Println(err)
				break
			}
		}
		wg.Done()
	}()
	go func() {
		for {
			err := client.Send()
			if err != nil {
				log.Println(err)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
