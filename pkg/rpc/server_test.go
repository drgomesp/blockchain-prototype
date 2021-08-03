package rpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
)

func testHelper_NewServer(t *testing.T, name string) *Server {
	t.Helper()
	s := NewServer(name)
	assert.NotNil(t, s)
	return s
}

func testHelper_NewListener(t *testing.T) Listener {
	t.Helper()
	port := 7000
	l := NewListener(port)
	defer l.Close()
	assert.NotNil(t, l)
	return l
}

func TestServer_Name(t *testing.T) {
	const name = "test"
	s := testHelper_NewServer(t, name)
	got := s.Name()
	assert.Equal(t, name, got)
}

func TestServer_Info(t *testing.T) {
	const name = "test"
	s := testHelper_NewServer(t, name)
	got := s.Info()
	assert.IsType(t, map[string]grpc.ServiceInfo{}, got)
}

func TestServer(t *testing.T) {
	// nettest local network listener
	netT := testHelper_NewListener(t)
	defer assert.NoError(t, netT.Close())

	// run Server
	s := testHelper_NewServer(t, "test")
	go s.Start(netT)
	defer s.Stop()
	netT.Accept()
	// setup gRPC connection
	conn, err := grpc.Dial(netT.Addr().String(), grpc.WithInsecure())
	if !assert.NoError(t, err) || !assert.NotNil(t, err) {
		t.FailNow()
	}
	defer assert.NoError(t, conn.Close())

	// setup client request
	stream, err := service.NewNodeClient(conn).GetBlock(context.Background())
	if !assert.NoError(t, err) || !assert.NotNil(t, err) {
		t.FailNow()
	}
	defer assert.NoError(t, stream.CloseSend())

	// send client requests
	for i := uint32(1); i <= 10; i++ {
		if !assert.NoError(t, stream.Send(&message.GetBlockRequest{Want: i})) {
			t.FailNow()
		}
		resp, err := stream.Recv()
		if !assert.NoError(t, err) || !assert.NotNil(t, resp) {
			t.FailNow()
		}
	}
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

func testHelper_sendToChan(ch chan<- []byte, b []byte, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	ch <- b
	return len(b), nil
}

func testHelper_makePipe() (c1, c2 net.Conn, stop func(), err error) {
	makeBuf := func() chan []byte { return make(chan []byte) }
	b1, b2 := makeBuf(), makeBuf()

	c1 = testHelper_newFakeConn(b1, b2)
	c2 = testHelper_newFakeConn(b2, b1)
	stop = func() { c1.Close(); c2.Close() }

	return c1, c2, stop, nil
}

func testHelper_nettest(t *testing.T) net.Listener {
	netT, err := nettest.NewLocalListener("tcp")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return netT
}
