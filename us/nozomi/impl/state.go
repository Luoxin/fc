package impl

import (
	"github.com/Luoxin/sexy/honoka"
	"github.com/darabuchi/log"
)

type State struct {
}

func InitState() (err error) {
	err = honoka.LoadEtcd()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = ServiceSync()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
