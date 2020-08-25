package utils

import (
	"ecm-sdk-go/constants"
	"ecm-sdk-go/flatten"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/types"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func GetDefaultAppGroupName() string {
	appGroupName := os.Getenv(constants.AppGroupNameEnvVar)
	return appGroupName
}

func GetDefaultConfigName() string {
	configName := os.Getenv(constants.ConfigNamesEnvVar)
	return configName
}

func ParseBackendInfo(maxRetryTimes int) (string, string, error) {

	backendInfo := &struct {
		BackendName string `json:"backendName"`
		Token       string `json:"token"`
	}{}

	// wait backend register
	var content []byte
	var err error
	fileName := constants.BackendRegisterInfoPath
	for i := 0; i < maxRetryTimes; i++ {
		content, err = ioutil.ReadFile(fileName)
		if err != nil {
			if i == maxRetryTimes-1 {
				return "", "", fmt.Errorf("failed to parse backend information from file:%s, err:%s! ", fileName, err.Error())
			}
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	if err = json.Unmarshal(content, backendInfo); err != nil {
		return "", "", err
	}
	return backendInfo.BackendName, backendInfo.Token, nil
}

func ParseConfigToMap(config, format string) (map[string]interface{}, error) {

	var flattenMap map[string]interface{}
	var err error

	if config != "" {
		var mapConfig map[string]interface{}

		switch format {
		case "json":
			if err = json.Unmarshal([]byte(config), &mapConfig); err != nil {
				log.Printf("[utils.parseConfigToMap] json unmarshal failed: " + err.Error())
				return nil, err
			}
		case "yaml":
			if err = yaml.Unmarshal([]byte(config), &mapConfig); err != nil {
				log.Printf("[utils.parseConfigToMap] yaml unmarshal failed: " + err.Error())
				return nil, err
			}
		case "toml":
			if err = toml.Unmarshal([]byte(config), &mapConfig); err != nil {
				log.Printf("[utils.parseConfigToMap] toml unmarshal failed: " + err.Error())
				return nil, err
			}
		default:
			log.Printf("[utils.parseConfigToMap] unsupported format")
			return nil, errors.New("unsupported format")
		}

		flattenMap, err = flatten.Flatten(mapConfig, "", flatten.DotStyle)
		if err != nil {
			log.Printf("[utils.parseConfigToMap] flatten failed: " + err.Error())
			return nil, err
		}
	}

	return flattenMap, nil
}

func GetServiceConfigKey(appGroupName, configName string) string {
	return appGroupName + "_" + configName
}

func GetServiceConfigKeyPrefix(appGroupName, configName string) string {
	return appGroupName + "_" + configName + "_" + "keyvalue"
}

func GetKeyValueConfig(serviceConfig *configproto.Config) *types.KeyValueConfig {
	flattenPrivate, err := ParseConfigToMap(serviceConfig.Private, serviceConfig.Format)
	if err != nil {
		log.Printf("[client.grpc_client] flatten private config failed: " + err.Error())
		return nil
	}

	flattenPublic, err := ParseConfigToMap(serviceConfig.Public, serviceConfig.PublicFormat)
	if err != nil {
		log.Printf("[client.grpc_client] flatten public config failed: " + err.Error())
		return nil
	}
	flattenServices, err := ParseConfigToMap(serviceConfig.Services, "json")
	if err != nil {
		log.Printf("[client.grpc_client] flatten services config failed: " + err.Error())
		return nil
	}

	keyValueConfig := &types.KeyValueConfig{
		Private:       flattenPrivate,
		Version:       serviceConfig.Version,
		Public:        flattenPublic,
		PublicVersion: serviceConfig.PublicVersion,
		Services:      flattenServices,
	}

	return keyValueConfig
}
