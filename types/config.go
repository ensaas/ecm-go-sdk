package types

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
