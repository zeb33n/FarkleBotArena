package game

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

// MockConn is a mock implementation of net.Conn for testing
type MockConn struct {
	ReadData []byte
	ReadErr  error
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	if m.ReadErr != nil {
		return 0, m.ReadErr
	}
	n = copy(b, m.ReadData)
	return n, io.EOF
}

// Other required methods of net.Conn interface...
func (m *MockConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return nil }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestClientRead(t *testing.T) {
	testCases := []struct {
		name     string
		readData []byte
		readErr  error
	}{
		{
			name:     "Successful read",
			readData: []byte("test data"),
			readErr:  nil,
		},
		{
			name:     "Read error",
			readData: nil,
			readErr:  errors.New("read error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConn := &MockConn{
				ReadData: tc.readData,
				ReadErr:  tc.readErr,
			}

			client := &Client{
				conn: mockConn,
				log:  log.New(io.Discard, "", 0), // Discard logs in test
			}

			dataCh, errCh := client.Read()

			// Use a timeout to prevent the test from hanging
			timeout := time.After(1 * time.Second)

			select {
			case data := <-dataCh:
				if tc.readErr == nil {
					if !bytes.Equal(data, tc.readData) {
						t.Errorf("Expected data %v, got %v", tc.readData, data)
					}
				} else {
					t.Errorf("Expected error, got data: %v", data)
				}
			case err := <-errCh:
				if tc.readErr == nil {
					t.Errorf("Expected data, got error: %v", err)
				} else if err.Error() != tc.readErr.Error() {
					t.Errorf("Expected error %v, got %v", tc.readErr, err)
				}
			case <-timeout:
				t.Error("Test timed out")
			}
		})
	}
}
