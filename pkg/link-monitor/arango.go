package link_monitor

import (
	"context"
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
		log.Println("graph exists already")
		return
	}

	g, err := a.db.CreateGraph(nil, name, nil)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}

	a.graph = g
	vertices, _ := a.graph.CreateVertexCollection(nil, "nodes")
	constraints := driver.VertexConstraints{
		From: []string{"nodes"},
		To:   []string{"nodes"},
	}
	a.vertices = vertices
	edges, _ := a.graph.CreateEdgeCollection(nil, "edges", constraints)
	a.edges = edges
}

func (a *ArangoClient) AddEdge(src, dst string, stats *ping.Statistics) {
	// Create edge
	edge := map[string]interface{}{
		"_from":  src,
		"_to":    dst,
		"rtt":    stats.AvgRtt,
		"jitter": stats.StdDevRtt,
		"loss":   stats.PacketLoss,
	}

	a.vertices.CreateDocument(nil, src)
	a.vertices.CreateDocument(nil, dst)

	a.edges.CreateDocument(nil, edge)
}
