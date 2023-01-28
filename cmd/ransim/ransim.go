package main

import (
	"flag"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("main")

func main() {
	log.Info("Starting RAN sim")

	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	grpcPort := flag.Int("grpcPort", 5150, "GRPC port for e2t server")
	modelName := flag.String("modelName", "model/two-cell-two-node-model.yaml", "RAN sim model file")
	metricName := flag.String("metricName", "metrics", "RAN sim metric file")
	hoLogic := flag.String("hoLogic", "mho", "the location of handover logic {local,mho}")
}