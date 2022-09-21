package honoka

import (
	"path/filepath"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"gopkg.in/yaml.v3"
)

type meta struct {
	EtcdAddrList []string `json:"etcd_addr_list,omitempty" yaml:"etcd_addr_list,omitempty" toml:"etcd_addr_list,omitempty"`
}

func (p *meta) Defaults() {
	if len(p.EtcdAddrList) == 0 {
		p.EtcdAddrList = append(p.EtcdAddrList, "http://127.0.0.1:2379")
	}
}

var (
	_meta = &meta{}

	ServiceName string
	BindIp      string
)

func LoadMeta() error {
	path := filepath.Join(utils.GetUserConfigDir(), "honoka", "meta.yaml")

	if utils.IsFile(path) {
		content, err := utils.FileRead(path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = yaml.Unmarshal(content, &_meta)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

	} else {
		log.Warnf("not found meta, use default")
	}

	_meta.Defaults()

	return nil
}

func Meta() meta {
	return *_meta
}
