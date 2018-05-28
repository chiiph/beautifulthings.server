SHELL:=/bin/bash

.PHONY: install test build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

test:
	go test -v ./...

#go build -ldflags "-X main.version=$(TAG)" -o news .
build:
	go build -o serv ./cmd/server

pack: build
	docker build -t docker.io/chiiph/beautifulthings:$(TAG) .

upload: pack
	docker push docker.io/chiiph/beautifulthings:$(TAG)

deploy:
	envsubst < ./kubernetes/deployment.yaml | kubectl apply -f -

clear-cluster:
	kubectl delete deployment,service beautifulthings

ship: test pack upload deploy
