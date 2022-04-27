package infrastructure

import (
	"bufio"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockTransServerNoHandler(t *testing.T) {
	srv := NewMockTransServer()
	defer srv.Close()

	conn, err := net.Dial("tcp", srv.Address)
	assert.NoError(t, err)
	defer conn.Close()

	welcome, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, WelcomeMessage, welcome)

	_, err = conn.Write([]byte("cmd:foo\ncommit:1\nend\n"))
	assert.NoError(t, err)

	end, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, EndMessage, end)
}

func TestMockTransServerHandler(t *testing.T) {
	expectedArgs := []byte("cmd:foo\ncommit:1\nend\n")

	expectedRes := []byte("foo:bar\n")

	handler := func(args []byte) []byte {
		assert.ElementsMatch(t, expectedArgs, args)
		return expectedRes
	}

	srv := NewMockTransServer()
	srv.SetHandler(handler)
	defer srv.Close()

	conn, err := net.Dial("tcp", srv.Address)
	assert.NoError(t, err)
	defer conn.Close()

	welcome, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, WelcomeMessage, welcome)

	_, err = conn.Write([]byte("cmd:foo\ncommit:1\nend\n"))
	assert.NoError(t, err)

	res, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, string(expectedRes), res)
}

func TestMockTransServerBusy(t *testing.T) {
	srv := NewMockTransServer()
	srv.SetBusy(true)
	defer srv.Close()

	conn, err := net.Dial("tcp", srv.Address)
	assert.NoError(t, err)
	defer conn.Close()

	welcome, err := bufio.NewReader(conn).ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, BusyMessage, welcome)

	_, err = conn.Write([]byte("cmd:foo\ncommit:1\nend\n"))
	assert.NoError(t, err)

	res, err := bufio.NewReader(conn).ReadString('\n')
	assert.Error(t, err)
	assert.Empty(t, res)
}
