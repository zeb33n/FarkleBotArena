package game

// I dont know if game is the right name for this package but it is the package that
// interfcts with the game server reading the current state and sending player decisions

import (
	"log"
	"net"
)

type Client struct {
	log  *log.Logger
	conn net.Conn
	// DataCh  chan []byte
	// ErrCh   chan error
	CloseCh chan struct{}
}

// returns an empty client. Connect/close will be called via tea.Cmd through the main model
func NewClient(log *log.Logger) *Client {
	return &Client{log: log}
}

func (c *Client) Connect(addr string) error {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// I think this is getting called too early, I think we only want to call read once 1 has been pressed?
func (c *Client) Read() (<-chan []byte, <-chan error) {

	c.log.Println("making channels")
	DataCh := make(chan []byte)
	ErrCh := make(chan error)
	// Channel is updated with new game state from the server
	go func() {

		c.log.Println("the go routine has started")

		defer close(DataCh)
		defer close(ErrCh)

		buffer := make([]byte, 1024)
		c.log.Println("buffer created inside routie")
		for {

			c.log.Println("why are we never inside the for loop?!")

			n, err := c.conn.Read(buffer)
			c.log.Println("read attempted")
			if err != nil {
				c.log.Println("error during read ", err)
				ErrCh <- err
				// do we want to return if we have an error reading?
			}

			c.log.Println("no error during reading of connection")

			cleanedBuff := []byte{}
			for _, b := range buffer[:n] {
				if !(b == 0) {
					cleanedBuff = append(cleanedBuff, b)
				} else {
					break
				}
			}

			c.log.Println(cleanedBuff)
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
