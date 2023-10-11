package smash

import (
	"log"
	"net/http"
)

type Node struct {
	moduleManager *ModuleManager
	scheduler     *Scheduler
	bus           *LocalBus
	managmentAPI  *ManagmentAPI

	cfg NodeConfig
}

type NodeConfig struct {
	BusAddr          string
	ManagmentAPIAddr string
}

func NewNode(cfg NodeConfig) (*Node, error) {
	bus, err := NewLocalBus(cfg.BusAddr)

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

	managmentAPI, err := NewManagmentAPI(ManagmentAPIConfig{
		Bus:           bus,
		ModuleManager: moduleManager,
	})

	if err != nil {
		return nil, err
	}

	return &Node{
		moduleManager: moduleManager,
		scheduler:     scheduler,
		bus:           bus,
		managmentAPI:  managmentAPI,
		cfg:           cfg,
	}, nil
}

func (n *Node) LoadModule(modulePath string) (*Module, error) {
	return n.moduleManager.Load(modulePath)
}

func (n *Node) RegisterModules() error {
	return n.moduleManager.RegisterModules()
}

func (n *Node) Link(remoteAddr string) (string, error) {
	return n.bus.Link(&RemoteBus{
		Addr: remoteAddr,
	}, true)
}

func (n *Node) Run() error {
	if err := n.scheduler.Start(); err != nil {
		return err
	}

	go func() {
		if err := n.bus.Serve(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(n.cfg.ManagmentAPIAddr, n.managmentAPI); err != nil {
			log.Println(err)
		}
	}()

	n.scheduler.Wait()

	return nil
}
