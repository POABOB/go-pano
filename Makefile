
 
.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) -t poabob/protoc-builder --proto_path=. --micro_out=. --go_out=:. ./protos/predict/predict.proto	
# docker run --rm -v $(pwd):$(pwd) -w $(pwd) -t poabob/protoc-builder --proto_path=. --micro_out=. --go_out=:. ./protos/predict/predict.proto
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-pano.output *.go

mac-build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o go-pano.output *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t predict:latest
