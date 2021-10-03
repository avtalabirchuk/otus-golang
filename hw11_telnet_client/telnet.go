package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TCPClient struct {
	address string
	timeout time.Duration

	inReader *bufio.Reader
	out      io.Writer

	conn       net.Conn
	connReader *bufio.Reader
}

func (tc *TCPClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	tc.connReader = bufio.NewReader(conn)
	return nil
}

func (tc *TCPClient) Close() error {
	if err := tc.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (tc *TCPClient) Send() error {
	data, err := tc.inReader.ReadBytes('\n')
	if err != nil {
		return err
	}
	_, err = tc.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TCPClient) Receive() error {
	data, connErr := tc.connReader.ReadBytes('\n')
	log.Println("data receive:", string(data))
	_, outWriteErr := tc.out.Write(data)
	if connErr != nil {
		return fmt.Errorf("error receive data: %w", connErr)
	}
	if outWriteErr != nil {
		return fmt.Errorf("error outwrited data: %w", outWriteErr)
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &TCPClient{
		address:  address,
		timeout:  timeout,
		inReader: bufio.NewReader(in),
		out:      out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
