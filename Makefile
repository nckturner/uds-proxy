PROXY_SRC = $(wildcard signer/*.go)
CLIENT_SRC = $(wildcard client/*.go)
SERVER_SRC = $(wildcard server/*.go)

all: bin/proxy bin/client bin/server

bin/proxy: $(PROXY_SRC)
	go build -o ./bin/proxy ./proxy/...

bin/client: $(CLIENT_SRC)
	go build -o ./bin/client ./client/...

bin/server: $(SERVER_SRC)
	go build -o ./bin/server ./server/...

.PHONY: bench
bench:
	cd uds && go test -bench=./uds/...
