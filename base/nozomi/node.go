package nozomi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Luoxin/sexy/honoka"
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/etcd"
	"github.com/elliotchance/pie/v2"
	"github.com/fsnotify/fsnotify"
)

type NodeState uint8

const (
	NodeStateDead NodeState = iota
	NodeStateAlive
)

type Node struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`

	State NodeState `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`
}

func (n Node) Id() string {
	return n.Address
}

type Server struct {
	Name     string  `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	NodeList []*Node `json:"node_list,omitempty" yaml:"node_list,omitempty" toml:"node_list,omitempty"`
}

var (
	currentNodeId string

	ErrorServerNotFound = errors.New("server not found")

	serverMap      sync.Map
	serviceLock    sync.Mutex
	serviceWatcher *fsnotify.Watcher
)

func RegisterServer(node *Node) error {
	log.Infof("register %s", honoka.ServiceName)

	err := honoka.LoadEtcd()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	currentNodeId = node.Id()

	node.State = NodeStateAlive

	key := fmt.Sprintf("nozomi_server_cfg_%s", honoka.ServiceName)

	err = etcd.NewMutex(fmt.Sprintf("nozomi_server_lock_%s", honoka.ServiceName)).WhenLocked(func() error {
		var service Server

		for i := 0; i < 3; i++ {
			err = etcd.GetJson(key, &service)
			if err != nil {
				if err != etcd.ErrNotFound {
					log.Errorf("err:%v", err)
					continue
				}
			}

			// TODO: 如果没有变化，不修改

			service.Name = honoka.ServiceName

			service.NodeList = pie.FilterNot[*Node](service.NodeList, func(n *Node) bool {
				return currentNodeId == n.Id()
			})

			service.NodeList = append(service.NodeList, node)

			err = etcd.Set(key, service)
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}

			return nil
		}

		return errors.New("")
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func UnregisterServer() {
	log.Warnf("unregister %s", honoka.ServiceName)

	key := fmt.Sprintf("nozomi_server_cfg_%s", honoka.ServiceName)
	err := etcd.NewMutex(fmt.Sprintf("nozomi_server_lock_%s", honoka.ServiceName)).WhenLocked(func() error {
		var service Server

		for i := 0; i < 3; i++ {
			err := etcd.GetJson(key, &service)
			if err != nil {
				if err != etcd.ErrNotFound {
					log.Errorf("err:%v", err)
					continue
				}
			}

			// TODO: 如果没有变化，不修改

			service.NodeList = pie.FilterNot[*Node](service.NodeList, func(n *Node) bool {
				return currentNodeId == n.Id()
			})

			err = etcd.Set(key, service)
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}

			return nil
		}

		return errors.New("")
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

func addWatch(name string) error {
	if serviceWatcher == nil {
		var err error
		serviceWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		go func(sign chan os.Signal) {
			for {
				select {
				case <-sign:
					log.Warn("service watcher stop")
					return

				case event := <-serviceWatcher.Events:
					switch event.Op {
					case fsnotify.Create, fsnotify.Write:
					case fsnotify.Remove:

					}
				case err := <-serviceWatcher.Errors:
					log.Errorf("service watcher error:%s", err)
				}
			}
		}(utils.GetExitSign())

	}

	serviceWatcher.Add(name)

	return nil
}

func GenServerPath(serviceName string) string {
	return filepath.Join(utils.GetUserConfigDir(), "honoka", "nozomi", serviceName+".json")
}

func loadServer(name string) (*Server, error) {
	val, ok := serverMap.Load(name)
	if ok {
		return val.(*Server), nil
	}

	serviceLock.Lock()
	defer serviceLock.Unlock()

	val, ok = serverMap.Load(name)
	if ok {
		return val.(*Server), nil
	}

	cf := GenServerPath(name)
	err := addWatch(cf)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if !utils.IsFile(cf) {
		return nil, ErrorServerNotFound
	}

	var service Server
	err = sonic.Unmarshal([]byte(cf), &service)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	serverMap.Store(name, &service)

	return &service, nil
}

func GetService(name string) (*Server, error) {
	return loadServer(name)
}
