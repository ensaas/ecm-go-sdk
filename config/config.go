package config

import (
	"ecm-sdk-go/cache"
	"ecm-sdk-go/constants"
	"errors"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	IpAddr string
	Port   uint64
}

type ClientConfig struct {
	CachePath            string
	UpdateEnvWhenChanged bool
	ListenInterval       uint64
}

type Config struct {
	clientConfigValid  bool
	clientConfig       ClientConfig
	serverConfigsValid bool
	serverConfig       ServerConfig
}

func (config *Config) SetClientConfig(clientConfig ClientConfig) (err error) {

	if clientConfig.CachePath == "" {
		clientConfig.CachePath = cache.GetCurrentPath() + string(os.PathSeparator) + "cache"
	}

	log.Printf("[config.SetClientConfig] cacheDir:<%s>", clientConfig.CachePath)

	if clientConfig.ListenInterval < 5*1000 {
		clientConfig.ListenInterval = constants.ListenInterval
	}

	config.clientConfig = clientConfig
	config.clientConfigValid = true

	return
}

func (config *Config) SetServerConfig(serverConfig ServerConfig) (err error) {
	if len(serverConfig.IpAddr) < 0 || serverConfig.Port < 0 || serverConfig.Port > 65535 {
		err = errors.New("[config.SetServerConfig] server config is invalid")
		return
	}

	if len(serverConfig.IpAddr) == 0 {
		if os.Getenv(constants.ConfigServerEnvVar) == "" {
			err = errors.New("[config.SetServerConfig] server ip address is empty")
			return
		}
		serverConfig.IpAddr = os.Getenv(constants.ConfigServerEnvVar)
	}

	if serverConfig.Port == 0 {
		if os.Getenv(constants.ConfigPortEnvVar) == "" {
			err = errors.New("[config.SetServerConfig] server port is empty")
			return
		}
		serverConfig.Port, err = strconv.ParseUint(os.Getenv(constants.ConfigPortEnvVar), 10, 64)
		if err != nil {
			err = errors.New("[config.SetServerConfig] server port is invalid")
			return
		}
	}

	// TODO: check connect server

	config.serverConfig = serverConfig
	config.serverConfigsValid = true

	return
}

func (config *Config) GetClientConfig() (clientConfig ClientConfig, err error) {
	clientConfig = config.clientConfig
	if !config.clientConfigValid {
		err = errors.New("[config.GetClientConfig] invalid client config")
	}
	return
}

func (config *Config) GetServerConfig() (serverConfig ServerConfig, err error) {
	serverConfig = config.serverConfig
	if !config.serverConfigsValid {
		err = errors.New("[config.GetServerConfig] invalid server configs")
	}
	return
}
