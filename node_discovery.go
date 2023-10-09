package smash

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/mdns"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type NodeDiscovery struct {
	remoteNodes map[string]*RemoteNode
	bus         *LocalBus
}

func NewNodeDiscovery(bus *LocalBus) (*NodeDiscovery, error) {
	return &NodeDiscovery{
		remoteNodes: map[string]*RemoteNode{},
		bus:         bus,
	}, nil
}

func (nd *NodeDiscovery) registerNode(entry *mdns.ServiceEntry) error {
	remoteAddr := fmt.Sprintf("http://%s:3000", entry.AddrV4.String())
	remoteNode, err := NewRemoteNode(remoteAddr)

	if err != nil {
		return err
	}

	_, err = nd.bus.Link(remoteNode.Bus)

	if err != nil {
		return err
	}

	nd.remoteNodes[entry.Name] = remoteNode

	return nil
}

func (nd *NodeDiscovery) Run() error {
	host, err := os.Hostname()

	if err != nil {
		return err
	}

	info := []string{"My awesome service"}

	service, err := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	if err != nil {
		return err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})

	if err != nil {
		return err
	}

	defer server.Shutdown()

	entriesCh := make(chan *mdns.ServiceEntry, 4)

	defer close(entriesCh)

	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)

			if _, ok := nd.remoteNodes[entry.Name]; !ok {
				fmt.Println("Registering", nd.remoteNodes, entry.Name)

				if err := nd.registerNode(entry); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	for range time.NewTicker(time.Second * 5).C {
		mdns.Lookup("_foobar._tcp", entriesCh)
	}

	return nil
}

type RemoteNode struct {
	ID  string
	Bus *RemoteBus
}

func NewRemoteNode(remoteAddr string) (*RemoteNode, error) {
	id, err := gonanoid.New()

	if err != nil {
		return nil, err
	}

	remoteBus, err := NewRemoteBus(remoteAddr)

	if err != nil {
		return nil, err
	}

	return &RemoteNode{
		ID:  id,
		Bus: remoteBus,
	}, nil
}
