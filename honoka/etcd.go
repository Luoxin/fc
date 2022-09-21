package honoka

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/etcd"
)

func LoadEtcd() error {
	err := Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = etcd.Connect(etcd.Config{
		Addrs: Meta().EtcdAddrList,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
