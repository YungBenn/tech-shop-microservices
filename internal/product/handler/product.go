package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/product/message"
	"github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
	log  *logrus.Logger
	repo repository.ProductRepository
	msg  message.KafkaProducerRepository
}

func NewProductServiceServer(
	log *logrus.Logger,
	repo repository.ProductRepository,
	msg message.KafkaProducerRepository,
) pb.ProductServiceServer {
	return &ProductServiceServer{
		UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{},
		log:                               log,
		repo:                              repo,
		msg:                               msg,
	}
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		s.log.Errorf(ErrAuthUser, err)
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	arg := entity.Product{
		Title:       req.Title,
		Price:       req.Price,
		Tag:         req.Tag,
		Discount:    req.Discount,
		Stock:       int(req.Stock),
		Image:       req.Image,
		Description: req.Description,
		CreatedBy:   authPayload.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	product, err := s.repo.InsertProduct(ctx, arg)
	if err != nil {
		s.log.Error("Error saving product: ", err)
		return nil, status.Errorf(codes.Internal, "Error saving product: %v", err)
	}

	err = s.msg.Publish(*product)
	if err != nil {
		s.log.Error("Error publishing product: ", err)
		return nil, status.Errorf(codes.Internal, "Error publishing product: %v", err)
	}

	s.log.Info("Product saved: ", product.ID)
	return &pb.CreateProductResponse{
		Status:  http.StatusOK,
		Message: "Product created successful",
		Product: utils.ConvertProduct(*product),
	}, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	products, paginationData, err := s.repo.FindAllProducts(ctx, int(req.Limit), int(req.Page))
	if err != nil {
		s.log.Error("Error listing products: ", err)
		return nil, status.Errorf(codes.Internal, "Error listing products: %v", err)
	}

	productList := make([]*pb.Product, len(products))
	for i, product := range products {
		productList[i] = utils.ConvertProduct(product)
	}

	s.log.Info("Listing products successful")
	return &pb.ListProductResponse{
		Limit:      req.Limit,
		Page:       req.Page,
		TotalRows:  paginationData.TotalRows,
		TotalPages: paginationData.TotalPages,
		Product:    productList,
	}, nil
}

func (s *ProductServiceServer) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	product, err := s.repo.FindOneProduct(ctx, req.Id)
	if err != nil {
		s.log.Error("Error reading product: ", err)
		return nil, status.Errorf(codes.Internal, "Error reading product: %v", err)
	}

	s.log.Info("Reading product successful")
	return &pb.ReadProductResponse{
		Product: utils.ConvertProduct(*product),
	}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		s.log.Errorf(ErrAuthUser, err)
		return nil, status.Errorf(codes.Internal, "Error authorizing user: %v", err)
	}

	findProduct, err := s.repo.FindOneProduct(ctx, req.Id)
	if err != nil {
		s.log.Errorf(ErrFindingProduct, err)
		return nil, status.Errorf(codes.NotFound, ErrFindingProduct, err)
	}

	if findProduct.CreatedBy != authPayload.UserID {
		s.log.Error("You are not authorized to update this product")
		return nil, status.Errorf(codes.PermissionDenied, "You are not authorized to update this product")
	}

	arg := entity.Product{
		Title:       req.Title,
		Price:       req.Price,
		Tag:         req.Tag,
		Discount:    req.Discount,
		Stock:       int(req.Stock),
		Image:       req.Image,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	product, err := s.repo.UpdateProduct(ctx, req.Id, arg)
	if err != nil {
		s.log.Error("Error updating product: ", err)
		return nil, status.Errorf(codes.Internal, "Error updating product: %v", err)
	}

	s.log.Info("Product updated: ", product.ID)
	return &pb.UpdateProductResponse{
		Product: utils.ConvertProduct(*product),
	}, nil
}

func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		s.log.Errorf(ErrAuthUser, err)
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	findProduct, err := s.repo.FindOneProduct(ctx, req.Id)
	if err != nil {
		s.log.Errorf(ErrFindingProduct, err)
		return nil, status.Errorf(codes.NotFound, ErrFindingProduct, err)
	}

	if findProduct.CreatedBy != authPayload.UserID {
		s.log.Error("You are not authorized to delete this product")
		return nil, status.Errorf(codes.PermissionDenied, "You are not authorized to delete this product")
	}

	err = s.repo.DeleteProduct(ctx, req.Id)
	if err != nil {
		s.log.Error("Error deleting product: ", err)
		return nil, status.Errorf(codes.Internal, "Error deleting product: %v", err)
	}

	s.log.Info("Product deleted: ", req.Id)
	return &pb.DeleteProductResponse{
		Message: "Product deleted successful: " + req.Id,
	}, nil
}

func (s *ProductServiceServer) DecreaseProductStock(ctx context.Context, req *pb.DecreaseProductStockRequest) (*pb.DecreaseProductStockResponse, error) {
	authPayload, err := utils.AuthorizeUser(ctx)
	if err != nil {
		s.log.Errorf(ErrAuthUser, err)
		return nil, status.Errorf(codes.Internal, ErrAuthUser, err)
	}

	findProduct, err := s.repo.FindOneProduct(ctx, req.Id)
	if err != nil {
		s.log.Errorf(ErrFindingProduct, err)
		return nil, status.Errorf(codes.NotFound, ErrFindingProduct, err)
	}

	if findProduct.CreatedBy != authPayload.UserID {
		s.log.Error("You are not authorized to manage this product")
		return nil, status.Errorf(codes.PermissionDenied, "You are not authorized to manage this product")
	}

	arg := entity.Product{
		Stock: findProduct.Stock - int(req.Quantity),
	}

	product, err := s.repo.UpdateProduct(ctx, req.Id, arg)
	if err != nil {
		s.log.Error("Error updating product: ", err)
		return nil, status.Errorf(codes.Internal, "Error updating product: %v", err)
	}

	s.log.Info("Product updated: ", product.ID)
	return &pb.DecreaseProductStockResponse{
		Product: utils.ConvertProduct(*product),
	}, nil
}
