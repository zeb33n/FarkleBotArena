package game

// I dont know if game is the right name for this package but it is the package that
// interfcts with the game server reading the current state and sending player decisions

import "net"

type client struct {
	conn    net.Conn
	DataCh  chan []byte
	ErrCh   chan error
	CloseCh chan struct{}
}

// returns an empty client. Connect/close will be called via tea.Cmd through the main model
func NewClient() *client {
	return &client{}
}

func (c *client) Connect(addr string) error {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *client) Read() {
	// Channel is updated with new game state from the server
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := c.conn.Read(buffer)
			if err != nil {
				c.ErrCh <- err
				// do we want to return if we have an error reading?
			}

			cleanedBuff := []byte{}
			for _, b := range buffer[:n] {
				if !(b == 0) {
					cleanedBuff = append(cleanedBuff, b)
				} else {
					break
				}
			}
			c.DataCh <- cleanedBuff

		}
	}()
}

// pretty sure the client at the moment only accepts/reads 1s and expects it for
// players roll decision

func (c *client) Respond(r []byte) error {
	_, err := c.conn.Write(r)
	if err != nil {
		return err
	}
	return nil
}
