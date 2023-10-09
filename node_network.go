package smash

import (
	"log"
	"net/http"
)

type NodeNetwork struct {
	nodeDiscover *NodeDiscovery
	managmentAPI *ManagmentAPI
}

func NewNodeNetwork(bus *LocalBus, moduleManager *ModuleManager) (*NodeNetwork, error) {
	nodeDiscovery, err := NewNodeDiscovery(bus)

	if err != nil {
		return nil, err
	}

	managmentAPI, err := NewManagmentAPI(moduleManager, bus)

	if err != nil {
		return nil, err
	}

	return &NodeNetwork{
		nodeDiscover: nodeDiscovery,
		managmentAPI: managmentAPI,
	}, nil
}

func (nn NodeNetwork) Run() error {
	go func() {
		err := nn.nodeDiscover.Run()

		if err != nil {
			log.Fatal(err)
		}
	}()

	return http.ListenAndServe(":3000", nn.managmentAPI.handler)
}
