package link_monitor

import (
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
)

const (
	probeCount = 3
	timeout    = time.Duration(10) * time.Second

	arangoPort           = 8529
	ARANGO_ROOT_USER     = "root"
	ARANGO_ROOT_PASSWORD = "cil_overly_v1"
	ARANGO_DATABASE      = "overly_mesh"
	ARANGO_GRAPH         = "link_performance"
)

type localPeerTopologyFunc func() []net.IP
type peerZoneFunc func(string) string

type Monitor struct {
	root   string
	client *ArangoClient
}

func (m *Monitor) InitializeMonitoring(brokerHost, rootAddr, role string) {
	host := fmt.Sprintf("http://%s:%d", brokerHost, arangoPort)
	conf := ArangoConfig{
		URL:      host,
		User:     ARANGO_ROOT_USER,
		Password: ARANGO_ROOT_PASSWORD,
		Database: ARANGO_DATABASE,
	}

	m.root = rootAddr
	m.client = NewArango(conf)

	if role == "broker" {
		m.client.CreateGraph(ARANGO_GRAPH)
	}
	m.client.LoadGraph(ARANGO_GRAPH)
}

func (m *Monitor) StartMonitor(granularity int, topoFunc localPeerTopologyFunc, zoneFunc peerZoneFunc) {
	go func() {
		for {
			m.peersProbe(topoFunc, zoneFunc)
			time.Sleep(time.Duration(granularity) * time.Second)
		}
	}()
}

func (m *Monitor) peersProbe(topoFunc localPeerTopologyFunc, zoneFunc peerZoneFunc) {
	peers := topoFunc()

	for _, peer := range peers {
		go m.probe(peer, zoneFunc)
	}
}

func (m *Monitor) probe(peer net.IP, zoneFunc peerZoneFunc) {
	pinger, err := ping.NewPinger(peer.String())
	if err != nil {
		log.Info("error creating pinger:", err)
		return
	}

	pinger.Count = probeCount
	pinger.Timeout = timeout

	err = pinger.Run()
	if err != nil {
		log.Println("failed to ping target host:", err)
	}

	stats := pinger.Statistics()

	src := fmt.Sprintf("%s--%s", m.root, zoneFunc(m.root))
	dst := fmt.Sprintf("%s--%s", peer.String(), zoneFunc(peer.String()))
	m.client.AddEdge(src, dst, stats)
}
