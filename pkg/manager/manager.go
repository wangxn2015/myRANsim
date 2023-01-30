package manager

import (
	"github.com/wangxn2015/myRANsim/pkg/model"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
	"github.com/wangxn2015/onos-lib-go/pkg/northbound"
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
	server *northbound.Server
}

func (m Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Error("Unable to run Manager:", err)
	}
}

func (m Manager) Close() {

}

func (m Manager) Start() error {
	if err := model.Load(m.model, m.config.ModelName); err != nil {
		log.Error(err)
		return err
	}
	if err := m.startNorthboundSever(); err != nil {
		log.Error(err)
	}
	return nil
}

func (m Manager) startNorthboundSever() error {
	m.server = northbound.NewServer(northbound.NewServerCfg(
		m.config.CAPath,
		m.config.KeyPath,
		m.config.CertPath,
		int16(m.config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	return nil
}
