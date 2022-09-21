package honoka

import (
	"sync"

	"github.com/darabuchi/log"
	"go.uber.org/atomic"
)

var (
	loaded   = atomic.NewBool(false)
	loadLock sync.Mutex
)

func Load() error {
	if loaded.Load() {
		return nil
	}

	loadLock.Lock()
	defer loadLock.Unlock()
	if loaded.Load() {
		return nil
	}

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

	loaded.Store(true)
	return nil
}
