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

type localPeerTopologyFunc func() ([]net.IP, error)

type Monitor struct {
	root   string
	client *ArangoClient
}

func (m *Monitor) StartMonitor(granularity int, brokerHost, rootAddr string, fun localPeerTopologyFunc) {
	host := fmt.Sprintf("http://%s:%d", brokerHost, arangoPort)
	conf := ArangoConfig{
		URL:      host,
		User:     ARANGO_ROOT_USER,
		Password: ARANGO_ROOT_PASSWORD,
		Database: ARANGO_DATABASE,
	}

	m.client = NewArango(conf)
	m.client.CreateGraph(ARANGO_GRAPH)

	m.root = rootAddr
	go func() {
		for {
			m.peersProbe(fun)
			time.Sleep(time.Duration(granularity) * time.Second)
		}
	}()
}

func (m *Monitor) peersProbe(fun localPeerTopologyFunc) {
	peers, err := fun()
	if err != nil {
		log.Errorf("error getting local peers: %v", err)
		return
	}

	for _, peer := range peers {
		go m.probe(peer)
	}
}

func (m *Monitor) probe(peer net.IP) {
	pinger, err := ping.NewPinger(peer.String())
	if err != nil {
		log.Info("error creating pinger:", err)
		return
	}

	pinger.OnFinish = m.uploadLinkPerformance

	pinger.Count = probeCount
	pinger.Timeout = timeout

	err = pinger.Run()
	if err != nil {
		log.Println("failed to ping target host:", err)
	}
}

func (m *Monitor) uploadLinkPerformance(stats *ping.Statistics) {
	m.client.AddEdge(m.root, stats.Addr, stats)
}
