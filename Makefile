proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/grpc/registry.proto

nats:
    docker run -d -p 4222:4222 -i nats

arango:
    docker run -d -p 8529:8529 -e ARANGO_ROOT_PASSWORD=cil_overly_v1 arangodb/arangodb:3.9.0