.PHONY: all rtm aws gcp auth clean docker

GO_VERSION="1.11"
RTM_NAME='rtm'
AUTH_NAME='oauth'
SERVERLESS_NAME='serverless'

go = @docker run \
		--rm \
		-v ${PWD}:/go/src/github.com/thorfour/stocktopus \
		-w /go/src/github.com/thorfour/stocktopus \
		-u $$(id -u) \
		-e XDG_CACHE_HOME=/tmp/.cache \
		-e CGO_ENABLED=0 \
	    -e GOOS=linux \
		-e GOPATH=/go \
		golang:$(GO_VERSION) \
		go

all:  rtm aws gcp docker
setup: 
	mkdir -p ./bin/
rtm: setup
	mkdir  -p ./bin/rtm
aws: setup
	mkdir -p ./bin/aws
	$(go) build -o ./bin/aws/$(SERVERLESS_NAME) ./cmd/awslambda/
	zip  -j ./bin/aws/stocktopus.zip ./build/serverless/aws/* ./bin/aws/*
	mkdir -p ./bin/aws/auth
	$(go) build -o ./bin/aws/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/aws/auth/authtopus.zip ./build/authtopus/aws/* ./bin/aws/auth/*
gcp: setup
	mkdir -p ./bin/gcp
	$(go) build -o ./bin/gcp/$(SERVERLESS_NAME) ./cmd/gcpfunction/
	zip  -j ./bin/gcp/stocktopus.zip ./build/serverless/gcp/* ./bin/gcp/*
	mkdir -p ./bin/gcp/auth
	$(go) build -o ./bin/gcp/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/gcp/auth/authtopus.zip ./build/authtopus/gcp/* ./bin/gcp/auth/*
	mkdir -p ./bin/gcp/kick
	zip -j ./bin/gcp/kick/kick.zip ./build/serverless/gcp/kick/*
docker: setup
	mkdir -p ./bin/docker
	$(go) build -o ./bin/docker/server ./cmd/server
	cp /etc/ssl/certs/ca-certificates.crt ./bin/docker/
	cp ./build/docker/Dockerfile ./bin/docker/
	docker build ./bin/docker/ -t quay.io/thorfour/stocktopus
docker-alpha:
	mkdir -p ./bin/docker
	$(go) build -tags ALPHA -o ./bin/docker/server ./cmd/server
	cp /etc/ssl/certs/ca-certificates.crt ./bin/docker/
	cp ./build/docker/Dockerfile ./bin/docker/
	docker build ./bin/docker/ -t quay.io/thorfour/stocktopus

clean:
	rm -r ./bin

test: 
	$(go) test -v ./...
