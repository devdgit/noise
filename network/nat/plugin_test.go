package nat

import (
	"testing"
	"time"

	"github.com/perlin-network/noise/network"
	"github.com/perlin-network/noise/network/discovery"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPlugin(t *testing.T) {
	t.Parallel()

	b := network.NewBuilder()
	RegisterPlugin(b)
	n, err := b.Build()
	assert.Equal(t, nil, err)
	p, ok := n.Plugins.Get(PluginID)
	assert.Equal(t, true, ok)
	natPlugin := p.(*plugin)
	assert.NotEqual(t, nil, natPlugin)
}

func TestNatConnect(t *testing.T) {
	t.Parallel()

	numNodes := 2
	nodes := make([]*network.Network, 0)
	for i := 0; i < numNodes; i++ {
		b := network.NewBuilder()
		port := network.GetRandomUnusedPort()
		b.SetAddress(network.FormatAddress("tcp", "localhost", uint16(port)))
		RegisterPlugin(b)
		b.AddPlugin(new(discovery.Plugin))
		n, err := b.Build()
		go n.Listen()

		assert.Equal(t, nil, err)
		pInt, ok := n.Plugins.Get(PluginID)
		assert.Equal(t, true, ok)
		p := pInt.(*plugin)
		assert.NotEqual(t, nil, p)
		nodes = append(nodes, n)
		n.BlockUntilListening()
	}

	nodes[1].Bootstrap(nodes[0].Address)
	pluginInt, ok := nodes[1].Plugin(discovery.PluginID)
	assert.Equal(t, true, ok)
	plugin := pluginInt.(*discovery.Plugin)
	routes := plugin.Routes
	peers := routes.GetPeers()
	for len(peers) < numNodes-1 {
		peers = routes.GetPeers()
		time.Sleep(50 * time.Millisecond)
	}

	assert.Equal(t, len(peers), 1)
}
