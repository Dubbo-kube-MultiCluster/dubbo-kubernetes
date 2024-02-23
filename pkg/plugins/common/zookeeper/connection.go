package zookeeper

import (
	config "github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/zookeeper"
	"github.com/dubbogo/go-zookeeper/zk"
)

func ConnectToZK(cfg config.ZookeeperStoreConfig) (*zk.Conn, error) {
	connect, _, err := zk.Connect(cfg.Servers, cfg.SessionTimeout)
	if err != nil {
		return nil, err
	}
	return connect, nil
}
