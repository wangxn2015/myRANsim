package pkg

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc/credentials"

	//"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/wangxn2015/onos-lib-go/pkg/errors"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"

	"github.com/onosproject/onos-ric-sdk-go/pkg/topo/connection"
	model "github.com/wangxn2015/myRANsim/api/ue_location"
	"github.com/wangxn2015/onos-lib-go/pkg/grpc/retry"
	"google.golang.org/grpc"
)

var log = logging.GetLogger()

type Config struct {
	CAPath   string
	KeyPath  string
	CertPath string
	GRPCPort int
	Insecure bool
}

func NewManager(cfg Config) *Manager {
	//client, err := NewClient(WithInsecure(cfg.Insecure))
	client, err := NewClient(
		WithHost(DefaultServiceHost),
		WithPort(cfg.GRPCPort),
		WithInsecure(cfg.Insecure),
		WithKeyPath(cfg.KeyPath),
		WithCaPath(cfg.CAPath),
		WithCertPath(cfg.CertPath),
	)

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
		log.Info("using default host and port")
	}

	dialOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor()),
	}
	if clientOptions.Service.Insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
		log.Info("using http insecure method")
	} else {
		log.Info("using http secure method")
		//tlsConfig, err := creds.GetClientCredentials()
		//if err != nil {
		//	log.Warn(err)
		//	return nil, err
		//}
		//
		//dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

		//------------------------------------------------------------------
		//------------change here ----wxn-----------------------------------
		log.Info("clientOptions.Service.Insecure: ", clientOptions.Service.Insecure)
		log.Infof("Loading client certs: %s %s", clientOptions.Service.CertPath, clientOptions.Service.KeyPath)
		clientCerts, err := tls.LoadX509KeyPair(clientOptions.Service.CertPath, clientOptions.Service.KeyPath)
		if err != nil {
			log.Info("Error loading default certs")
		}

		//var clientCAs *x509.CertPool
		clientCAs, err := GetCertPool(clientOptions.Service.CaPath) //CAPath: CA机构
		if err != nil {
			log.Error("error here")
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{clientCerts},
			ClientCAs:          clientCAs,
			InsecureSkipVerify: true,
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		//------------------------------------------------------------
		//------------------------------------------------------------------
	}
	conns := connection.NewManager()
	log.Info("server address: ", clientOptions.Service.GetAddress())
	conn, err := conns.Connect(clientOptions.Service.GetAddress(), dialOpts...)
	if err != nil {
		log.Error("client connection error", err)
		return nil, err
	}

	cl := model.NewUeLocationServiceClient(conn)
	log.Info("UeLocationServiceClient established")

	return &topo{
		client: cl,
	}, nil
}

// GetCertPool loads the Certificate Authority from the given path
func GetCertPool(CaPath string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CaPath)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("failed to append CA certificate from %s", CaPath)
	}
	return certPool, nil
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
