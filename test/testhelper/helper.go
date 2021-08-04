package testhelper

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"

	"github.com/drgomesp/rhizom/test/fake"
)

func FailOnError(t *testing.T, err error) {
	t.Helper()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func LocalListener(t *testing.T) net.Listener {
	t.Helper()
	netT, err := nettest.NewLocalListener("tcp")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return netT
}

func MakePipe() (c1, c2 net.Conn, stop func(), err error) {
	b1, b2 := make(chan []byte), make(chan []byte)
	c1 = fake.NewFakeConn(b1, b2)
	c2 = fake.NewFakeConn(b2, b1)
	stop = func() { c1.Close(); c2.Close() }
	return c1, c2, stop, nil
}
