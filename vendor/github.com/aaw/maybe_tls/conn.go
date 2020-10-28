package maybe_tls

import (
	"net"
	"time"
)

type StreamReplay struct {
	data []byte
	err  error
}

// Wraps a TCP connection and implements the net.Conn interface.
type Conn struct {
	conn   net.Conn
	replay *StreamReplay
}

// Read reads data from the connection.
// Read can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c Conn) Read(b []byte) (n int, err error) {
	if c.replay.err != nil {
		err = c.replay.err
		c.replay.err = nil
		return
	}
	rn := len(c.replay.data)
	if rn > 0 {
		n = len(b)
		if rn < n {
			n = rn
		}
		copy(b, c.replay.data)
		c.replay.data = c.replay.data[n:]
	} else {
		n, err = c.conn.Read(b)
	}
	return
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c Conn) Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c Conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future I/O, not just
// the immediately following call to Read or Write.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
// A zero value for t means Read will not time out.
func (c Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
