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
	availableAddresses []net.IP
}

func (a *AddressDistribution) InitializeAddresses() {
	a.availableAddresses = make([]net.IP, 254)
	currAddr := net.ParseIP(baseAddr)
	lastAddr := net.ParseIP(maxAddr)

	for {
		a.availableAddresses = append(a.availableAddresses, currAddr)
		currAddr = a.incrementPeerAddress(currAddr)

		if currAddr.Equal(lastAddr) {
			break
		}
	}
}

func (a *AddressDistribution) GetLastAddress() net.IP {
	return net.ParseIP(maxAddr)
}

func (a *AddressDistribution) GetAvailableAddress() (net.IP, error) {
	if len(a.availableAddresses) == 0 {
		return nil, fmt.Errorf("no available addresses")
	}

	addr := a.availableAddresses[0]
	a.availableAddresses[0] = a.availableAddresses[len(a.availableAddresses)-1]
	a.availableAddresses = a.availableAddresses[:len(a.availableAddresses)-1]

	return addr, nil
}

func (a *AddressDistribution) ReturnAddress(addr net.IP) {
	a.availableAddresses = append(a.availableAddresses, addr)
}

func (a *AddressDistribution) incrementPeerAddress(currIP net.IP) net.IP {
	ip := currIP.To4()

	if ip == nil {
		log.Fatalf("error parsing ip: %v", ip)
	}

	ip[3]++

	return ip
}

func AddressToNet(addr net.IP) net.IPNet {
	return net.IPNet{
		IP:   addr,
		Mask: net.CIDRMask(mask, mask),
	}
}
