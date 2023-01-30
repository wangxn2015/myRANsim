package main

import (
	"flag"
	"github.com/wangxn2015/onos-lib-go/pkg/certs"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("client main")

func main() {
	log.Info("start client")

	ready := make(chan bool)

	caPath := flag.String("caPath", "/home/baicells/go_project/myRANsim/cmd/ransim/.onos/config/certs/ca-cert.pem", "path to CA certificate")
	keyPath := flag.String("keyPath", "/home/baicells/go_project/myRANsim/cmd/ransim/.onos/config/certs/server-key.pem", "path to client private key")
	certPath := flag.String("certPath", "/home/baicells/go_project/myRANsim/cmd/ransim/.onos/config/server-cert.pem/", "path to client certificate")
	grpcPort := flag.Int("grpcPort", 5150, "GRPC port for e2t server")

	_, err := certs.HandleCertPaths(*caPath, *keyPath, *certPath, true)
	if err != nil {
		log.Fatal(err)
	}
	cfg := Config{
		CAPath:   *caPath,
		KeyPath:  *keyPath,
		CertPath: *certPath,
		GRPCPort: *grpcPort,
	}

	manager := NewManager(cfg)
	manager.Run()
	<-ready

}
