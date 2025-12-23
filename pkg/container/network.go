package container

import (
	"fmt"
	"github.com/samama/firaaq/pkg/filesystem"
	"github.com/samama/firaaq/pkg/network"
	"path/filepath"
	"strconv"
)

const netnsPath = "/var/run/vessel/netns"

// SetupNetwork creates a pair of veth interfaces and moves one end into the container namespace.
// The bridge argument is retained for future flexibility but this implementation always uses firaaq0.
func (c *Container) SetupNetwork(bridge string) (filesystem.Unmounter, error) {
	nsMountTarget := filepath.Join(netnsPath, c.Digest)
	vethName := fmt.Sprintf("veth%.7s", c.Digest)
	peerName := fmt.Sprintf("P%s", vethName)
	masterName := "firaaq0"

	if err := network.SetupVirtualEthernet(vethName, peerName); err != nil {
		return nil, err
	}
	if err := network.LinkSetMaster(vethName, masterName); err != nil {
		return nil, err
	}
	unmount, err := network.MountNewNetworkNamespace(nsMountTarget)
	if err != nil {
		return unmount, err
	}
	if err := network.LinkSetNsByFile(nsMountTarget, peerName); err != nil {
		return unmount, err
	}

	// Change current network namespace to setup the veth
	unset, err := network.SetNetNSByFile(nsMountTarget)
	if err != nil {
		return unmount, nil
	}
	defer unset()

	ctrEthName := "eth0"
	ctrEthIPAddr := c.GetIP()
	if err := network.LinkRename(peerName, ctrEthName); err != nil {
		return unmount, err
	}
	if err := network.LinkAddAddr(ctrEthName, ctrEthIPAddr); err != nil {
		return unmount, err
	}
	if err := network.LinkSetup(ctrEthName); err != nil {
		return unmount, err
	}
	if err := network.LinkAddGateway(ctrEthName, "172.30.0.1"); err != nil {
		return unmount, err
	}
	if err := network.LinkSetup("lo"); err != nil {
		return unmount, err
	}

	return unmount, nil
}

func (c *Container) SetNetworkNamespace() (network.Unsetter, error) {
	netns := filepath.Join(netnsPath, c.Digest)
	return network.SetNetNSByFile(netns)
}

// GetIP derives a stable IP in the 172.30.0.0/16 range from the container digest.
// This keeps networking deterministic while avoiding collisions for the short-lived network.
func (c *Container) GetIP() string {
	a, _ := strconv.ParseInt(c.Digest[:2], 10, 64)
	b, _ := strconv.ParseInt(c.Digest[62:], 10, 64)
	return fmt.Sprintf("172.30.%d.%d/16", a, b)
}