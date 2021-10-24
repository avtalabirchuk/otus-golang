package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "use it to specify dial timeout")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Please use: telnet host port [--timeout=2s]\n")
	}
	flag.Parse()
	flag.Usage()
	telnetArgs := flag.Args()
	if len(telnetArgs) < 2 {
		log.Fatal("Correct usage: telnet host port [--timeout=2s]")
	}

	host := telnetArgs[0]
	port := telnetArgs[1]
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(
		address,
		*timeout,
		os.Stdin,
		os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Connected to %s...", address)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		client.Send()
		err := client.Send()
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		client.Receive()
		log.Println("...Connection was closed by peer")
		cancel()
	}()

	<-ctx.Done()
}
