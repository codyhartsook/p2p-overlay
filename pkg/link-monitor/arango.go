package link_monitor

import (
	"context"
	"fmt"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
)

type ArangoConfig struct {
	URL      string `desc:"Arangodb server URL (http://127.0.0.1:8529)"`
	User     string `desc:"Arangodb server username"`
	Password string `desc:"Arangodb server user password"`
	Database string `desc:"Arangodb database name"`
}

type ArangoClient struct {
	db       driver.Database
	graph    driver.Graph
	edges    driver.Collection
	vertices driver.Collection
}

type VertexNode struct {
	Key string `json:"_key"` // mandatory field (handle) - short name
}

type EdgeLink struct {
	Key  string `json:"_key"`  // mandatory field (handle)
	From string `json:"_from"` // mandatory field
	To   string `json:"_to"`   // mandatory field

	// other fields … e.g.

	Rtt    int64   `json:"rtt"`
	Jitter int64   `json:"jitter"`
	Loss   float64 `json:"loss"`
}

func NewArango(cfg ArangoConfig) *ArangoClient {
	// Connect to DB
	if cfg.URL == "" || cfg.User == "" || cfg.Password == "" || cfg.Database == "" {
		log.Fatal("ArangoDB Config has an empty field")
	}
	if !strings.Contains(cfg.URL, "http") {
		cfg.URL = "http://" + cfg.URL
	}
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{cfg.URL},
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}

	// Authenticate with DB
	conn, err = conn.SetAuthentication(driver.BasicAuthentication(cfg.User, cfg.Password))
	if err != nil {
		log.Fatalf("Failed to authenticate with arango: %v", err)
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	db, err := ensureDatabase(c, cfg)
	if err != nil {
		log.Fatalf("Failed to create DB")
	}

	return &ArangoClient{db: db}
}

func ensureDatabase(c driver.Client, cfg ArangoConfig) (driver.Database, error) {
	var db driver.Database

	exists, err := c.DatabaseExists(context.Background(), cfg.Database)
	if err != nil {
		return db, err
	}

	if !exists {
		// Create database
		db, err = c.CreateDatabase(context.Background(), cfg.Database, nil)
		if err != nil {
			return db, err
		}
	} else {
		db, err = c.Database(context.Background(), cfg.Database)
		if err != nil {
			return db, err
		}
	}
	return db, nil
}

func (a *ArangoClient) CreateGraph(name string) {
	g_exists, err := a.db.GraphExists(nil, name)
	if g_exists {
		return
	}

	g, err := a.db.CreateGraph(nil, name, nil)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}

	a.graph = g
	vertices, err := a.graph.CreateVertexCollection(nil, "nodes")
	if err != nil {
		log.Fatalf("Failed to create vertex collection: %v", err)
	}

	constraints := driver.VertexConstraints{
		From: []string{"nodes"},
		To:   []string{"nodes"},
	}
	a.vertices = vertices
	edges, err := a.graph.CreateEdgeCollection(nil, "edges", constraints)
	if err != nil {
		log.Fatalf("Failed to create edge collection: %v", err)
	}
	a.edges = edges
}

func (a *ArangoClient) LoadGraph(gName string) {
	g, err := a.db.Graph(nil, gName)
	if err != nil {
		log.Fatalf("Failed to load graph: %v", err)
	}
	a.graph = g

	vertices, err := a.graph.VertexCollection(nil, "nodes")
	if err != nil {
		log.Fatalf("Failed to load vertex collection: %v", err)
	}
	a.vertices = vertices

	edges, _, err := a.graph.EdgeCollection(nil, "edges")
	if err != nil {
		log.Fatalf("Failed to load edge collection: %v", err)
	}
	a.edges = edges
}

func (a *ArangoClient) AddEdge(src, dst string, stats *ping.Statistics) {
	srcV := VertexNode{Key: src}
	if exists, _ := a.vertices.DocumentExists(context.TODO(), src); !exists {
		_, err := a.vertices.CreateDocument(context.TODO(), srcV)
		if err != nil {
			log.Fatalf("Failed to create document: %v", err)
		}
	}

	dstV := VertexNode{Key: dst}
	if exists, _ := a.vertices.DocumentExists(context.TODO(), dst); !exists {
		_, err := a.vertices.CreateDocument(context.TODO(), dstV)
		if err != nil {
			log.Fatalf("Failed to create document: %v", err)
		}
	}

	srcNode := fmt.Sprintf("nodes/%s", src)
	dstNode := fmt.Sprintf("nodes/%s", dst)
	edge := EdgeLink{
		Key:    src + "-" + dst,
		From:   srcNode,
		To:     dstNode,
		Rtt:    stats.AvgRtt.Milliseconds(),
		Jitter: stats.StdDevRtt.Milliseconds(),
		Loss:   float64(stats.PacketLoss),
	}

	if exists, _ := a.edges.DocumentExists(context.TODO(), edge.Key); !exists {
		_, err := a.edges.CreateDocument(context.TODO(), edge)
		if err != nil {
			log.Fatalf("Failed to create document: %v", err)
		}
	} else {
		_, err := a.edges.UpdateDocument(context.TODO(), edge.Key, edge)
		if err != nil {
			log.Fatalf("Failed to update document: %v", err)
		}
	}
}
