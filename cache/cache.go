package cache

import (
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	util "ecm-sdk-go/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GetFileName(cacheDir, cacheFilePrefix string) string {
	return cacheDir + string(os.PathSeparator) + cacheFilePrefix + "_" + constants.CachFileName
}

func WriteConfigToFile(cacheDir, cacheFilePrefix, content string) {
	mkdirIfNecessary(cacheDir)
	fileName := GetFileName(cacheDir, cacheFilePrefix)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		log.Printf("[ERROR]:faild to write config  cache:%s ,value:%s ,err:%s \n", fileName, string(content), err.Error())
	}
}

func ReadConfigFromFile(cacheDir, cacheFilePrefix string) (string, error) {
	fileName := GetFileName(cacheDir, cacheFilePrefix)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read config cache file:%s,err:%s! ", fileName, err.Error()))
	}
	return string(b), nil
}

func WriteConfigToCache(cachePath, appGroupName, configName string, serviceConfig *configproto.Config) {
	// write raw config to cache
	content, err := json.Marshal(serviceConfig)
	if err != nil {
		log.Printf("[client.grpc_client] json marshal failed: " + err.Error())
		return
	}
	WriteConfigToFile(cachePath, util.GetServiceConfigKey(appGroupName, configName), string(content))

	// write key value config to cache
	keyValueConfig := util.GetKeyValueConfig(serviceConfig)
	if keyValueConfig == nil {
		return
	}

	keyContent, err := json.Marshal(keyValueConfig)
	if err != nil {
		log.Printf("[client.grpc_client] json marshal failed: " + err.Error())
		return
	}
	WriteConfigToFile(cachePath, util.GetServiceConfigKeyPrefix(appGroupName, configName), string(keyContent))
}

func ReadConfigFromCache(cachePath, appGroupName, configName string) (*configproto.Config, error) {
	content, err := ReadConfigFromFile(cachePath, util.GetServiceConfigKey(appGroupName, configName))
	if err != nil {
		return nil, err
	}

	serviceConfig := &configproto.Config{}
	if err := json.Unmarshal([]byte(content), serviceConfig); err != nil {
		log.Printf("[client.readConfigFromCache] json unmarshal failed")
		return nil, err
	}

	return serviceConfig, nil
}
