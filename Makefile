.PHONY: all rtm aws gcp auth clean

RTM_NAME='rtm'
AUTH_NAME='oauth'
SERVERLESS_NAME='serverless'

all:  rtm aws gcp docker test 
setup: 
	mkdir -p ./bin/
rtm: setup
	mkdir  -p ./bin/rtm
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/rtm/$(RTM_NAME) ./cmd/rtm/
aws: setup
	mkdir -p ./bin/aws
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/aws/$(SERVERLESS_NAME) ./cmd/awslambda/
	zip  -j ./bin/aws/stocktopus.zip ./build/serverless/aws/* ./bin/aws/*
	mkdir -p ./bin/aws/auth
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/aws/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/aws/auth/authtopus.zip ./build/authtopus/aws/* ./bin/aws/auth/*
gcp: setup
	mkdir -p ./bin/gcp
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/gcp/$(SERVERLESS_NAME) ./cmd/gcpfunction/
	zip  -j ./bin/gcp/stocktopus.zip ./build/serverless/gcp/* ./bin/gcp/*
	mkdir -p ./bin/gcp/auth
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/gcp/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/gcp/auth/authtopus.zip ./build/authtopus/gcp/* ./bin/gcp/auth/*
	mkdir -p ./bin/gcp/kick
	zip -j ./bin/gcp/kick/kick.zip ./build/serverless/gcp/kick/*
docker: setup
	mkdir -p ./bin/docker
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/docker/server ./cmd/server
	cp /etc/ssl/certs/ca-certificates.crt ./bin/docker/
	cp ./build/docker/Dockerfile ./bin/docker/
	docker build ./bin/docker/ -t quay.io/thorfour/stocktopus

clean:
	rm -r ./bin

test: 
	go test -v ./...
