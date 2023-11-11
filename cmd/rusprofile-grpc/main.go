package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/DmitriiKumancev/rusprofile-grpc/pkg"
)

var (
	grpcPort        = flag.Int("grpc-port", 8888, "port for gRPC server to listen on")
	httpPort        = flag.Int("http-port", 8080, "port for HTTP server to listen on")
	swaggerSpecPath = flag.String("swagger-spec", "./api/openapiv2/rusprofile-grpc.swagger.json", "path to Swagger spec file")
	staticFilesPath = flag.String("static-files", "./static/web", "path to static files")
)

func main() {
	flag.Parse()

	go startHTTPServer()
	startGRPCServer()
}

func startHTTPServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *swaggerSpecPath)
	})
	mux.HandleFunc("/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(*staticFilesPath, "index.html"))
	})
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/static")
		p = path.Join(*staticFilesPath, p)
		http.ServeFile(w, r, p)
	})

	gw, err := registerGatewayEndpoints()
	if err != nil {
		log.Fatalf("couldn't register gateway endpoints: %v", err)
	}

	mux.Handle("/", gw)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), mux); err != nil {
		log.Fatalf("HTTP server stopped with error: %v", err)
	}
}

func registerGatewayEndpoints() (http.Handler, error) {
	h := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pkg.RegisterCompanyFinderHandlerFromEndpoint(context.Background(), h, fmt.Sprintf(":%d", *grpcPort), opts); err != nil {
		return nil, fmt.Errorf("couldn't register HTTP handler: %w", err)
	}

	return h, nil
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("couldn't start gRPC server: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pkg.RegisterCompanyFinderServer(grpcServer, pkg.CompanyFinder{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server stopped with error: %v", err)
	}
}
