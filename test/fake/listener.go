package fake

import (
	"net"
)

type FakeListener struct {
	FakeConn *FakeConn
	Err      error
}

func NewFakeListener(fake net.Conn) *FakeListener {
	return &FakeListener{FakeConn: fake.(*FakeConn)}
}

func (f *FakeListener) Accept() (net.Conn, error) { return f.FakeConn, f.Err }
func (f *FakeListener) Close() error              { return f.FakeConn.Close() }
func (f *FakeListener) Addr() net.Addr            { return f.FakeConn.Addr }
