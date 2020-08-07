package client

import (
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/types"
	"ecm-sdk-go/utils"
	"errors"
	"log"
	"strconv"

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
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		return client, err
	}

	// get Grpc Client
	configServer := serverConfig.IpAddr
	configPort := strconv.FormatUint(serverConfig.Port, 10)
	grpcClient, err := newGrpcClient(configServer, configPort, clientConfig)
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

func (client *ConfigClient) GetConfig(serviceName, groupId string) (*types.Config, error) {
	// check service name and group id
	if serviceName == "" {
		serviceName = utils.GetDefaultServiceName()
		if serviceName == "" {
			return nil, errors.New("[client.GetConfig] the service name can not be empty")
		}
	}

	if groupId == "" {
		groupId = utils.GetDefaultGroupId()
		if groupId == "" {
			groupId = constants.DefaultGroupId
		}
	}

	serviceKey := getServiceConfigKey(serviceName, groupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(serviceName, groupId, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("grpc server can not be connected")
	}

	// json unmarsh services
	services := map[string]map[string]*types.ServiceAddress{}
	if client.serviceConfig[serviceKey].Services != "" {
		if err := json.Unmarshal([]byte(client.serviceConfig[serviceKey].Services), &services); err != nil {
			return nil, errors.New("JSON unmarshal services failed")
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

func (client *ConfigClient) GetKeyValueConfig(serviceName, groupId string) (*types.KeyValueConfig, error) {
	// check service name and group id
	if serviceName == "" {
		serviceName = utils.GetDefaultServiceName()
		if serviceName == "" {
			return nil, errors.New("[client.GetConfig] the service name can not be empty")
		}
	}

	if groupId == "" {
		groupId = utils.GetDefaultGroupId()
		if groupId == "" {
			groupId = constants.DefaultGroupId
		}
	}

	serviceKey := getServiceConfigKey(serviceName, groupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(serviceName, groupId, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("grpc server can not be connected")
	}

	return getKeyValueConfig(client.serviceConfig[serviceKey]), nil
}

func (client *ConfigClient) GetPublicConfig(serviceName, groupId string) (string, error) {
	// check service name and group id
	if serviceName == "" {
		serviceName = utils.GetDefaultServiceName()
		if serviceName == "" {
			return "", errors.New("[client.GetPublicConfig] the service name can not be empty")
		}
	}

	if groupId == "" {
		groupId = utils.GetDefaultGroupId()
		if groupId == "" {
			groupId = constants.DefaultGroupId
		}
	}

	var public string
	serviceKey := getServiceConfigKey(serviceName, groupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}

		err := client.grpcClient.getConfig(serviceName, groupId, client.serviceConfig[serviceKey])
		if err != nil {
			return "", err
		}
		public = client.serviceConfig[serviceKey].Public
	} else {
		return "", errors.New("grpc server can not be connected")
	}

	return public, nil
}

func (client *ConfigClient) GetPrivateConfig(serviceName, groupId string) (string, error) {
	// check service name and group id
	if serviceName == "" {
		serviceName = utils.GetDefaultServiceName()
		if serviceName == "" {
			return "", errors.New("[client.GetPrivateConfig] the service name can not be empty")
		}
	}

	if groupId == "" {
		groupId = utils.GetDefaultGroupId()
		if groupId == "" {
			groupId = constants.DefaultGroupId
		}
	}

	var private string
	serviceKey := getServiceConfigKey(serviceName, groupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(serviceName, groupId, client.serviceConfig[serviceKey])
		if err != nil {
			return "", err
		}
		private = client.serviceConfig[serviceKey].Private
	} else {
		return "", errors.New("grpc server can not be connected")
	}

	return private, nil
}

func (client *ConfigClient) GetServiceAddress(serviceName, groupId, service string) (map[string]*types.ServiceAddress, error) {
	// check service name and group id
	if serviceName == "" {
		serviceName = utils.GetDefaultServiceName()
		if serviceName == "" {
			return nil, errors.New("[client.GetServiceAddress] the service name can not be empty")
		}
	}

	if groupId == "" {
		groupId = utils.GetDefaultGroupId()
		if groupId == "" {
			groupId = constants.DefaultGroupId
		}
	}

	var serviceAddress map[string]*types.ServiceAddress
	serviceKey := getServiceConfigKey(serviceName, groupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		err := client.grpcClient.getConfig(serviceName, groupId, client.serviceConfig[serviceKey])
		if err != nil {
			return nil, err
		}
		// json unmarsh services
		services := map[string]map[string]*types.ServiceAddress{}
		if client.serviceConfig[serviceKey].Services != "" {
			if err := json.Unmarshal([]byte(client.serviceConfig[serviceKey].Services), &services); err != nil {
				return nil, errors.New("JSON unmarshal services failed")
			}
		}
		for key, value := range services {
			if service == key {
				serviceAddress = value
				break
			}
		}
	} else {
		return nil, errors.New("grpc server can not be connected")
	}

	return serviceAddress, nil
}

func (client *ConfigClient) PublishConfig(publishConfigRequest *configproto.PublishConfigRequest) error {

	// check service name and group id
	if publishConfigRequest.ServiceName == "" {
		publishConfigRequest.ServiceName = utils.GetDefaultServiceName()
		if publishConfigRequest.ServiceName == "" {
			return errors.New("[client.PublishConfig] the service name can not be empty")
		}
	}

	if publishConfigRequest.GroupId == "" {
		publishConfigRequest.GroupId = utils.GetDefaultGroupId()
		if publishConfigRequest.GroupId == "" {
			publishConfigRequest.GroupId = constants.DefaultGroupId
		}
	}
	if client.grpcClient != nil {
		err := client.grpcClient.publishConfig(publishConfigRequest)
		if err != nil {
			return err
		}
	} else {
		return errors.New("grpc server can not be connected")
	}

	return nil
}

func (client *ConfigClient) ListenConfig(param config.ListenConfigParam) error {

	// check service name and group id
	if param.ServiceName == "" {
		param.ServiceName = utils.GetDefaultServiceName()
		if param.ServiceName == "" {
			return errors.New("[client.ListenChangedConfig] the service name can not be empty")
		}
	}

	if param.GroupId == "" {
		param.GroupId = utils.GetDefaultGroupId()
		if param.GroupId == "" {
			param.GroupId = constants.DefaultGroupId
		}
	}

	serviceKey := getServiceConfigKey(param.ServiceName, param.GroupId)

	if client.grpcClient != nil {
		if client.serviceConfig[serviceKey] == nil {
			client.serviceConfig[serviceKey] = &configproto.Config{}
		}
		client.grpcClient.listenConfig(client.serviceConfig[serviceKey], &param)
	} else {
		return errors.New("grpc server can not be connected")
	}

	return nil
}
