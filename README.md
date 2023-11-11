# gRPC Wrapper for Rusprofile.ru

This wrapper provides seamless access to [rusprofile.ru](https://www.rusprofile.ru/) data through gRPC, with an additional HTTP API accessible via the HTTP-to-gRPC gateway.

## Table of Contents
- [Run](#run)
- [Usage](#usage)
  - [Browser](#browser)
  - [gRPC](#grpc)
  - [HTTP](#http)
- [Project Structure](#project-structure)
- [Tasks](#tasks)

## Run

To launch the Docker container, execute the following command:

```shell
docker run -it -p 8080:8080 -p 8888:8888 --rm dmkumantsev/rusprofile-grpc:v1.0.1
```

## Usage

### Browser

Open [Swagger UI](http://localhost:8080/swagger-ui/) in your browser.

### gRPC

Utilize `grpcurl` (equivalent to `curl` for gRPC) for testing the gRPC API:

```shell
# Returns 'not found'
grpcurl -plaintext -d '{"INN": "123"}' localhost:8888 rusprofile.v1.CompanyFinder/ByINN

# Returns '200'
grpcurl -plaintext -d '{"INN": "5003028028"}' localhost:8888 rusprofile.v1.CompanyFinder/ByINN
```

### HTTP

Use `curl` to test the HTTP API:

```shell
# Returns 'not found'
curl localhost:8080/v1/company/123

# Returns '200'
curl localhost:8080/v1/company/5003028028
```

## Tasks

### Technologies
- The service is written in Go using Go Modules.
- Provides API through gRPC.
- Offers API through HTTP using grpc-gateway.
- Includes Swagger UI with documentation generated from the .proto file using protoc-gen-swagger. Documentation is accessible at /swaggerui.
- Packaged in a Docker container.

## Project Structure

- **api:** Proto files and buf configuration.
  - **openapiv2:** Generated swagger.json file.
- **build:** Dockerfile.
- **cmd:** Main package.
- **pkg:** Generated Go files and implementation of gRPC server for retrieving data from rusprofile.ru.
- **static/web:** Static files for Swagger UI page.
- **tools:** Imports for tools used in code generation.

Feel free to explore and contribute to this project!

#
**Author:** [Dmitrii Kumancev](https://github.com/DmitriiKumancev)


