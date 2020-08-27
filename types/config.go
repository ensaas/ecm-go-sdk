package types

import configproto "ecm-sdk-go/proto"

type ServiceAddress struct {
	InternalAddress string `json:"internalAddress"`
	ExternalAddress string `json:"externalAddress"`
	SVCAddress      string `json:"svcAddress"`
	Port            int    `json:"port"`
	TargetPort      int    `json:"targetPort"`
}

type Config struct {
	Private       string                                `json:"private"`
	Version       string                                `json:"version"`
	Format        string                                `json:"format"`
	Public        string                                `json:"public"`
	PublicVersion string                                `json:"publicVersion"`
	PublicFormat  string                                `json:"publicFormat"`
	Services      map[string]map[string]*ServiceAddress `json:"services"`
}

type KeyValueConfig struct {
	Private       map[string]interface{} `json:"private"`
	Version       string                 `json:"version"`
	Public        map[string]interface{} `json:"public"`
	PublicVersion string                 `json:"publicVersion"`
	Services      map[string]interface{} `json:"services"`
}

type BackendRegisterResult struct {
	Token          string                   `json:"token"`
	BackendName    string                   `json:"backendName"`
	EnableTracing  bool                     `json:"enableTracing"`
	ServiceName    string                   `json:"serviceName"`
	APPGroupId     string                   `json:"appGroupId"`
	AppGroupConfig *AppGroupInBackendResult `json:"appGroupConfig"`
}

type AppGroupInBackendResult struct {
	AppGroupName string                   `json:"appGroupName"`
	Configs      []*ConfigInBackendResult `json:"configs"`
}

type ConfigInBackendResult struct {
	ConfigName string              `json:"configName"`
	WriteAble  bool                `json:"writeAble"`
	Config     *configproto.Config `json:"config"`
}
