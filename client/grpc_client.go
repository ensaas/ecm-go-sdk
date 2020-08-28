package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"time"

	"ecm-sdk-go/cache"
	"ecm-sdk-go/config"
	"ecm-sdk-go/constants"
	configproto "ecm-sdk-go/proto"
	"ecm-sdk-go/utils"
	util "ecm-sdk-go/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcClient struct {
	ecmServerHost      string
	config             config.ClientConfig
	client             configproto.ConfigServiceClient
	ctx                context.Context
	listenConfigClient map[string]configproto.ConfigService_ListenConfigClient
	putConfigClient    map[string]configproto.ConfigService_PutConfigClient
	conn               *grpc.ClientConn
	serviceConfigMutex sync.RWMutex
	streamClientMutex  sync.RWMutex
	cancel             context.CancelFunc
	deleteChan         chan int
	listenRecvChan     chan int
	listenSendChan     chan int
	putRecvChan        chan int
	putSendChan        chan int
	chanCount          int
}

func newGrpcClient(clientConfig config.ClientConfig) (*GrpcClient, error) {

	ecmServerHost := clientConfig.EcmServerHost
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// use custom credential
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial(ecmServerHost, opts...)
	if err != nil {
		return nil, err
	}

	c := configproto.NewConfigServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	return &GrpcClient{
		ecmServerHost:      ecmServerHost,
		config:             clientConfig,
		conn:               conn,
		client:             c,
		ctx:                ctx,
		cancel:             cancel,
		listenConfigClient: make(map[string]configproto.ConfigService_ListenConfigClient),
		putConfigClient:    make(map[string]configproto.ConfigService_PutConfigClient),
		chanCount:          0,
	}, nil

}

func (c *GrpcClient) deleteGrpcClient() {

	// stop send and recv thread
	if c.listenRecvChan != nil {
		c.listenRecvChan <- 1
	}

	if c.listenSendChan != nil {
		c.listenSendChan <- 1
	}

	if c.putSendChan != nil {
		c.putSendChan <- 1
	}

	if c.putRecvChan != nil {
		c.putRecvChan <- 1
	}

	if c.deleteChan != nil {
		for i := 0; i < c.chanCount; i++ {
			select {
			case <-c.deleteChan:
				log.Printf("[client.grpc_client] stop signal")
			}
		}
	}

	c.closeStreamClient()
}

func (c *GrpcClient) reconnect() {
	// stop grpc client before reconnect
	c.closeStreamClient()

	interval := time.Second
	for {
		log.Printf("[client.grpc_client] Reconnect")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		// use custom credential
		opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
		conn, err := grpc.Dial(c.ecmServerHost, opts...)
		if err != nil {
			time.Sleep(interval + time.Duration(rand.Intn(1000))*time.Millisecond)
			interval = computeInterval(interval)
			continue
		}

		client := configproto.NewConfigServiceClient(conn)

		ctx, cancel := context.WithCancel(context.Background())

		isFailed := false
		c.streamClientMutex.Lock()
		for key := range c.listenConfigClient {
			if _, ok := c.listenConfigClient[key]; ok {
				if c.listenConfigClient[key], err = client.ListenConfig(ctx); err != nil {
					isFailed = true
					break
				}
			}
		}
		c.streamClientMutex.Unlock()
		if isFailed {
			time.Sleep(interval + time.Duration(rand.Intn(1000))*time.Millisecond)
			interval = computeInterval(interval)
			continue
		}

		c.streamClientMutex.Lock()
		for key := range c.putConfigClient {
			if _, ok := c.putConfigClient[key]; ok {
				if c.putConfigClient[key], err = client.PutConfig(ctx); err != nil {
					isFailed = true
					break
				}
			}
		}
		c.streamClientMutex.Unlock()
		if isFailed {
			time.Sleep(interval + time.Duration(rand.Intn(1000))*time.Millisecond)
			interval = computeInterval(interval)
			continue
		}

		c.client = client
		c.ctx = ctx
		c.conn = conn
		c.cancel = cancel
		log.Printf("[client.grpc_client] Connected")
		break
	}
}

func (c *GrpcClient) closeStreamClient() {
	if c == nil {
		return
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	c.streamClientMutex.Lock()
	for key := range c.listenConfigClient {
		if _, ok := c.listenConfigClient[key]; ok {
			c.listenConfigClient[key] = nil
		}
	}
	for key := range c.putConfigClient {
		if _, ok := c.putConfigClient[key]; ok {
			c.putConfigClient[key] = nil
		}
	}
	c.streamClientMutex.Unlock()

	c.cancel()

}

func (c *GrpcClient) getConfig(appGroupName, configName string, serviceConfig *configproto.Config) error {

	// send rpc
	c.serviceConfigMutex.RLock()
	data, err := c.client.GetConfig(c.ctx, &configproto.ConfigVersion{
		Version:       serviceConfig.Version,
		AppGroupName:  appGroupName,
		ConfigName:    configName,
		PublicVersion: serviceConfig.PublicVersion,
	})
	c.serviceConfigMutex.RUnlock()

	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == codes.NotFound {
			// write empty string to cache file
			cache.WriteConfigToCache(c.config.CachePath, appGroupName, configName, &configproto.Config{})
			log.Printf("[client.getConfig] " + errStatus.Message())
			return err
		} else if errStatus.Code() == codes.Internal || errStatus.Code() == codes.Unavailable {
			// get config from cache
			data, err = cache.ReadConfigFromCache(c.config.CachePath, appGroupName, configName)
			if err != nil {
				log.Printf("[ERROR] get config from cache  error:%s ", err.Error())
				return errors.New("read config from both server and cache fail")
			}
		} else {
			log.Printf("[client.getConfig] " + err.Error())
			return err
		}
	}

	if data != nil && !reflect.DeepEqual(data, &configproto.Config{}) {
		c.serviceConfigMutex.Lock()

		// update service config and set env
		if err := c.updateServiceConfig(serviceConfig, data, nil); err != nil {
			return err
		}

		// write config to cache file
		cache.WriteConfigToCache(c.config.CachePath, appGroupName, configName, serviceConfig)
		c.serviceConfigMutex.Unlock()
	}

	return nil
}

func (c *GrpcClient) publishConfig(publishConfigRequest *configproto.PublishConfigRequest) error {

	response, err := c.client.PublishConfig(c.ctx, publishConfigRequest)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == codes.Unavailable {
			// retry send rpc
			c.reconnect()
			response, err = c.client.PublishConfig(c.ctx, publishConfigRequest)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if response.Result != constants.GrpcResponseSuccess {
		return errors.New("Publish config failed: " + response.Result)
	}

	return nil
}

func (c *GrpcClient) listenConfig(serviceConfig *configproto.Config, param *config.ListenConfigParam) {
	// initial delete channel
	if c.deleteChan == nil {
		c.deleteChan = make(chan int)
	}

	// initial receive channel
	if c.listenRecvChan == nil {
		c.listenRecvChan = make(chan int)
	}
	c.chanCount++

	// create listen config client
	listenClientKey := utils.GetServiceConfigKeyAddRandom(param.AppGroupName, param.ConfigName)
	c.streamClientMutex.Lock()
	listenConfigClient, err := c.client.ListenConfig(c.ctx)
	if err != nil {
		log.Printf(err.Error())
		return

	}
	c.listenConfigClient[listenClientKey] = listenConfigClient
	c.streamClientMutex.Unlock()

	go func() {
		for {
			select {
			case <-c.listenRecvChan:
				log.Printf("[client.listenConfig] listen receive thread receive graceful shutdown signal")
				c.deleteChan <- c.chanCount
				return
			default:
				c.streamClientMutex.RLock()
				listenConfigClient := c.listenConfigClient[listenClientKey]
				c.streamClientMutex.RUnlock()

				if listenConfigClient == nil {
					time.Sleep(time.Second)
					continue
				}
				data, err := listenConfigClient.Recv()
				if err != nil {
					if err == io.EOF {
						log.Printf("[client.listenConfig] listen receive thread failed: " + err.Error())
						return
					}
					errStatus, _ := status.FromError(err)
					if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
						errStatus.Code() == codes.Canceled {
						log.Printf("[client.listenConfig] listen receive thread failed: " + err.Error())
						return
					}
					log.Printf("[client.listenConfig] listen receive thread failed: " + err.Error())
					time.Sleep(time.Second)
					continue
				}

				c.serviceConfigMutex.Lock()
				// update service config and set env
				if err := c.updateServiceConfig(serviceConfig, data, param); err != nil {
					c.serviceConfigMutex.Unlock()
					continue
				}

				// write config to cache file
				if param != nil {
					cache.WriteConfigToCache(c.config.CachePath, param.AppGroupName, param.ConfigName, serviceConfig)
				}
				c.serviceConfigMutex.Unlock()
			}
		}
	}()

	// initial send channel
	if c.listenSendChan == nil {
		c.listenSendChan = make(chan int)
	}
	c.chanCount++
	go func() {
		c.streamClientMutex.RLock()
		listenConfigClient := c.listenConfigClient[listenClientKey]
		c.streamClientMutex.RUnlock()

		c.serviceConfigMutex.RLock()
		err := listenConfigClient.Send(&configproto.ConfigVersion{
			Version:       serviceConfig.Version,
			AppGroupName:  param.AppGroupName,
			ConfigName:    param.ConfigName,
			PublicVersion: serviceConfig.PublicVersion,
		})
		c.serviceConfigMutex.RUnlock()
		if err != nil {
			if err == io.EOF {
				log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
				return
			}
			errStatus, _ := status.FromError(err)
			if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
				errStatus.Code() == codes.Canceled {
				log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
				return
			}
			log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
			c.reconnect()
		}

		t1 := time.NewTimer(time.Duration(c.config.ListenInterval) * time.Second)
		for {
			select {
			case <-c.listenSendChan:
				log.Printf("[client.listenConfig] listen send thread receive graceful shutdown signal")
				c.deleteChan <- c.chanCount
				return
			case <-t1.C:
				c.streamClientMutex.RLock()
				listenConfigClient := c.listenConfigClient[listenClientKey]
				c.streamClientMutex.RUnlock()
				if listenConfigClient == nil {
					time.Sleep(time.Second)
					continue
				}

				c.serviceConfigMutex.RLock()
				err := listenConfigClient.Send(&configproto.ConfigVersion{
					Version:       serviceConfig.Version,
					AppGroupName:  param.AppGroupName,
					ConfigName:    param.ConfigName,
					PublicVersion: serviceConfig.PublicVersion,
				})
				c.serviceConfigMutex.RUnlock()
				if err != nil {
					if err == io.EOF {
						log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
						return
					}
					errStatus, _ := status.FromError(err)
					if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
						errStatus.Code() == codes.Canceled {
						log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
						return
					}
					log.Printf("[client.listenConfig] listen send thread failed: " + err.Error())
					c.reconnect()
				}
				t1.Reset(time.Duration(c.config.ListenInterval) * time.Second)
			}
		}
	}()

	if c.putSendChan == nil {
		c.putSendChan = make(chan int)
	}
	c.chanCount++

	// create put config client
	c.streamClientMutex.Lock()
	putConfigClient, err := c.client.PutConfig(c.ctx)
	if err != nil {
		log.Printf(err.Error())
	}
	c.putConfigClient[listenClientKey] = putConfigClient
	c.streamClientMutex.Unlock()

	go func() {
		putConfigRequest := &configproto.PutConfigRequest{
			AppGroupName: param.AppGroupName,
			ConfigName:   param.ConfigName,
		}
		c.streamClientMutex.RLock()
		putConfigClient := c.putConfigClient[listenClientKey]
		c.streamClientMutex.RUnlock()

		err := putConfigClient.Send(putConfigRequest)
		if err != nil {
			if err == io.EOF {
				log.Printf("[client.listenConfig] put send thread failed: " + err.Error())
				return
			}
			errStatus, _ := status.FromError(err)
			if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
				errStatus.Code() == codes.Canceled {
				log.Printf("[client.listenConfig] put send thread failed: " + err.Error())
				return
			}
			log.Printf("[client.listenConfig] put send thread failed: " + err.Error())
			c.reconnect()
		}

		// send heartbeat package
		t1 := time.NewTimer(time.Duration(constants.HeartBeatInterval) * time.Second)
		for {
			select {
			case <-c.putSendChan:
				log.Printf("[client.listenConfig] put send thread receive graceful shutdown signal")
				c.deleteChan <- c.chanCount
				return
			case <-t1.C:
				c.streamClientMutex.RLock()
				putConfigClient := c.putConfigClient[listenClientKey]
				c.streamClientMutex.RUnlock()

				if putConfigClient == nil {
					t1.Reset(time.Duration(constants.HeartBeatInterval) * time.Second)
					continue
				}
				putConfigRequest := &configproto.PutConfigRequest{
					AppGroupName:     param.AppGroupName,
					ConfigName:       param.ConfigName,
					HeartBeatPackage: constants.HeartBeatPackage,
				}
				err := putConfigClient.Send(putConfigRequest)
				if err != nil {
					if err == io.EOF {
						log.Printf("[client.listenConfig] put send thread failed: " + err.Error())
						return
					}
					errStatus, _ := status.FromError(err)
					if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
						errStatus.Code() == codes.Canceled {
						log.Printf("[client.listenConfig] put send thread failed: " + err.Error())
						return
					}
					log.Printf("[client.listenconfig] put send thread failed: " + err.Error())
				}
				t1.Reset(time.Duration(constants.HeartBeatInterval) * time.Second)
			}
		}
	}()

	// initial send channel
	if c.putRecvChan == nil {
		c.putRecvChan = make(chan int)
	}
	c.chanCount++
	go func() {
		for {
			select {
			case <-c.putRecvChan:
				log.Printf("[client.listenConfig] put receive thread receive graceful shutdown signal")
				c.deleteChan <- c.chanCount
				return
			default:
				c.streamClientMutex.RLock()
				putConfigClient := c.putConfigClient[listenClientKey]
				c.streamClientMutex.RUnlock()

				if putConfigClient == nil {
					time.Sleep(time.Second)
					continue
				}

				data, err := putConfigClient.Recv()
				if err != nil {
					if err == io.EOF {
						log.Printf("[client.listenConfig] put receive thread failed: " + err.Error())
						return
					}
					errStatus, _ := status.FromError(err)
					if errStatus.Code() == codes.NotFound || errStatus.Code() == codes.PermissionDenied ||
						errStatus.Code() == codes.Canceled {
						log.Printf("[client.listenConfig] put receive thread failed: " + err.Error())
						return
					}
					log.Printf("[client.listenconfig] receive data from put config request failed: " + err.Error())
					time.Sleep(time.Second)
					continue
				}
				if data == nil {
					log.Printf("[client.listenconfig] receive data from put config request is empty")
					continue
				}

				// delete message of config server
				deleteMessageRequest := &configproto.UpdateConfigMessage{
					Key:   data.UpdateConfigMessage.Key,
					Value: data.UpdateConfigMessage.Value,
				}

				response, err := c.client.DeleteMessage(c.ctx, deleteMessageRequest)
				if err != nil || response.Result != constants.GrpcResponseSuccess {
					// retry
					c.client.DeleteMessage(c.ctx, deleteMessageRequest)
				}

				c.serviceConfigMutex.Lock()
				// update service config

				if err := c.updateServiceConfig(serviceConfig, data.Config, param); err != nil {
					c.serviceConfigMutex.Unlock()
					continue
				}

				// write config to cache file
				if param != nil {
					cache.WriteConfigToCache(c.config.CachePath, param.AppGroupName, param.ConfigName, serviceConfig)
				}
				c.serviceConfigMutex.Unlock()
			}
		}
	}()

}

func computeInterval(t time.Duration) time.Duration {
	return t * 2
}

func (c *GrpcClient) updateServiceConfig(serviceConfig, changedConfig *configproto.Config, param *config.ListenConfigParam) error {
	// update public
	// check changed and added keys
	if changedConfig.PublicVersion != "" {
		changedPublic := make(map[string]interface{})
		if changedConfig.Public != "" {
			var err error
			changedPublic, err = util.ParseConfigToMap(changedConfig.Public, changedConfig.PublicFormat)
			if err != nil {
				return err
			}
		}

		public := make(map[string]interface{})
		if serviceConfig.Public != "" {
			var err error
			public, err = util.ParseConfigToMap(serviceConfig.Public, serviceConfig.PublicFormat)
			if err != nil {
				return err
			}
		}
		for key, value := range changedPublic {
			if public[key] != changedPublic[key] {
				// call onChange function
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.PublicObjectName, key, fmt.Sprintf("%v", value))
				}

				// set public env
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, fmt.Sprintf("%v", value))
				}
			}
		}

		// check deleted keys
		for key := range public {
			if _, ok := changedPublic[key]; !ok {
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.PublicObjectName, key, "")
				}
				// set public env
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, "")
				}
			}
		}
		serviceConfig.Public = changedConfig.Public
		serviceConfig.PublicVersion = changedConfig.PublicVersion
		serviceConfig.PublicFormat = changedConfig.PublicFormat
	}

	// update private
	if changedConfig.Version != "" {
		changedPrivate := make(map[string]interface{})
		if changedConfig.Private != "" {
			var err error
			changedPrivate, err = util.ParseConfigToMap(changedConfig.Private, changedConfig.Format)
			if err != nil {
				return err
			}
		}

		private := make(map[string]interface{})
		if serviceConfig.Private != "" {
			var err error
			private, err = util.ParseConfigToMap(serviceConfig.Private, serviceConfig.Format)
			if err != nil {
				return err
			}
		}

		for key, value := range changedPrivate {
			if private[key] != changedPrivate[key] {
				// call onChange function
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.PrivateObjectName, key, fmt.Sprintf("%v", value))
				}

				// set private env
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, fmt.Sprintf("%v", value))
				}
			}
		}

		// check deleted keys
		for key := range private {
			if _, ok := changedPrivate[key]; !ok {
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.PrivateObjectName, key, "")
				}
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, "")
				}
			}
		}
		serviceConfig.Private = changedConfig.Private
		serviceConfig.Version = changedConfig.Version
		serviceConfig.Format = changedConfig.Format
	}

	// update services
	if changedConfig.PublicVersion != "" {
		changedServices := make(map[string]interface{})
		if changedConfig.Services != "" {
			var err error
			changedServices, err = util.ParseConfigToMap(changedConfig.Services, "json")
			if err != nil {
				return err
			}
		}

		services := make(map[string]interface{})
		if serviceConfig.Services != "" {
			var err error
			services, err = util.ParseConfigToMap(serviceConfig.Services, "json")
			if err != nil {
				return err
			}
		}

		for key, value := range changedServices {
			if services[key] != changedServices[key] {
				// call onChange function
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.ServicesObjectName, key, fmt.Sprintf("%v", value))
				}

				// set private env
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, fmt.Sprintf("%v", value))
				}
			}
		}

		// check deleted keys
		for key := range services {
			if _, ok := changedServices[key]; !ok {
				if param != nil && param.OnChange != nil {
					param.OnChange(constants.ServicesObjectName, key, "")
				}
				if c.config.UpdateEnvWhenChanged {
					os.Setenv(key, "")
				}
			}
		}
		serviceConfig.Services = changedConfig.Services
		serviceConfig.PublicVersion = changedConfig.PublicVersion
	}
	return nil
}
