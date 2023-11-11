IMAGE = iskorotkov/rusprofile
TAG = v0.1.0

.PHONY: build run test build-image push-image run-image protoc-update

protoc-update:
	cd pkg && buf generate

build:
	go build ./...

run:
	go run ./cmd/rusprofilegrpc

test:
	go test ./...

build-image:
	docker build -t $(IMAGE):$(TAG) -f ./build/rusprofilegrpc.dockerfile .

push-image:
	docker push $(IMAGE):$(TAG)

run-image:
	docker run -it -p 8888:8888 --rm $(IMAGE):$(TAG)
