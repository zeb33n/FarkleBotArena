package game

// I dont know if game is the right name for this package but it is the package that
// interfcts with the game server reading the current state and sending player decisions

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	// DataCh  chan []byte
	// ErrCh   chan error
	CloseCh chan struct{}
}

// returns an empty client. Connect/close will be called via tea.Cmd through the main model
func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Client) Read() (<-chan []byte, <-chan error) {
	DataCh := make(chan []byte)
	ErrCh := make(chan error)
	// Channel is updated with new game state from the server
	go func() {
		defer close(DataCh)
		defer close(ErrCh)

		for {
			buffer := make([]byte, 1024)
			n, err := c.conn.Read(buffer)
			if err != nil {
				ErrCh <- err
				// do we want to return if we have an error reading?
			}
			fmt.Printf("normal:\n")
			fmt.Printf("%s\n", buffer)

			cleanedBuff := []byte{}
			for _, b := range buffer[:n] {
				if b != 0 {
					cleanedBuff = append(cleanedBuff, b)
				} else {
					break
				}
			}
			fmt.Printf("cleaned:\n")
			fmt.Printf("%s\n", cleanedBuff)
			DataCh <- cleanedBuff

		}
	}()

	return DataCh, ErrCh
}

// pretty sure the client at the moment only accepts/reads 1s and expects it for
// players roll decision

func (c *Client) Respond(r []byte) error {
	_, err := c.conn.Write(r)
	if err != nil {
		return err
	}
	return nil
}
