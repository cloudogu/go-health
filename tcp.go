package health

import (
	"net"
	"strconv"

	"time"

	"github.com/pkg/errors"
)

// TCPHealthCheckOptions defines options for the tcp health check
type TCPHealthCheckOptions struct {
	Hostname string
	Port     int
	Timeout  time.Duration
}

// TCPHealthCheckBuilder build for tcp health checks
type TCPHealthCheckBuilder struct {
	options *TCPHealthCheckOptions
}

// WithHostname sets the hostname for the tcp connection
func (builder *TCPHealthCheckBuilder) WithHostname(hostname string) *TCPHealthCheckBuilder {
	builder.options.Hostname = hostname
	return builder
}

// WithTimeout sets the timeout for the tcp connection
func (builder *TCPHealthCheckBuilder) WithTimeout(timeout time.Duration) *TCPHealthCheckBuilder {
	builder.options.Timeout = timeout
	return builder
}

// Builder creates a tcp http health checker with the configured options
func (builder *TCPHealthCheckBuilder) Build() HealthChecker {
	return NewTCPHealthCheckerWithOptions(builder.options)
}

// NewTCPHealthCheckBuilder creates a builder for tcp health checks
func NewTCPHealthCheckBuilder(port int) *TCPHealthCheckBuilder {
	return &TCPHealthCheckBuilder{
		options: createTCPHealthCheckOptions(port),
	}
}

// NewTCPHealthChecker creates a new tcp health check for the given host and port
func NewTCPHealthChecker(hostname string, port int) HealthChecker {
	options := createTCPHealthCheckOptions(port)
	options.Hostname = hostname
	return NewTCPHealthCheckerWithOptions(options)
}

func createTCPHealthCheckOptions(port int) *TCPHealthCheckOptions {
	return &TCPHealthCheckOptions{
		Hostname: "localhost",
		Port:     port,
		Timeout:  30 * time.Second,
	}
}

// NewTCPHealthCheckerWithOptions creates a new tcp health check for the given options
func NewTCPHealthCheckerWithOptions(options *TCPHealthCheckOptions) HealthChecker {
	return func() error {
		conn, err := net.DialTimeout("tcp", options.Hostname+":"+strconv.Itoa(options.Port), options.Timeout)
		if err != nil {
			return errors.Wrapf(err, "tcp connection to %s:%s failed", options.Hostname, options.Port)
		}
		defer conn.Close()
		return nil
	}
}
