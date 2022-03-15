package subnet

import (
	"bytes"
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
	availableAddresses []net.IP
}

func (a *AddressDistribution) InitializeAddresses() {
	ips := make([]string, 0)
	currAddr := baseAddr
	lastAddr := maxAddr
	for {
		currIP := net.ParseIP(currAddr)
		currAddr = a.incrementPeerAddress(currIP).String()
		if currAddr == lastAddr {
			break
		}
		ips = append(ips, currAddr)
	}

	a.availableAddresses = make([]net.IP, len(ips))
	for i, ip := range ips {
		a.availableAddresses[i] = net.ParseIP(ip)
	}
}

func (a *AddressDistribution) GetBrokerAddress() net.IP {
	return net.ParseIP(baseAddr)
}

func (a *AddressDistribution) GetAvailableAddress() (net.IP, error) {
	if len(a.availableAddresses) == 0 {
		return nil, fmt.Errorf("no available addresses")
	}

	addr := a.availableAddresses[0]
	a.availableAddresses[0] = a.availableAddresses[len(a.availableAddresses)-1]
	a.availableAddresses = a.availableAddresses[:len(a.availableAddresses)-1]

	sort.Slice(a.availableAddresses, func(i, j int) bool {
		return bytes.Compare(a.availableAddresses[i], a.availableAddresses[j]) < 0
	})

	return addr, nil
}

func (a *AddressDistribution) ReturnAddress(addr string) {
	ip := net.ParseIP(addr)
	a.availableAddresses = append(a.availableAddresses, ip)
}

func (a *AddressDistribution) incrementPeerAddress(currIP net.IP) net.IP {
	ip := currIP.To4()

	if ip == nil {
		log.Fatalf("error parsing ip: %v", ip)
	}

	ip[3]++

	return ip
}

func AddressToNet(addr string, mask int) net.IPNet {
	ip := net.ParseIP(addr)
	return net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(mask, mask),
	}
}
