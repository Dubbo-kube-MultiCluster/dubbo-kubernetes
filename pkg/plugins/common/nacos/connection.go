package nacos

import (
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	config "github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/nacos"
)

func ConnectToNacos(cfg config.NacosStoreConfig) (*zk.Conn, error) {
	return nil, nil
}
