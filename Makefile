
GOPATH=$(shell pwd)/.go
#GOPATH=$(shell pwd)/vendor

echo:
	@echo $(GOPATH)

deps:
	go get -u github.com/eyedeekay/torcat

torcat:
	go build -o ./torcat ./cmd
