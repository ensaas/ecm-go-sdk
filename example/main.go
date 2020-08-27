package main

import (
	"ecm-sdk-go/client"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	_ "ecm-sdk-go/global"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func createConfigClientTest() client.ConfigClient {
	var clientConfigTest = config.ClientConfig{
		CachePath:            "cache",
		ListenInterval:       10,
		UpdateEnvWhenChanged: true,
		EcmServerHost:        os.Getenv(constants.EcmServerHostEnvVar),
	}

	conf := config.Config{}
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

	// get client config
	appGroupName, err := utils.GetDefaultAppGroupName()
	if err != nil {
		fmt.Println(fmt.Sprintf("[global.init] Get app group name failed, errMessage = %s", err.Error()))
		return
	}
	configNames, err := utils.GetDefaultConfigName()
	if err != nil {
		fmt.Println(fmt.Sprintf("[global.init] Get config names failed, errMessage = %s", err.Error()))
		return
	}

	for _, configNameTmp := range configNames {

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
