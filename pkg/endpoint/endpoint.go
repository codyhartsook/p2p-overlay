package endpoint

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type Endpoint struct{}

type EndpointSpec struct {
	Hostname  string   `json:"hostname"`
	Subnets   []string `json:"subnets"`
	PrivateIP string   `json:"private_ip"`
	PublicIP  string   `json:"public_ip"`
}

func GetLocalEndpoint() (EndpointSpec, error) {

	return EndpointSpec{}, nil
}

func GetLocalIP() net.IP {
	return GetLocalIPForDestination("8.8.8.8")
}

func GetLocalIPForDestination(dst string) net.IP {
	conn, err := net.Dial("udp", dst+":53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
