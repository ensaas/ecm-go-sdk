package client

import (
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/utils"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	DefaultServiceConfig = configproto.Config{}
	DefaultGrpcClient    = &GrpcClient{}
)

func init() {
	// get server addrss
	configServer := os.Getenv(constants.ConfigServerEnvVar)
	configPort := os.Getenv(constants.ConfigPortEnvVar)
	if configServer == "" || configPort == "" {
		log.Printf(fmt.Sprintf("[client.init] Warning the environment variables %s or %s is empty", constants.ConfigServerEnvVar, constants.ConfigPortEnvVar))
		return
	}

	appGroupName := utils.GetDefaultAppGroupName()
	configName := utils.GetDefaultConfigName()

	if configName == "" || appGroupName == "" {
		log.Printf(fmt.Sprintf("[client.init] Warning the environment variables %s or %s is empty", constants.AppGroupNameEnvVar, constants.ConfigNameEnvVar))
		return
	}

	clientConfig := getDefaultClientConfig()
	DefaultGrpcClient, err := newGrpcClient(configServer, configPort, clientConfig)
	if err != nil {
		log.Printf("[client.init] Warning creating grpc client. errMessage = %s" + err.Error())
		return
	}

	param := &config.ListenConfigParam{
		AppGroupName: appGroupName,
		ConfigName:   configName,
	}

	DefaultGrpcClient.listenConfig(&DefaultServiceConfig, param)
}

func getDefaultClientConfig() config.ClientConfig {
	cachePath := constants.CachePath

	if os.Getenv(constants.CachePathEnvVar) != "" {
		cachePath = os.Getenv(constants.CachePathEnvVar)
	}

	var err error
	updateEnvWhenChanged := constants.UpdateEnvWhenChanged
	if os.Getenv(constants.UpdateEnvWhenChangedEnvVar) != "" {
		updateEnvWhenChanged, err = strconv.ParseBool(os.Getenv(constants.UpdateEnvWhenChangedEnvVar))
		if err != nil {
			updateEnvWhenChanged = constants.UpdateEnvWhenChanged
		}
	}
	listenInterval := constants.ListenInterval
	if os.Getenv(constants.ListenIntervalEnvVar) != "" {
		listenInterval, err = strconv.ParseUint(os.Getenv(constants.ListenIntervalEnvVar), 10, 0)
		if err != nil {
			listenInterval = constants.ListenInterval
		}
	}

	clientConfig := config.ClientConfig{
		CachePath:            cachePath,
		UpdateEnvWhenChanged: updateEnvWhenChanged,
		ListenInterval:       listenInterval,
	}

	return clientConfig
}
