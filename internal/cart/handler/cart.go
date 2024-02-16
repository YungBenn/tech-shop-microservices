package handler

import (
	"context"
	"net/http"

	"github.com/YungBenn/tech-shop-microservices/internal/cart/entity"
	cartPb "github.com/YungBenn/tech-shop-microservices/internal/cart/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/cart/repository"
	productPb "github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartServiceServer struct {
	cartPb.UnimplementedCartServiceServer
	client productPb.ProductServiceClient
	log    *logrus.Logger
	rdb    repository.CartRepository
}

func NewCartServiceServer(log *logrus.Logger, rdb repository.CartRepository) cartPb.CartServiceServer {
	return &CartServiceServer{
		UnimplementedCartServiceServer: cartPb.UnimplementedCartServiceServer{},
		log:                            log,
		rdb:                            rdb,
	}
}

func (s *CartServiceServer) AddToCart(ctx context.Context, req *cartPb.AddToCartRequest) (*cartPb.AddToCartResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	readProductReq := &productPb.ReadProductRequest{
		Id: req.ProductId,
	}

	res, err := s.client.ReadProduct(ctx, readProductReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error reading product: %v", err)
	}

	value := entity.Products{
		ProductID: res.Product.Id,
		Quantity:  req.Quantity,
	}

	err = s.rdb.Add(authPayload.UserID, value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error adding to cart: %v", err)
	}

	return &cartPb.AddToCartResponse{
		Status:  http.StatusOK,
		Message: "Successfully added to cart",
	}, nil
}

func (s *CartServiceServer) RemoveFromCart(ctx context.Context, req *cartPb.RemoveFromCartRequest) (*cartPb.RemoveFromCartResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	err = s.rdb.Remove(authPayload.UserID, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error removing from cart: %v", err)
	}

	return &cartPb.RemoveFromCartResponse{
		Status:  http.StatusOK,
		Message: "Successfully removed from cart",
	}, nil
}

func (s *CartServiceServer) CartList(ctx context.Context, req *cartPb.CartListRequest) (*cartPb.CartListResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	cart, err := s.rdb.Get(authPayload.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting cart: %v", err)
	}

	cartList := make([]*cartPb.Product, len(cart.Products))
	for i, product := range cart.Products {
		cartList[i] = utils.ConvertCartProduct(product)
	}

	return &cartPb.CartListResponse{
		Status:  http.StatusOK,
		Message: "Successfully retrieved cart",
		UserId: authPayload.UserID,
		Cart: cartList,
	}, nil
}
