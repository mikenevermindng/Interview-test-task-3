package api_conf

import (
	"monitor-service/conf"
)

var apiConf = (*conf.ApiConfiguration)(nil)

func LoadConfig() {
	loadedConf := conf.NewConfiguration()
	apiConf = &loadedConf.Api
}

func GetConfig() *conf.ApiConfiguration {
	return apiConf
}
