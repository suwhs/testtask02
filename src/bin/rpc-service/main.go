package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"google.golang.org/grpc"
	"whs.su/rusprofile/src/rest"
	"whs.su/rusprofile/src/rpc"
	"whs.su/rusprofile/src/server"
)

func main() {
	var port int
	var rest_port int = 80

	ctx := context.TODO()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		log.Printf("graceful shutdown....")
		cancel()
		log.Printf("done")
	}()

	flag.IntVar(&port, "port", 8193, "define grpc srvice port (env var RUSPROFILE_GRPC_PORT)")
	flag.IntVar(&rest_port, "rest service port", 8080, "define rest http service port (env var RUSPROFILE_REST_PORT)")

	if _port := os.Getenv("RUSPROFILE_GRPC_PORT"); _port != "" {
		if _i, err := strconv.ParseInt(_port, 10, 32); err == nil {
			port = int(_i)
		}
	}
	if _rest_port := os.Getenv("RUSPROFILE_REST_PORT"); _rest_port != "" {
		if _i, err := strconv.ParseInt(_rest_port, 10, 32); err == nil {
			rest_port = int(_i)
		}
	}
	flag.Parse()
	grpc_addr := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterRusprofileServer(grpcServer, server.NewServer())

	rest.RunRestServer(ctx, grpc_addr, rest_port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc server error: %s", err.Error())
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
}
