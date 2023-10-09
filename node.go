package smash

import "log"

type Node struct {
	moduleManager *ModuleManager
	scheduler     *Scheduler
	bus           *LocalBus
	nodeNetwork   *NodeNetwork
}

func NewNode() (*Node, error) {
	bus, err := NewLocalBus()

	if err != nil {
		return nil, err
	}

	scheduler, err := NewScheduler()

	if err != nil {
		return nil, err
	}

	moduleManager, err := NewModuleManager(bus, scheduler)

	if err != nil {
		return nil, err
	}

	nodeNetwork, err := NewNodeNetwork(bus, moduleManager)

	if err != nil {
		return nil, err
	}

	return &Node{
		moduleManager: moduleManager,
		nodeNetwork:   nodeNetwork,
		scheduler:     scheduler,
		bus:           bus,
	}, nil
}

func (n *Node) LoadModule(modulePath string) (*Module, error) {
	return n.moduleManager.Load(modulePath)
}

func (n *Node) Run() error {
	if err := n.moduleManager.RegisterModules(); err != nil {
		return err
	}

	if err := n.scheduler.Start(); err != nil {
		return err
	}

	go func() {
		if err := n.nodeNetwork.Run(); err != nil {
			log.Println(err)
		}
	}()

	n.scheduler.Wait()

	return nil
}
