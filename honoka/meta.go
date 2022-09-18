package honoka

import (
	"path/filepath"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

type meta struct {
}

var _meta = meta{}

func LoadMeta() error {
	path := filepath.Join(utils.GetUserConfigDir(), "sexy", "meta.yaml")

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
		defaults.SetDefaults(_meta)
	}

	return nil
}

func Meta() meta {
	return _meta
}