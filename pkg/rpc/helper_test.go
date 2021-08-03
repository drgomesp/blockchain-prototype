package rpc

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
)

func testHelper_NewListener(t *testing.T) Listener {
	t.Helper()
	port := 7000
	l := NewListener(port)
	defer l.Close()
	assert.NotNil(t, l)
	return l
}

func testHelper_NewServer(t *testing.T, name string, svc *BlockService) *Server {
	t.Helper()
	s := NewServer(name, svc)
	assert.NotNil(t, s)
	return s
}

func testHelper_NewStreamService(t *testing.T) *BlockService {
	t.Helper()
	s := NewStreamService()
	assert.NotNil(t, s)
	return s
}

func testHelper_localListener(t *testing.T) net.Listener {
	netT, err := nettest.NewLocalListener("tcp")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return netT
}

func testHelper_makePipe() (c1, c2 net.Conn, stop func(), err error) {
	makeBuf := func() chan []byte { return make(chan []byte) }
	b1, b2 := makeBuf(), makeBuf()

	c1 = testHelper_newFakeConn(b1, b2)
	c2 = testHelper_newFakeConn(b2, b1)
	stop = func() { c1.Close(); c2.Close() }

	return c1, c2, stop, nil
}

func testHelper_sendToChan(ch chan<- []byte, b []byte, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	ch <- b
	return len(b), nil
}

type (
	fakeListener struct {
		FakeConn *fakeConn
		Err      error
	}
	fakeConn struct {
		Rd, Wr chan []byte
		Addr   net.Addr
		Err    error
	}
)

func testHelper_newFakeConn(rd, wr chan []byte) *fakeConn {
	return &fakeConn{Rd: rd, Wr: wr}
}
func testHelper_newFakeListener(fake net.Conn) *fakeListener {
	return &fakeListener{FakeConn: fake.(*fakeConn)}
}

func (f *fakeConn) Close() error {
	if f.Err != nil {
		return f.Err
	}
	close(f.Rd)
	close(f.Wr)
	return nil
}

func (f *fakeConn) Read(b []byte) (n int, Err error)   { b = <-f.Rd; return len(b), f.Err }
func (f *fakeConn) Write(b []byte) (n int, Err error)  { return testHelper_sendToChan(f.Wr, b, f.Err) }
func (f *fakeConn) LocalAddr() net.Addr                { return f.Addr }
func (f *fakeConn) RemoteAddr() net.Addr               { return f.Addr }
func (f *fakeConn) SetDeadline(_ time.Time) error      { return f.Err }
func (f *fakeConn) SetReadDeadline(_ time.Time) error  { return f.Err }
func (f *fakeConn) SetWriteDeadline(_ time.Time) error { return f.Err }
func (f *fakeListener) Accept() (net.Conn, error)      { return f.FakeConn, f.Err }
func (f *fakeListener) Close() error                   { return f.FakeConn.Close() }
func (f *fakeListener) Addr() net.Addr                 { return f.FakeConn.Addr }
