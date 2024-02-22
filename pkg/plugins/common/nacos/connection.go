package nacos

import (
	config "github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/nacos"
	"github.com/dubbogo/go-zookeeper/zk"
)

func ConnectToNacos(cfg config.NacosStoreConfig) (*zk.Conn, error) {
	return nil, nil
}
