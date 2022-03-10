package addresses

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	baseAddr = "10.0.0.0"
	maxAddr  = "10.0.0.254"
	mask     = 32
)

type AddressDistribution struct {
	availableAddresses []string
}

func (a *AddressDistribution) InitializeAddresses() {
	a.availableAddresses = make([]string, 0)
	currAddr := baseAddr
	lastAddr := maxAddr

	for {
		currAddr = a.incrementPeerAddress(currAddr)
		if currAddr == lastAddr {
			break
		}

		a.availableAddresses = append(a.availableAddresses, currAddr)
	}
}

func (a *AddressDistribution) GetLastAddress() net.IP {
	return net.ParseIP(maxAddr)
}

func (a *AddressDistribution) GetAvailableAddress() (net.IP, error) {
	if len(a.availableAddresses) == 0 {
		return nil, fmt.Errorf("no available addresses")
	}

	addr := net.ParseIP(a.availableAddresses[0])
	a.availableAddresses[0] = a.availableAddresses[len(a.availableAddresses)-1]
	a.availableAddresses = a.availableAddresses[:len(a.availableAddresses)-1]

	log.Info(a.availableAddresses[0])

	return addr, nil
}

func (a *AddressDistribution) ReturnAddress(addr string) {
	a.availableAddresses = append(a.availableAddresses, addr)
}

func (a *AddressDistribution) incrementPeerAddress(currIP string) string {
	ip := net.ParseIP(currIP)
	ip = ip.To4()

	if ip == nil {
		log.Fatalf("error parsing ip: %v", ip)
	}

	ip[3]++

	return ip.String()
}

func AddressToNet(addr net.IP) net.IPNet {
	return net.IPNet{
		IP:   addr,
		Mask: net.CIDRMask(mask, mask),
	}
}
