package nozomi

import (
	"errors"
	"fmt"

	"github.com/Luoxin/sexy/honoka"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/etcd"
	"github.com/elliotchance/pie/v2"
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
	NodeList []*Node
}

var currentNodeId string

func RegisterServer(node *Node) error {
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
					return err
				}
			}

			// TODO: 如果没有变化，不修改

			service.NodeList = pie.FilterNot[*Node](service.NodeList, func(n *Node) bool {
				return currentNodeId == n.Id()
			})

			service.NodeList = append(service.NodeList, node)

			err = etcd.Set(key, service)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
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
	key := fmt.Sprintf("nozomi_server_cfg_%s", honoka.ServiceName)
	err := etcd.NewMutex(fmt.Sprintf("nozomi_server_lock_%s", honoka.ServiceName)).WhenLocked(func() error {
		var service Server

		for i := 0; i < 3; i++ {
			err := etcd.GetJson(key, &service)
			if err != nil {
				if err != etcd.ErrNotFound {
					log.Errorf("err:%v", err)
					return err
				}
			}

			// TODO: 如果没有变化，不修改

			service.NodeList = pie.FilterNot[*Node](service.NodeList, func(n *Node) bool {
				return currentNodeId == n.Id()
			})

			err = etcd.Set(key, service)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
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
