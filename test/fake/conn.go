package fake

import (
	"net"
	"time"
)

type FakeConn struct {
	Rd, Wr chan []byte
	Addr   net.Addr
	Err    error
}

func NewFakeConn(rd, wr chan []byte) *FakeConn {
	return &FakeConn{Rd: rd, Wr: wr}
}

func (f *FakeConn) Close() error {
	if f.Err != nil {
		return f.Err
	}
	close(f.Rd)
	close(f.Wr)
	return nil
}

func (f *FakeConn) Write(b []byte) (n int, Err error) {
	if f.Err != nil {
		return 0, f.Err
	}
	f.Wr <- b
	return len(b), nil
}

func (f *FakeConn) Read(b []byte) (n int, Err error)   { b = <-f.Rd; return len(b), f.Err }
func (f *FakeConn) LocalAddr() net.Addr                { return f.Addr }
func (f *FakeConn) RemoteAddr() net.Addr               { return f.Addr }
func (f *FakeConn) SetDeadline(_ time.Time) error      { return f.Err }
func (f *FakeConn) SetReadDeadline(_ time.Time) error  { return f.Err }
func (f *FakeConn) SetWriteDeadline(_ time.Time) error { return f.Err }
