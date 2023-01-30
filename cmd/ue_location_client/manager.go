package main

import (
	"context"
	"github.com/onosproject/onos-ric-sdk-go/pkg/utils/creds"

	"github.com/onosproject/onos-ric-sdk-go/pkg/topo/connection"
	model "github.com/wangxn2015/myRANsim/api/ue_location"
	"github.com/wangxn2015/onos-lib-go/pkg/grpc/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	CAPath   string
	KeyPath  string
	CertPath string
	GRPCPort int
}

func NewManager(cfg Config) *Manager {

	client, err := NewClient()
	if err != nil {
		log.Error("create client error: ", err)
	}

	return &Manager{
		config: cfg,
		client: client,
	}
}

type Manager struct {
	config Config
	client Client
}

func (m Manager) Run() {
	err := m.start()
	if err != nil {
		log.Error("error when running manager: %v", err)
	}
}

func (m Manager) start() error {

	return nil
}

//---------------------------------------------------------------------------------
//---------------------------------------------------------------------------------
type Client interface {
	GetUE(ctx context.Context) error
}

// NewClient creates a new topo client
func NewClient(opts ...Option) (Client, error) {
	clientOptions := Options{}

	for _, opt := range opts {
		opt.apply(&clientOptions)
	}

	if clientOptions.Service.Host == "" || clientOptions.Service.Port == 0 {
		clientOptions.Service.Host = DefaultServiceHost
		clientOptions.Service.Port = DefaultServicePort
	}

	dialOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor()),
	}
	if clientOptions.Service.Insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		tlsConfig, err := creds.GetClientCredentials()
		if err != nil {
			log.Warn(err)
			return nil, err
		}

		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}
	conns := connection.NewManager()
	conn, err := conns.Connect(clientOptions.Service.GetAddress(), dialOpts...)
	if err != nil {
		log.Error("client connection error", err)
		return nil, err
	}

	cl := model.NewUeLocationServiceClient(conn)
	log.Info("connection established")

	return &topo{
		client: cl,
	}, nil
}

type topo struct {
	client model.UeLocationServiceClient
}

func (t topo) GetUE(ctx context.Context) error {

	return nil
}
