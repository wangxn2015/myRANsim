package main

import "fmt"

const (
	// DefaultServicePort :
	DefaultServicePort = 5150

	// DefaultServiceHost :
	DefaultServiceHost = "onos-topo"
)

// Option topo client
type Option interface {
	apply(*Options)
}

type funcOption struct {
	f func(*Options)
}

func (f funcOption) apply(options *Options) {
	f.f(options)
}

func newOption(f func(*Options)) Option {
	return funcOption{
		f: f,
	}
}

// WithOptions sets the client options
func WithOptions(opts Options) Option {
	return newOption(func(options *Options) {
		*options = opts
	})
}

// Options topo SDK options
type Options struct {
	// Service is the service options
	Service ServiceOptions
}

// ServiceOptions are the options for a service
type ServiceOptions struct {
	// Host is the service host
	Host string
	// Port is the service port
	Port int

	Insecure bool
}

// GetHost gets the service host
func (o ServiceOptions) GetHost() string {
	return o.Host
}

// GetPort gets the service port
func (o ServiceOptions) GetPort() int {
	if o.Port == 0 {
		return DefaultServicePort
	}
	return o.Port
}

// IsInsecure is topo connection secure
func (o ServiceOptions) IsInsecure() bool {
	return o.Insecure
}

// GetAddress gets the service address
func (o ServiceOptions) GetAddress() string {
	return fmt.Sprintf("%s:%d", o.GetHost(), o.GetPort())
}
