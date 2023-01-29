package manager

import (
	"github.com/wangxn2015/myRANsim/pkg/model"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger()

type Config struct {
	CAPath     string
	KeyPath    string
	CertPath   string
	GRPCPort   int
	ModelName  string
	MetricName string
	HOLogic    string
}

func NewManager(cfg *Config) (*Manager, error) {
	log.Info("Creating Manager")
	mgr := &Manager{
		config: *cfg,
		model:  &model.Model{},
	}
	return mgr, nil
}

type Manager struct {
	config Config
	model  *model.Model
}

func (m Manager) Run() {

}

func (m Manager) Close() {

}
