package config

type ListenConfigParam struct {
	AppGroupName string
	ConfigName   string
	OnChange     func(object, key, value string)
}
