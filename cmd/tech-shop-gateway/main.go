package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/YungBenn/tech-shop-microservices/config"
	authSvc "github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ctx = context.Background()

func initAuthServiceClient(c config.EnvVars) (*grpc.ClientConn, error) {
	authServerUrl := fmt.Sprintf("%s:%s", c.AuthServiceHost, c.AuthServicePort)
	conn, err := grpc.DialContext(ctx, authServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initProductServiceClient(c config.EnvVars) (*grpc.ClientConn, error) {
	productServerUrl := fmt.Sprintf("%s:%s", c.ProductServiceHost, c.ProductServicePort)
	conn, err := grpc.DialContext(ctx, productServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initCartServiceClient(c config.EnvVars) (*grpc.ClientConn, error) {
	cartServerUrl := fmt.Sprintf("%s:%s", c.CartServiceHost, c.CartServicePort)
	conn, err := grpc.DialContext(ctx, cartServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	log := logrus.New()

	config, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	authConn, err := initAuthServiceClient(config)
	if err != nil {
		log.Error("Error initializing auth service client: ", err)
	}

	mux := runtime.NewServeMux()

	err = authSvc.RegisterAuthServiceHandler(ctx, mux, authConn)
	if err != nil {
		log.Panic("Error registering auth service handler: ", err)
	}

	clientUrl := fmt.Sprintf("%s:%s", config.ClientHost, config.ClientPort)
	gwServer := &http.Server{
		Addr: clientUrl,
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}