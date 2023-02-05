package main

import (
	"flag"
	"github.com/wangxn2015/myRANsim/pkg/manager"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
	"os/user"
)

var log = logging.GetLogger("main")

func main() {
	logging.SetLevel(logging.DebugLevel)
	log.Info("Starting RAN sim")
	ready := make(chan bool)

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("home dir:", u.HomeDir)

	//caPath := flag.String("caPath", "../ransim/.onos/config/certs/ca-cert.pem", "path to CA certificate")
	//keyPath := flag.String("keyPath", "../ransim/.onos/config/certs/server-key.pem", "path to client private key")
	//certPath := flag.String("certPath", "../ransim/.onos/config/certs/server-cert.pem", "path to client certificate")

	caPath := flag.String("caPath", u.HomeDir+"/go_project/myRANsim/cmd/ransim/.onos/config/certs/ca-cert.pem", "path to CA certificate")
	keyPath := flag.String("keyPath", u.HomeDir+"/go_project/myRANsim/cmd/ransim/.onos/config/certs/server-key.pem", "path to client private key")
	certPath := flag.String("certPath", u.HomeDir+"/go_project/myRANsim/cmd/ransim/.onos/config/certs/server-cert.pem", "path to client certificate")

	grpcPort := flag.Int("grpcPort", 5150, "GRPC port for e2t server")
	modelName := flag.String("modelName", "model/two-cell-two-node-model.yaml", "RAN sim model file")
	metricName := flag.String("metricName", "metrics", "RAN sim metric file")
	hoLogic := flag.String("hoLogic", "mho", "the location of handover logic {local,mho}")

	cfg := &manager.Config{
		CAPath:     *caPath,
		KeyPath:    *keyPath,
		CertPath:   *certPath,
		GRPCPort:   *grpcPort,
		ModelName:  *modelName,
		MetricName: *metricName,
		HOLogic:    *hoLogic,
	}
	mgr, err := manager.NewManager(cfg)
	if err == nil {
		mgr.Run()
		<-ready
		mgr.Close()
	}

}
