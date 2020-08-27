package config

import (
	"ecm-sdk-go/cache"
	"ecm-sdk-go/constants"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

type ClientConfig struct {
	EcmServerHost        string
	CachePath            string
	UpdateEnvWhenChanged bool
	ListenInterval       uint64
}

type Config struct {
	clientConfigValid bool
	clientConfig      ClientConfig
}

func (config *Config) SetClientConfig(clientConfig ClientConfig) (err error) {

	if clientConfig.EcmServerHost != "" {
		// remove http:// or https://
		clientConfig.EcmServerHost = strings.Replace(clientConfig.EcmServerHost, "http://", "", 1)
		clientConfig.EcmServerHost = strings.Replace(clientConfig.EcmServerHost, "https://", "", 1)
		// check ecm server address and port
		ecmServerIP, ecmServerPort, err := getEcmIpAndPort(clientConfig.EcmServerHost)
		if err != nil {
			return err
		}
		if len(ecmServerIP) < 0 || ecmServerPort <= 0 || ecmServerPort > 65535 {
			return errors.New("[config.SetClientConfig] ecm server host is invalid")
		}

		if len(ecmServerIP) == 0 {
			return errors.New("[config.SetServerConfig] server ip address is empty")
		}
	} else {
		// if do not define the ecm server host, use env variale
		clientConfig.EcmServerHost = os.Getenv(constants.ConfigServerEnvVar)
	}

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

func (config *Config) GetClientConfig() (clientConfig ClientConfig, err error) {
	clientConfig = config.clientConfig
	if !config.clientConfigValid {
		err = errors.New("[config.GetClientConfig] invalid client config")
	}
	return
}

func getEcmIpAndPort(ecmServerHost string) (string, int64, error) {
	arr := strings.Split(ecmServerHost, ":")
	if len(arr) != 2 {
		return "", 0, errors.New("[config.getEcmIpAndPort] The ecm server host is invalid")
	}

	port, err := strconv.ParseInt(arr[1], 10, 0)
	if err != nil {
		return "", 0, err

	}

	return arr[0], port, nil
}
