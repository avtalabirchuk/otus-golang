package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalln("tuc")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		log.Fatalln("failed to connect:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		client.Send()
		cancel()
	}()

	go func() {
		client.Receive()
		cancel()
	}()

	<-ctx.Done()
}
