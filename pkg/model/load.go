package model

import (
	"github.com/spf13/viper"
)

const configDir = ".onos"

func ViperConfigure(configname string) {
	// Set the file type of the configurations file
	viper.SetConfigType("yaml")

	// Set the file name of the configurations file
	viper.SetConfigName(configname)

	// Set the path to look for the configurations file
	viper.AddConfigPath("./" + configDir + "/config")
	viper.AddConfigPath("./cmd/" + configDir + "/config")
	//viper.AddConfigPath("$HOME/" + configDir + "/config")
	//viper.AddConfigPath("/etc/onos/config")
	//viper.AddConfigPath(".")
	//viper.AddConfigPath("./cmd/ransim/" + configDir + "/config")
}

func Load(model *Model, modelName string) error {
	return LoadConfig(model, modelName)
}

func LoadConfig(model *Model, modelName string) error {
	var err error

	ViperConfigure(modelName)

	if err = viper.ReadInConfig(); err != nil {

		return err
	}

	err = viper.Unmarshal(model)

	//// Convert the MCC-MNC format into numeric PLMNID
	//model.PlmnID = types.PlmnIDFromString(model.Plmn)
	//
	//// initialize neighbor's Ocn value - for mlb/handover
	//for k, v := range model.Cells {
	//	v.MeasurementParams.NCellIndividualOffsets = make(map[types.NCGI]int32)
	//	for _, n := range v.Neighbors {
	//		v.MeasurementParams.NCellIndividualOffsets[n] = 0
	//	}
	//	model.Cells[k] = v
	//}

	return err

}
