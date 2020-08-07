package config

type ListenConfigParam struct {
	ServiceName string
	GroupId     string
	OnChange    func(object, key, value string)
}
