# cleanup dependencies and download missing ones
.PHONY: deps
deps:
	go mod tidy
	go mod download

# run dependency cleanup, followed by updating the patch version
.PHONY: deps-update
deps-update: deps
	go get -u=patch
	
# run tests
.PHONY: tests
tests:
	go test -race -cover -count 1 ./...

# run standard go tooling for better rcode hygiene
.PHONY: tidy
tidy: imports fmt
	go vet ./...
	golint ./...

# automatically add missing imports
.PHONY: imports
imports:
	find . -type f -name '*.go' -exec goimports -w {} \;

# format code and simplify if possible
.PHONY: fmt
fmt:
	find . -type f -name '*.go' -exec gofmt -s -w {} \;

verifiers: staticcheck

staticcheck:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/staticcheck ./...

# runs the 
.PHONY: run-swarm
run-swarm:
	docker run --network host --name temporal_swarm -d -it -v ${PWD}/swarmtest/datadir:/data \
					-v ${PWD}/swarmtest/passwordfile:/password \
					ethersphere/swarm:0.5.7 \
								--datadir /data \
								--password /password \
								--debug

.PHONY: stop-swarm
stop-swarm:
	docker stop temporal_swarm && docker rm temporal_swarm