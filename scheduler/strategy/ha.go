package strategy

import (
	"sort"

	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/scheduler/node"

	log "github.com/Sirupsen/logrus"
)

// HAPlacementStrategy places a container on the node with the fewest running containers.
type HAPlacementStrategy struct {
}

// Initialize a HAPlacementStrategy.
func (p *HAPlacementStrategy) Initialize() error {
	return nil
}

// Name returns the name of the strategy.
func (p *HAPlacementStrategy) Name() string {
	return "ha"
}

// RankAndSort sorts nodes based on the ha strategy applied to the container config.
func (p *HAPlacementStrategy) RankAndSort(config *cluster.ContainerConfig, nodes []*node.Node) ([]*node.Node, error) {
	// for ha, a healthy node should decrease its weight to increase its chance of being selected
	// set healthFactor to -10 to make health degree [0, 100] overpower cpu + memory (each in range [0, 100])
	const healthFactor int64 = -10
	var ContainerHAGroup string = config.Labels["com.docker.swarm.hagroup"]
	log.WithFields(log.Fields{"Label_HA": ContainerHAGroup, "Algo": "plop"}).Info("Strategy debug")

	weightedNodes, err := weighNodes(config, nodes, healthFactor)
	if err != nil {
		return nil, err
	}

	sort.Sort(weightedNodes)
	output := make([]*node.Node, len(weightedNodes))
	for i, n := range weightedNodes {
		for _, conf := range n.Node.Containers {
			if conf.Config.Labels["com.docker.swarm.hagroup"]==ContainerHAGroup && conf.Info.State.Running {
				log.WithFields(log.Fields{"Name": conf.Info.Name, "State": conf.Info.State, "Label_HA": conf.Config.Labels["com.docker.swarm.hagroup"]}).Info("Strategy debug - Found One !!!")
			}
		}
		output[i] = n.Node
	}
	return output, nil
}
