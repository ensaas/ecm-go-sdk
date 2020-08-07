package main

import (
	"ecm-sdk-go/client"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	serviceName = os.Getenv(constants.ServiceNameEnvVar)
	groupId     = "default"
)

var clientConfigTest = config.ClientConfig{
	CachePath:            "cache",
	ListenInterval:       10,
	UpdateEnvWhenChanged: true,
}

var serverConfigTest = config.ServerConfig{
	IpAddr: os.Getenv(constants.ConfigServerEnvVar),
	Port:   9000,
}

func cretateConfigClientTest() client.ConfigClient {
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

	c := cretateConfigClientTest()
	defer c.DeleteConfigClient()

	content, err := c.GetConfig(serviceName, groupId)
	if err != nil {
		fmt.Println("get config fail. errMessage = " + err.Error())
	}
	configBytes, err := json.Marshal(content)
	if err != nil {
		fmt.Println("json marshal fail. errMessage = " + err.Error())
	}
	fmt.Println("raw config is:")
	fmt.Println(string(configBytes))

	keyValueContent, err := c.GetKeyValueConfig(serviceName, groupId)
	if err != nil {
		fmt.Println("get config fail. errMessage = " + err.Error())
	}
	keyValueBytes, err := json.Marshal(keyValueContent)
	if err != nil {
		fmt.Println("json marshal fail. errMessage = " + err.Error())
	}
	fmt.Println("key value config is:")
	fmt.Println(string(keyValueBytes))

	c.ListenConfig(config.ListenConfigParam{
		ServiceName: serviceName,
		GroupId:     groupId,
		OnChange: func(object, key, value string) {
			fmt.Println("config changed object:" + object + ", key:" + key + ", value:" + fmt.Sprint(value))
		},
	})

	// publish config
	publishConfig := &configproto.PublishConfigRequest{
		ServiceName: serviceName,
		GroupId:     groupId,
		Private:     "key1: val1\nfield:\n  key2: val2\n  key3: val3\nkey4: val4\n",
		TagName:     fmt.Sprintf("v0.0.1-sdk-%s", time.Now().Format("2006/01/02/15:04:05")),
		Format:      "yaml",
		Description: "test",
	}

	err = c.PublishConfig(publishConfig)
	if err != nil {
		fmt.Println("publish config fail. errMessage = " + err.Error())
	}

	wg.Wait()
}
