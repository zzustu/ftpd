package ftpd

import (
	"io"
	"net"
	"sync"
)

type DataConn interface {
	Read([]byte) (int, error)
	ReadFrom(io.Reader) (int64, error)
	Write([]byte) (int, error)
	Close() error
}

type portModeConn struct {
	conn       *net.TCPConn
	remoteAddr net.TCPAddr
}

func (c *portModeConn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *portModeConn) ReadFrom(r io.Reader) (int64, error) {
	return c.conn.ReadFrom(r)
}

func (c *portModeConn) Write(b []byte) (int, error) {
	sz, err := c.conn.Write(b)
	if err != nil {
		return sz, err
	}
	err = c.conn.CloseWrite()
	return sz, err
}

func (c *portModeConn) Close() error {
	return c.conn.Close()
}

type pasvModeConn struct {
	conn       *net.TCPConn
	remoteAddr net.TCPAddr
	mutex      sync.Mutex
}

func (c *pasvModeConn) Read(b []byte) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.conn.Read(b)
}

func (c *pasvModeConn) ReadFrom(r io.Reader) (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.conn.ReadFrom(r)
}

func (c *pasvModeConn) Write(b []byte) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.conn.Write(b)
}

func (c *pasvModeConn) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.conn.Close()
}

func newPortModeConn(addr *net.TCPAddr) (DataConn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	c := new(portModeConn)
	c.conn = conn
	c.remoteAddr = *addr

	return c, nil
}

func newPasvModeConn(addr *net.TCPAddr) (DataConn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	c := new(portModeConn)
	c.conn = conn
	c.remoteAddr = *addr

	return c, nil
}
