package global

import (
	"ecm-sdk-go/client"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	"ecm-sdk-go/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	// get client config
	appGroupName, err := utils.GetDefaultAppGroupName()
	if err != nil {
		log.Printf("[global.init] Get app group name failed, errMessage = %s", err.Error())
		return
	}
	configNames, err := utils.GetDefaultConfigNames()
	if err != nil {
		log.Printf("[global.init] Get config names failed, errMessage = %s", err.Error())
		return
	}

	if len(configNames) == 0 {
		log.Printf("[global.init] Warning the backend does not have any config name permissions")
		return
	}

	if appGroupName == "" {
		log.Printf("[global.init] Warning the app group name is empty")
		return
	}
	clientConfig := getDefaultClientConfig()
	conf := config.Config{}
	conf.SetClientConfig(clientConfig)

	configClient, err := client.NewConfigClient(&conf)
	if err != nil {
		log.Println(err)
		return
	}

	for _, configNameTmp := range configNames {
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
