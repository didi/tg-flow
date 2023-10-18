package apollo

import (
	"go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2"
	"go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2/model"
)

type AppConfig struct {
	config *model.ConfResult
}

func NewApolloConfig(namespace string, configName string) (*AppConfig, error){
	cfg, err := apollo.GetConfig(namespace, configName)
	if err != nil {
		return nil, err
	}

	return &AppConfig{config: cfg}, nil
}

func (a *AppConfig) IsVersion(version string) bool {
	if version == a.config.GetVersion(){
		return true
	}

	return false
}

func (a *AppConfig) GetConfigs() map[string]string {
	return a.config.GetConfigs()
}