package constants

const (
	BackendRegisterInfoPath           = "/mosn/register"
	EnvPrefix                         = "ENSAASMESH_"
	EcmServerAddrEnvVar               = EnvPrefix + "CONFIG_ADDR" // sidecar auto inject
	CachePathEnvVar                   = EnvPrefix + "CACHE_PATH"
	UpdateEnvWhenChangedEnvVar        = EnvPrefix + "UPDATE_ENV_WHEN_CHANGED"
	ListenIntervalEnvVar              = EnvPrefix + "LISTEN_INTERNAL"
	CachePath                         = "global_cache"
	CachFileName                      = "config"
	UpdateEnvWhenChanged              = true
	NotLoadCacheAtStart               = true
	ListenInterval             uint64 = 10 //unit: s
	PublicObjectName                  = "public"
	PrivateObjectName                 = "private"
	ServicesObjectName                = "services"
	GrpcResponseSuccess               = "success"
	HeartBeatPackage                  = "\n"
	HeartBeatInterval                 = 40
)
