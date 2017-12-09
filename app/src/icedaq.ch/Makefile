# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: cbuild sbuild fmt clean copy lint vet

# DIND_CONTAINER=contained-dind

all: clean fmt lint cbuild sbuild copy

# build the client for linux
cbuild:
	@echo "+ $@"
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ticketclient client.go monitoring.go

# build the server for linux
sbuild:
	@echo "+ $@"
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ticketserver main.go server.go databases.go monitoring.go

rserver:
	@echo "+ $@"
	go run main.go server.go databases.go monitoring.go

rclient:
	@echo "+ $@"
	go run client.go monitoring.go

# copy the static binaries to the corresponding docker folders.
copy:
	@echo "+ $@"
	cp ticketclient ../../../docker/client
	cp ticketserver ../../../docker/server

fmt:
	@echo "+ $@"
	@gofmt -s -l *.go

 lint:
	@echo "+ $@"
	@golint *.go

# test: fmt lint vet
# 	@echo "+ $@"
# 	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet *.go

clean:
	@echo "+ $@"
	@rm -rf ticketserver
	@rm -rf ticketclient

# build the containers and push them to the gke registry.
containers:
	@echo "+ $@"
	sudo docker build -t ticketserver ../../../docker/server/
	sudo docker build -t ticketclient ../../../docker/client/
	sudo docker tag ticketclient gcr.io/default-1296/ticketclient
	sudo docker tag ticketserver gcr.io/default-1296/ticketserver
	gcloud docker -- push gcr.io/default-1296/ticketclient
	gcloud docker -- push gcr.io/default-1296/ticketserver

# run: dbuild
# 	docker run --rm -it \
# 		-v $(CURDIR)/.certs:/etc/docker/ssl:ro \
# 		--net container:$(DIND_CONTAINER) \
# 		$(DOCKER_IMAGE) -d \
# 		--dcacert=/etc/docker/ssl/cacert.pem \
# 		--dcert=/etc/docker/ssl/client.cert \
# 		--dkey=/etc/docker/ssl/client.key

# devbuild:
# 	docker build --rm --force-rm -f Dockerfile.dev -t $(DOCKER_IMAGE):dev .