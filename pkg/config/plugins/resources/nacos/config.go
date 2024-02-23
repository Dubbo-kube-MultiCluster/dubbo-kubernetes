package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/common/logger"

	"gopkg.in/natefinch/lumberjack.v2"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/config"
)

type NacosStoreConfig struct {
	config.BaseConfig
	TimeoutMs            uint64                 `json:"timeoutMs" envconfig:"dubbo_store_nacos_timeout_ms"`
	ListenInterval       uint64                 `json:"listenInterval" envconfig:"dubbo_store_nacos_listen_interval"`
	BeatInterval         int64                  `json:"beatInterval" envconfig:"dubbo_store_nacos_beat_interval"`
	NamespaceId          string                 `json:"namespaceId" envconfig:"dubbo_store_nacos_namespace_id"`
	AppName              string                 `json:"appName" envconfig:"dubbo_store_nacos_app_name"`
	Endpoint             string                 `json:"endpoint" envconfig:"dubbo_store_nacos_endpoint"`
	RegionId             string                 `json:"regionId" envconfig:"dubbo_store_nacos_region_id"`
	AccessKey            string                 `json:"accessKey" envconfig:"dubbo_store_nacos_access_key"`
	SecretKey            string                 `json:"secretKey" envconfig:"dubbo_store_nacos_secret_key"`
	OpenKMS              bool                   `json:"openKMS" envconfig:"dubbo_store_nacos_open_KMS"`
	CacheDir             string                 `json:"cacheDir" envconfig:"dubbo_store_nacos_cache_dir"`
	UpdateThreadNum      int                    `json:"updateThreadNum" envconfig:"dubbo_store_nacos_update_threadNum"`
	NotLoadCacheAtStart  bool                   `json:"notLoadCacheAtStart" envconfig:"dubbo_store_nacos_notLoad_cache_at_start"`
	UpdateCacheWhenEmpty bool                   `json:"updateCacheWhenEmpty" envconfig:"dubbo_store_nacos_update_cache_when_empty"`
	Username             string                 `json:"username" envconfig:"dubbo_store_nacos_username"`
	Password             string                 `json:"password" envconfig:"dubbo_store_nacos_password"`
	LogDir               string                 `json:"logDir" envconfig:"dubbo_store_nacos_log_dir"`
	LogLevel             string                 `json:"logLevel" envconfig:"dubbo_store_nacos_log_level"`
	LogSampling          *logger.SamplingConfig `json:"logSampling" envconfig:"dubbo_store_nacos_log_sampling"`
	ContextPath          string                 `json:"contextPath" envconfig:"dubbo_store_nacos_context_path"`
	LogRollingConfig     *lumberjack.Logger     `json:"logRollingConfig" envconfig:"dubbo_store_nacos_log_rolling_config"`
	CustomLogger         logger.Logger          `json:"customLogger" envconfig:"dubbo_store_nacos_custom_logger"`
	AppendToStdout       bool                   `json:"appendToStdout" envconfig:"dubbo_store_nacos_append_to_stdoutl"`
}
