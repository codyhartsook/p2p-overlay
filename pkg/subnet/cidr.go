package subnet

import (
	"fmt"
	"net"
	"sort"

	log "github.com/sirupsen/logrus"
)

const (
	baseAddr = "10.0.0.1"
	maxAddr  = "10.0.0.254"
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

func (a *AddressDistribution) GetBrokerAddress() net.IP {
	return net.ParseIP(baseAddr)
}

func (a *AddressDistribution) GetAvailableAddress() (net.IP, error) {
	if len(a.availableAddresses) == 0 {
		return nil, fmt.Errorf("no available addresses")
	}

	addr := net.ParseIP(a.availableAddresses[0])
	a.availableAddresses[0] = a.availableAddresses[len(a.availableAddresses)-1]
	a.availableAddresses = a.availableAddresses[:len(a.availableAddresses)-1]

	sort.Strings(a.availableAddresses)

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

func AddressToNet(addr string, mask int) net.IPNet {
	ip := net.ParseIP(addr)
	return net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(mask, mask),
	}
}
