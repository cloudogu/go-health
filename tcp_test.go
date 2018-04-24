package health_test

import (
	"net"
	"testing"

	"time"

	"github.com/cloudogu/go-health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTCPHealthChecker(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	checker := health.NewTCPHealthChecker("localhost", port)
	err = checker()
	assert.NoError(t, err)
}

func TestNewTCPHealthCheckerWithNonOpenPort(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	checker := health.NewTCPHealthChecker("localhost", port)
	err = checker()
	assert.Error(t, err)
}

func TestNewTCPHealthCheckBuilderWithShortTimeout(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	checker := health.NewTCPHealthCheckBuilder(port).
		WithTimeout(2 * time.Second).
		Build()

	err = checker()
	assert.NoError(t, err)
}

func TestNewTCPHealthCheckBuilderWithUnknownHostname(t *testing.T) {
	checker := health.NewTCPHealthCheckBuilder(80).
		WithHostname("sorbot.haxor.local").
		WithTimeout(1 * time.Second).
		Build()

	err := checker()
	assert.Error(t, err)
}
