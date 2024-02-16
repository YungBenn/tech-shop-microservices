package handler

import (
	"context"
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/search/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/search/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/search/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
	log *logrus.Logger
	es  repository.EsProduct
}

func NewSearchServiceServer(
	log *logrus.Logger,
	es repository.EsProduct,
) pb.SearchServiceServer {
	return &SearchServiceServer{
		UnimplementedSearchServiceServer: pb.UnimplementedSearchServiceServer{},
		log:                              log,
		es:                               es,
	}
}

func (s *SearchServiceServer) CreateIndex(ctx context.Context, req *pb.CreateIndexRequest) (*pb.CreateIndexResponse, error) {
	err := s.es.CreateIndex(ctx, req.Index)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating index: %v", err)
	}

	return &pb.CreateIndexResponse{
		Message: "Index created successfully",
	}, nil
}

func (s *SearchServiceServer) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (*pb.SearchProductResponse, error) {
	products, err := s.es.SearchProduct(ctx, req.Query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error searching product: %v", err)
	}

	pbProducts := make([]*pb.ProductData, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.ProductData{
			Id:          product.ID,
			Title:       product.Title,
			Price:       product.Price,
			Tag:         product.Tag,
			Discount:    product.Discount,
			Image:       product.Image,
			Description: product.Description,
			CreatedBy:   product.CreatedBy,
			CreatedAt:   timestamppb.New(time.Unix(product.CreatedAt, 0)),
			UpdatedAt:   timestamppb.New(time.Unix(product.UpdatedAt, 0)),
		}
	}

	return &pb.SearchProductResponse{
		Product: pbProducts,
	}, nil
}

func (s *SearchServiceServer) UpdateIndexProduct(ctx context.Context, req *pb.UpdateIndexProductRequest) (*pb.UpdateIndexProductResponse, error) {
	product := entity.ProductData{
		ID:          req.Product.Id,
		Title:       req.Product.Title,
		Price:       req.Product.Price,
		Tag:         req.Product.Tag,
		Discount:    req.Product.Discount,
		Image:       req.Product.Image,
		Description: req.Product.Description,
		CreatedBy:   req.Product.CreatedBy,
		UpdatedAt:   time.Now().Unix(),
	}

	err := s.es.UpdateProduct(ctx, product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error updating product: %v", err)
	}

	return &pb.UpdateIndexProductResponse{
		Message: "Product updated successfully",
	}, nil
}
