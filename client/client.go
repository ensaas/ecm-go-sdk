package client

import (
	"ecm-sdk-go/config"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/types"
	"ecm-sdk-go/utils"
	"errors"
	"log"

	"k8s.io/apimachinery/pkg/util/json"
)

type ConfigClient struct {
	serviceConfig map[string]*configproto.Config
	grpcClient    *GrpcClient
}

func NewConfigClient(config *config.Config) (ConfigClient, error) {
	client := ConfigClient{}
	// init service config
	client.serviceConfig = map[string]*configproto.Config{}

	clientConfig, err := config.GetClientConfig()
	if err != nil {
		return client, err
	}

	// get Grpc Client
	grpcClient, err := newGrpcClient(clientConfig)
	if err != nil {
		log.Printf("[client.client] grpc server cannot be connected %s", err.Error())
		return client, err
	}
	client.grpcClient = grpcClient

	return client, err
}

func (client *ConfigClient) DeleteConfigClient() {
	if client.grpcClient != nil {
		client.grpcClient.deleteGrpcClient()
	}
}

func (client *ConfigClient) GetConfig(appGroupName, configName string) (*types.Config, error) {
	// check service name and group id
	if appGroupName == "" {
		var err error
		appGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return nil, err
		}
		if appGroupName == "" {
			return nil, errors.New("[client.GetConfig] the app group name can not be empty")
		}
	}

	if configName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return nil, err
		}
		if len(configNames) == 1 {
			configName = configNames[0]
		}
		if configName == "" {
			return nil, errors.New("[client.GetConfig] the config name can not be empty")
		}
	}

	serviceKey := utils.GetServiceConfigKey(appGroupName, configName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(appGroupName, configName, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("[client.GetConfig] grpc server can not be connected")
	}

	// json unmarsh services
	services := map[string]map[string]*types.ServiceAddress{}
	if client.serviceConfig[serviceKey].Services != "" {
		if err := json.Unmarshal([]byte(client.serviceConfig[serviceKey].Services), &services); err != nil {
			return nil, errors.New("[client.GetConfig] JSON unmarshal services failed")
		}
	}
	config := &types.Config{
		Private:       client.serviceConfig[serviceKey].Private,
		Version:       client.serviceConfig[serviceKey].Version,
		Format:        client.serviceConfig[serviceKey].Format,
		Public:        client.serviceConfig[serviceKey].Public,
		PublicVersion: client.serviceConfig[serviceKey].PublicVersion,
		PublicFormat:  client.serviceConfig[serviceKey].PublicFormat,
		Services:      services,
	}

	return config, nil
}

func (client *ConfigClient) GetKeyValueConfig(appGroupName, configName string) (*types.KeyValueConfig, error) {
	// check service name and group id
	if appGroupName == "" {
		var err error
		appGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return nil, err
		}
		if appGroupName == "" {
			return nil, errors.New("[client.GetKeyValueConfig] the app group name can not be empty")
		}
	}

	if configName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return nil, err
		}
		if len(configNames) == 1 {
			configName = configNames[0]
		}
		if configName == "" {
			return nil, errors.New("[client.GetKeyValueConfig] the config name can not be empty")
		}
	}

	serviceKey := utils.GetServiceConfigKey(appGroupName, configName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(appGroupName, configName, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("[client.GetKeyValueConifg] grpc server can not be connected")
	}

	return utils.GetKeyValueConfig(client.serviceConfig[serviceKey]), nil
}

func (client *ConfigClient) GetPublicConfig(appGroupName, configName string) (string, error) {
	// check service name and group id
	if appGroupName == "" {
		var err error
		appGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return "", err
		}
		if appGroupName == "" {
			return "", errors.New("[client.GetPublicConfig] the app group name can not be empty")
		}
	}

	if configName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return "", err
		}
		if len(configNames) == 1 {
			configName = configNames[0]
		}
		if configName == "" {
			return "", errors.New("[client.GetPublicConfig] the config name can not be empty")
		}
	}

	var public string
	serviceKey := utils.GetServiceConfigKey(appGroupName, configName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}

		err := client.grpcClient.getConfig(appGroupName, configName, client.serviceConfig[serviceKey])
		if err != nil {
			return "", err
		}
		public = client.serviceConfig[serviceKey].Public
	} else {
		return "", errors.New("[client.GetPublicConfig] grpc server can not be connected")
	}

	return public, nil
}

func (client *ConfigClient) GetPrivateConfig(appGroupName, configName string) (string, error) {
	// check service name and group id
	if appGroupName == "" {
		var err error
		appGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return "", err
		}
		if appGroupName == "" {
			return "", errors.New("[client.GetPrivateConfig] the app group name can not be empty")
		}
	}

	if configName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return "", err
		}
		if len(configNames) == 1 {
			configName = configNames[0]
		}
		if configName == "" {
			return "", errors.New("[client.GetPrivateConfig] the config name can not be empty")
		}
	}

	var private string
	serviceKey := utils.GetServiceConfigKey(appGroupName, configName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(appGroupName, configName, client.serviceConfig[serviceKey])
		if err != nil {
			return "", err
		}
		private = client.serviceConfig[serviceKey].Private
	} else {
		return "", errors.New("[client.GetPrivateConfig] grpc server can not be connected")
	}

	return private, nil
}

func (client *ConfigClient) GetServiceAddress(appGroupName, configName, service string) (map[string]*types.ServiceAddress, error) {
	// check service name and group id
	if appGroupName == "" {
		var err error
		appGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return nil, err
		}
		if appGroupName == "" {
			return nil, errors.New("[client.GetServiceAddress] the app group name can not be empty")
		}
	}

	if configName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return nil, err
		}
		if len(configNames) == 1 {
			configName = configNames[0]
		}
		if configName == "" {
			return nil, errors.New("[client.GetServiceAddress] the config name can not be empty")
		}
	}

	var serviceAddress map[string]*types.ServiceAddress
	serviceKey := utils.GetServiceConfigKey(appGroupName, configName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(appGroupName, configName, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
		// json unmarsh services
		services := map[string]map[string]*types.ServiceAddress{}
		if client.serviceConfig[serviceKey].Services != "" {
			if err := json.Unmarshal([]byte(client.serviceConfig[serviceKey].Services), &services); err != nil {
				return nil, errors.New("[client.GetServiceAddress] JSON unmarshal services failed")
			}
		}
		for key, value := range services {
			if service == key {
				serviceAddress = value
				break
			}
		}
	} else {
		return nil, errors.New("[client.GetServiceAddress] grpc server can not be connected")
	}

	return serviceAddress, nil
}

func (client *ConfigClient) PublishConfig(publishConfigRequest *configproto.PublishConfigRequest) error {

	// check service name and group id
	if publishConfigRequest.AppGroupName == "" {
		var err error
		publishConfigRequest.AppGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return err
		}
		if publishConfigRequest.AppGroupName == "" {
			return errors.New("[client.PublishConfig] the app group name name can not be empty")
		}
	}

	if publishConfigRequest.ConfigName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return err
		}
		if len(configNames) == 1 {
			publishConfigRequest.ConfigName = configNames[0]
		}
		if publishConfigRequest.ConfigName == "" {
			return errors.New("[client.PublishConfig] the config name can not be empty")
		}
	}
	if client.grpcClient != nil {
		err := client.grpcClient.publishConfig(publishConfigRequest)
		if err != nil {
			return err
		}
	} else {
		return errors.New("[client.PublishConfig] grpc server can not be connected")
	}

	return nil
}

func (client *ConfigClient) ListenConfig(param config.ListenConfigParam) error {

	// check service name and group id
	if param.AppGroupName == "" {
		var err error
		param.AppGroupName, err = utils.GetDefaultAppGroupName()
		if err != nil {
			return err
		}
		if param.AppGroupName == "" {
			return errors.New("[client.ListenConfig] the app group name can not be empty")
		}
	}

	if param.ConfigName == "" {
		configNames, err := utils.GetDefaultConfigNames()
		if err != nil {
			return err
		}
		if len(configNames) == 1 {
			param.ConfigName = configNames[0]
		}
		if param.ConfigName == "" {
			return errors.New("[client.ListenConfig] the config name can not be empty")
		}
	}

	serviceKey := utils.GetServiceConfigKey(param.AppGroupName, param.ConfigName)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		client.grpcClient.listenConfig(client.serviceConfig[serviceKey], &param)
	} else {
		return errors.New("[client.ListenConfig] grpc server can not be connected")
	}

	return nil
}
