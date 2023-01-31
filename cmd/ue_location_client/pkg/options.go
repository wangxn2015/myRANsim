package pkg

import "fmt"

const (
	// DefaultServicePort :
	DefaultServicePort = 5150

	// DefaultServiceHost :
	DefaultServiceHost = "127.0.0.1"
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

// WithTopoHost sets the host for the topo service
func WithHost(host string) Option {
	return newOption(func(options *Options) {
		options.Service.Host = host
	})
}

// WithTopoPort sets the port for the topo service
func WithPort(port int) Option {
	return newOption(func(options *Options) {
		options.Service.Port = port
	})
}

func WithInsecure(is_insecure bool) Option {
	return newOption(func(options *Options) {
		options.Service.Insecure = is_insecure
	})
}

func WithCaPath(path string) Option {
	return newOption(func(options *Options) {
		options.Service.CaPath = path
	})
}

func WithCertPath(path string) Option {
	return newOption(func(options *Options) {
		options.Service.CertPath = path
	})
}

func WithKeyPath(path string) Option {
	return newOption(func(options *Options) {
		options.Service.KeyPath = path
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
	// used for testing ---- wxn, don't use it for production
	CaPath   string
	CertPath string //CA Cert path
	KeyPath  string
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
