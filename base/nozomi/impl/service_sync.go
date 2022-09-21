package impl

import (
	"os"
	"path/filepath"

	"github.com/Luoxin/sexy/base/nozomi"
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/etcd"
)

func ServiceSync() error {
	etcd.WatchPrefix("nozomi_server_cfg_", func(event etcd.Event) {
		var service nozomi.Server
		err := sonic.Unmarshal([]byte(event.Value), &service)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		if service.Name == "" {
			return
		}

		log.Info(string(event.Value))

		cf := nozomi.GenServerPath(service.Name)

		if !utils.IsDir(filepath.Dir(cf)) {
			err = os.MkdirAll(filepath.Dir(cf), 0777)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		}

		err = utils.FileWrite(cf, event.Value)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
	return nil
}
