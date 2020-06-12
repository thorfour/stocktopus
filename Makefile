.PHONY: all clean docker

GO_VERSION="1.14"

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

all: bin 
setup: 
	mkdir -p ./bin/
bin: setup
	$(go) build -o ./bin/stocktopus ./cmd/stocktopus
docker: bin
	docker build . -t quay.io/thorfour/stocktopus
clean:
	rm -r ./bin
test: 
	$(go) test -v ./...
push:
	echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USERNAME) --password-stdin quay.io
	docker push quay.io/thorfour/stocktopus

# Make targets for circle ci builds
circle-ci-bin: setup
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/stocktopus ./cmd/stocktopus
circle-ci-test:
	go test -v ./...
circle-ci-docker: circle-ci-bin
	docker build . -t quay.io/thorfour/stocktopus
