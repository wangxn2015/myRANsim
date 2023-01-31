package pkg

import (
	"context"
	"github.com/onosproject/onos-ric-sdk-go/pkg/utils/creds"
	"github.com/wangxn2015/onos-lib-go/pkg/errors"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
	"google.golang.org/grpc/status"
	"io"

	"github.com/onosproject/onos-ric-sdk-go/pkg/topo/connection"
	model "github.com/wangxn2015/myRANsim/api/ue_location"
	"github.com/wangxn2015/onos-lib-go/pkg/grpc/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logging.GetLogger()

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
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ch := make(chan model.UeInfo)
		err := m.client.GetUE(ctx, ch)
		if err != nil {
			log.Error("error getting UE: ", err)
		}
		for resp := range ch {
			log.Info("manager receive resp: %v", resp)
		}
	}()

	return nil
}

type Client interface {
	GetUE(ctx context.Context, ch chan<- model.UeInfo) error
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

func (t topo) GetUE(ctx context.Context, ch chan<- model.UeInfo) error {
	req := model.UeLocationRequest{
		ProcessureCode: 1,
		Imsi:           1,
	}
	stream, err := t.client.GetUe(ctx, &req)
	if err != nil {
		return errors.FromGRPC(err)
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF || err == context.Canceled {
				break
			}
			if err != nil {
				stat, ok := status.FromError(err)
				if ok {
					err = errors.FromStatus(stat)
					if errors.IsCanceled(err) || errors.IsTimeout(err) {
						break
					}
				}
				log.Error("An error occurred here", err)
			}
			log.Info("recv msg: %+v", resp)
			ch <- *resp
		}
	}()

	return nil
}
