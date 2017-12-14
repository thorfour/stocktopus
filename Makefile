.PHONY: all rtm aws gcp auth clean

RTM_NAME='rtm'
AUTH_NAME='oauth'
SERVERLESS_NAME='serverless'
AWS_TAG='AWS'
GCP_TAG='GCP'

setup: 
	mkdir -p ./bin/
all:  rtm aws gcp 
rtm: setup
	mkdir  -p ./bin/rtm
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/rtm/$(RTM_NAME) ./cmd/rtm/
aws: setup
	mkdir -p ./bin/aws
	CGO_ENABLED=0 GOOS=linux go build -tags $(AWS_TAG) -o ./bin/aws/$(SERVERLESS_NAME) ./cmd/serverless/
	zip  -j ./bin/aws/stocktopus.zip ./build/serverless/aws/* ./bin/aws/*
	mkdir -p ./bin/aws/auth
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/aws/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/aws/auth/authtopus.zip ./build/authtopus/aws/* ./bin/auth/*
gcp: setup
	mkdir -p ./bin/gcp
	CGO_ENABLED=0 GOOS=linux go build -tags $(GCP_TAG) -o ./bin/gcp/$(SERVERLESS_NAME) ./cmd/serverless/
	zip  -j ./bin/gcp/stocktopus.zip ./build/serverless/gcp/* ./bin/gcp/*
	mkdir -p ./bin/gcp/auth
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/gcp/auth/$(AUTH_NAME) ./cmd/oauth/
	zip -j ./bin/gcp/auth/authtopus.zip ./build/authtopus/gcp/* ./bin/auth/*
	mkdir -p ./bin/gcp/kick
	zip -j ./bin/gcp/kick/kick.zip ./build/serverless/gcp/kick/*
clean:
	rm -r ./bin
