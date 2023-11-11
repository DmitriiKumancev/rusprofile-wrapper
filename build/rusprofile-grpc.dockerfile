# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app

# restore packages
COPY ["go.mod", "go.sum", "./"]
RUN go get -d -v ./...

# build
COPY . .
RUN go install -v ./cmd/rusprofile-grpc

# final stage
FROM alpine:latest as run
COPY --from=builder /go/bin/rusprofile-grpc /app
COPY --from=builder /go/src/app/api/openapiv2 /api/openapiv2
COPY --from=builder /go/src/app/static /static
ENTRYPOINT ["/app"]
