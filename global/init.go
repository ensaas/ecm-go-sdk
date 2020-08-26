package global

import (
	"ecm-sdk-go/client"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	"ecm-sdk-go/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	// get client config
	appGroupName := utils.GetDefaultAppGroupName()
	configNames := utils.GetDefaultConfigName()

	if configNames == "" || appGroupName == "" {
		log.Printf(fmt.Sprintf("[client.init] Warning the environment variables %s or %s is empty", constants.AppGroupNameEnvVar, constants.ConfigNamesEnvVar))
		return
	}
	clientConfig := getDefaultClientConfig()

	// get server config
	configServer := os.Getenv(constants.ConfigServerEnvVar)
	configPort := os.Getenv(constants.ConfigPortEnvVar)
	if configServer == "" || configPort == "" {
		log.Printf(fmt.Sprintf("[client.init] Warning the environment variables %s or %s is empty", constants.ConfigServerEnvVar, constants.ConfigPortEnvVar))
		return
	}
	port, err := strconv.ParseUint(configPort, 10, 0)
	if err != nil {
		log.Println("The config port invalid. errMessage = " + err.Error())
		return
	}

	var serverConfig = config.ServerConfig{
		IpAddr: configServer,
		Port:   port,
	}

	conf := config.Config{}
	conf.SetServerConfig(serverConfig)
	conf.SetClientConfig(clientConfig)

	configClient, err := client.NewConfigClient(&conf)
	if err != nil {
		log.Println(err)
		return
	}

	configNamesList := parseConfigNames(configNames)
	for _, configNameTmp := range configNamesList {
		configName := strings.Trim(configNameTmp, " ")
		configClient.ListenConfig(config.ListenConfigParam{
			AppGroupName: appGroupName,
			ConfigName:   configName,
		})
	}
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

func parseConfigNames(configNames string) []string {
	arr := strings.Split(configNames, ",")
	return arr
}
