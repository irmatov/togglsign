package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/irmatov/togglsign/app"
	"github.com/irmatov/togglsign/infra/adapters/storage"
	"github.com/irmatov/togglsign/infra/ports/rpc"
	"google.golang.org/grpc"
)

const (
	envNamePrefix  = "SIGNER_"
	envDSN         = envNamePrefix + "DSN"
	envJWTKey      = envNamePrefix + "JWT_KEY"
	envSignKey     = envNamePrefix + "SIGN_KEY"
	envGRPCAddress = envNamePrefix + "GRPC_ADDRESS"
)

func run() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()

	db, err := storage.New(ctx, os.Getenv(envDSN))
	if err != nil {
		return err
	}
	application := app.New(db, []byte(os.Getenv(envJWTKey)), os.Getenv(envSignKey))
	srv := rpc.New(application)
	grpcSrv := grpc.NewServer()
	rpc.RegisterSignerServer(grpcSrv, srv)
	lis, err := net.Listen("tcp", os.Getenv(envGRPCAddress))
	if err != nil {
		return err
	}
	go func() {
		<-sig
		grpcSrv.GracefulStop()
	}()
	grpcSrv.Serve(lis)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
