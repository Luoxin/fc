package honoka

import (
	"github.com/darabuchi/log"
)

func Load() error {
	err := LoadMeta()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = LoadConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
