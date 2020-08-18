package constants

const (
	InitWaitTimeout                   = 10 // second
	BackendRegisterInfoPath           = "/mosn/register"
	EnvPrefix                         = "ENSAASMESH_"
	ConfigServerEnvVar                = EnvPrefix + "CONFIG_SERVER"
	ConfigPortEnvVar                  = EnvPrefix + "CONFIG_PORT"
	AppGroupNameEnvVar                = EnvPrefix + "APPGROUP_NAME"
	ConfigNameEnvVar                  = EnvPrefix + "CONFIG_NAME"
	CachePathEnvVar                   = EnvPrefix + "CACHE_PATH"
	UpdateEnvWhenChangedEnvVar        = EnvPrefix + "UPDATE_ENV_WHEN_CHANGED"
	NotLoadCacheAtStartEnvVar         = EnvPrefix + "NOTLOAD_CACHE_AT_START"
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
