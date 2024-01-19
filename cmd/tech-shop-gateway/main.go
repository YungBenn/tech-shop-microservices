package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/YungBenn/tech-shop-microservices/config"
	authSvc "github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	productSvc "github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ctx = context.Background()

func initAuthServiceClient(c config.EnvVars, mux *runtime.ServeMux, log *logrus.Logger) error {
	authServerUrl := fmt.Sprintf("%s:%s", c.AuthServiceHost, c.AuthServicePort)
	conn, err := grpc.DialContext(ctx, authServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	err = authSvc.RegisterAuthServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Error("Error registering auth service handler: ", err)
	}

	return nil
}

func initProductServiceClient(c config.EnvVars, mux *runtime.ServeMux, log *logrus.Logger) error {
	productServerUrl := fmt.Sprintf("%s:%s", c.ProductServiceHost, c.ProductServicePort)
	conn, err := grpc.DialContext(ctx, productServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	err = productSvc.RegisterProductServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Error("Error registering product service handler: ", err)
	}

	return nil
}

// func initCartServiceClient(c config.EnvVars) (*grpc.ClientConn, error) {
// 	cartServerUrl := fmt.Sprintf("%s:%s", c.CartServiceHost, c.CartServicePort)
// 	conn, err := grpc.DialContext(ctx, cartServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

func main() {
	log := logrus.New()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	mux := runtime.NewServeMux()

	err = initAuthServiceClient(conf, mux, log)
	if err != nil {
		log.Error("Error initializing auth service client: ", err)
	}

	err = initProductServiceClient(conf, mux, log)
	if err != nil {
		log.Error("Error initializing product service client: ", err)
	}

	clientUrl := fmt.Sprintf("%s:%s", conf.ClientHost, conf.ClientPort)
	gwServer := &http.Server{
		Addr:    clientUrl,
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
