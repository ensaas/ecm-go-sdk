package main

import (
	"ecm-sdk-go/client"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	appGroupName = os.Getenv(constants.AppGroupNameEnvVar)
	configNames  = os.Getenv(constants.ConfigNamesEnvVar)
)

func createConfigClientTest() client.ConfigClient {
	var clientConfigTest = config.ClientConfig{
		CachePath:            "cache",
		ListenInterval:       10,
		UpdateEnvWhenChanged: true,
	}

	port, err := strconv.ParseUint(os.Getenv(constants.ConfigPortEnvVar), 10, 0)
	if err != nil {
		fmt.Println("The config port of ensaas cp is invalid. errMessage = " + err.Error())
		return client.ConfigClient{}
	}

	var serverConfigTest = config.ServerConfig{
		IpAddr: os.Getenv(constants.ConfigServerEnvVar),
		Port:   port,
	}

	conf := config.Config{}
	conf.SetServerConfig(serverConfigTest)
	conf.SetClientConfig(clientConfigTest)

	client, err := client.NewConfigClient(&conf)
	if err != nil {
		fmt.Println("Create grpc failed. errMessage = " + err.Error())
	}
	return client
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	c := createConfigClientTest()
	defer c.DeleteConfigClient()

	configNameList := strings.Split(configNames, ",")
	for _, configNameTmp := range configNameList {

		configName := strings.Trim(configNameTmp, " ")
		content, err := c.GetConfig(appGroupName, configName)
		if err != nil {
			fmt.Println(fmt.Sprintf("get raw config of config '%s' fail. errMessage = %s", configName, err.Error()))
		}
		configBytes, err := json.Marshal(content)
		if err != nil {
			fmt.Println(fmt.Sprintf("json marshal of config '%s' fail. errMessage = %s ", configName, err.Error()))
		}
		fmt.Println(fmt.Sprintf("raw config of config '%s' is:", configName))
		fmt.Println(string(configBytes))

		keyValueContent, err := c.GetKeyValueConfig(appGroupName, configName)
		if err != nil {
			fmt.Println(fmt.Sprintf("get key value config of '%s' fail. errMessage = %s", configName, err.Error()))
		}
		keyValueBytes, err := json.Marshal(keyValueContent)
		if err != nil {
			fmt.Println(fmt.Sprintf("json marshal of config '%s' fail. errMessage = %s ", configName, err.Error()))
		}
		fmt.Println(fmt.Sprintf("key value config of config '%s' is:", configName))
		fmt.Println(string(keyValueBytes))

		c.ListenConfig(config.ListenConfigParam{
			AppGroupName: appGroupName,
			ConfigName:   configName,
			OnChange: func(object, key, value string) {
				fmt.Println(fmt.Sprintf("config '%s' changed object: %s, key: %s, value: %s", configName, object, key, value))
			},
		})

		// publish config
		publishConfig := &configproto.PublishConfigRequest{
			AppGroupName: appGroupName,
			ConfigName:   configName,
			Private:      "key1: val1\nfield:\n  key2: val2\n  key3: val3\nkey4: val4\n",
			TagName:      fmt.Sprintf("v0.0.1-sdk-%s", time.Now().Format("2006/01/02/15:04:05")),
			Format:       "yaml",
			Description:  "test",
		}

		err = c.PublishConfig(publishConfig)
		if err != nil {
			fmt.Println(fmt.Sprintf("publish config '%s' fail. errMessage = %s", configName, err.Error()))
		}
	}

	wg.Wait()
}
