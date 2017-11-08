.PHONY: all rtm aws gcp auth clean

RTM_NAME='rtm'
AUTH_NAME='oauth'
SERVERLESS_NAME='serverless'

setup: 
	mkdir -p ./bin/
all:  rtm aws gcp auth
rtm: setup
	mkdir ./bin/rtm
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/rtm/$(RTM_NAME) ./cmd/rtm/*.go
aws: setup
	mkdir ./bin/aws
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/aws/$(SERVERLESS_NAME) ./cmd/serverless/*.go
	zip  -j ./bin/aws/stocktopus.zip ./build/aws/* ./bin/aws/*
gcp: setup
	mkdir ./bin/gcp
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/gcp/$(SERVERLESS_NAME) ./cmd/serverless/*.go
	zip  -j ./bin/gcp/stocktopus.zip ./build/gcp/* ./bin/gcp/*
auth: setup
	mkdir ./bin/auth
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/auth/$(AUTH_NAME) ./cmd/oauth/*.go
	zip -j ./bin/auth/authtopus.zip ./build/authtopus/* ./bin/auth/*
clean:
	rm -r ./bin
