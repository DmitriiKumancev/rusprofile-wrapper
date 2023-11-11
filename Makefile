IMAGE = dmkumantsev/rusprofile-grpc
TAG = v1.0.1

.PHONY: proto-generate build run test build-image push-image run-image build-push-image

proto-generate:
	cd ./api && buf generate

build:
	go build ./...

run:
	go run ./cmd/rusprofile-grpc

test:
	go test ./...

build-image:
	docker build -t $(IMAGE):$(TAG) -f ./build/rusprofile-grpc.dockerfile .

push-image:
	docker push $(IMAGE):$(TAG)

run-image:
	docker run -it -p 8080:8080 -p 8888:8888 --rm $(IMAGE):$(TAG)

build-push-image: build-image push-image


