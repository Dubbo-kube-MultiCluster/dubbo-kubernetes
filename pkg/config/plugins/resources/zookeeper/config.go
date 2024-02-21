package zookeeper

import (
	"github.com/apache/dubbo-kubernetes/pkg/config"
	"time"
)

type ZookeeperStoreConfig struct {
	config.BaseConfig
	Servers        []string      `json:"servers" envconfig:"dubbo_store_zookeeper_servers"`
	SessionTimeout time.Duration `json:"sessionTimeout" envconfig:"dubbo_store_zookeeper_session_timeout"`
}
