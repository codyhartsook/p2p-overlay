package link_monitor

import (
	"net"
	"time"

	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
)

const (
	arangoHost = "localhost"
	probeCount = 3
	timeout    = time.Duration(10) * time.Second
)

type PerformanceSpec struct {
	ExecutedAt time.Time `json:"executed_at"`
	PeerAddr   string    `json:"peer_addr"`
	LocalAddr  string    `json:"local_addr"`
	Latency    float32   `json:"latency"`
	Jitter     float32   `json:"jitter"`
	Loss       float32   `json:"loss"`
}

type localPeerTopologyFunc func() (net.IP, []net.IP, error)

type Monitor struct {
}

func (m *Monitor) StartMonitor(granularity int, fun localPeerTopologyFunc) {
	go func() {
		for {
			time.Sleep(time.Duration(granularity) * time.Second)
			m.peersProbe(fun)
		}
	}()
}

func (m *Monitor) peersProbe(fun localPeerTopologyFunc) {
	root, peers, err := fun()
	if err != nil {
		log.Errorf("error getting local peers: %v", err)
		return
	}

	for _, peer := range peers {
		go m.probe(root, peer)
	}
}

func (m *Monitor) probe(root, peer net.IP) {
	pinger, err := ping.NewPinger(peer.String())
	if err != nil {
		log.Info("error creating pinger:", err)
		return
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		log.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		log.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = probeCount
	pinger.Timeout = timeout

	err = pinger.Run()
	if err != nil {
		log.Println("failed to ping target host:", err)
	}

}
