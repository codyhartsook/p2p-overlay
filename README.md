# p2p-overlay

Zero configuration wireguard overlay system. 

## network paths
Overlay configuration is achieved using a one time registry and subsequent  
publisher-subscriber event updates. New nodes perform a grpc handshake with  
a persistent broker node to exchange wireguard public keys and allowed-ips.  
The new peer configuration is then broadcasted over secure tunnels to all  
other peers, propogating the nodes's public-key and allowed-ips. After the  
initial handshake and broadcast, wireguard tunnels are established between  
the new node and all peers in the system. 

## connectivity requirements
Peer nodes must be able to reach the broker's public ip. After the initial  
handshake, internal ips and overlay subnet adresses are used between the  
peers. 

## broker usage  
```
cd cmd/broker/
go build  
sudo ./broker -zone=zone_string
```

## peer usage
```
cd cmd/peer/
go build
sudo ./peer -broker=broker_public_ip -zone=zone_string
```

## overlay topology and tunnel monitoring
The overlay topology is dynamically published to a graph database and can  
be viewed via the arangodb UI console.  
```
http://broker_ip:8529
user: root
pass: cil_overly_v1
```
Graph edges attributes include tunnel round-trip latency, jitter, and  
packet loss. 

## infrastructure
The system utilizes an Arango graph database and a NATS pubsub server.  
These must be deployed as containers on the broker host.  
```
make arango
make nats
```
