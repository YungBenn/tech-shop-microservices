package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/utils"
	"github.com/sirupsen/logrus"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
	log  *logrus.Logger
	repo repository.ProductRepository
}

func NewProductServiceServer(log *logrus.Logger, repo repository.ProductRepository) pb.ProductServiceServer {
	return &ProductServiceServer{
		UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{},
		log:                               log,
		repo:                              repo,
	}
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	arg := entity.Product{
		Title:       req.Title,
		Price:       req.Price,
		Tag:         req.Tag,
		Discount:    req.Discount,
		Image:       req.Image,
		Description: req.Description,
	}

	product, err := s.repo.InsertProduct(ctx, arg)
	if err != nil {
		s.log.Error("Error saving product: ", err)
		return nil, err
	}

	s.log.Info("Product saved: ", product.ID)
	return &pb.CreateProductResponse{
		Status:  http.StatusOK,
		Message: "Product created successful",
		Product: utils.ConvertProduct(arg),
	}, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
    products, err := s.repo.FindAllProducts(ctx, int(req.Limit), int(req.Page))
    if err != nil {
        s.log.Error("Error listing products: ", err)
        return nil, err
    }

    productList := make([]*pb.Product, len(products))
    for i, product := range products {
        productList[i] = utils.ConvertProduct(product)
    }

    s.log.Info("Listing products successful")
    return &pb.ListProductResponse{
        Limit:   req.Limit,
        Page:    req.Page,
        Product: productList,
    }, nil
}

func (s *ProductServiceServer) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	product, err := s.repo.FindOneProduct(ctx, req.Id)
	if err != nil {
		s.log.Error("Error reading product: ", err)
		return nil, err
	}

	s.log.Info("Reading product successful")
	return &pb.ReadProductResponse{
		Product: utils.ConvertProduct(*product),
	}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	arg := entity.Product{
		Title:       req.Title,
		Price:       req.Price,
		Tag:         req.Tag,
		Discount:    req.Discount,
		Image:       req.Image,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	product, err := s.repo.UpdateProduct(ctx, req.Id, arg)
	if err != nil {
		s.log.Error("Error updating product: ", err)
		return nil, err
	}

	s.log.Info("Product updated: ", product.ID)
	return &pb.UpdateProductResponse{
		Product: utils.ConvertProduct(arg),
	}, nil
}

func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {	
	err := s.repo.DeleteProduct(ctx, req.Id)
	if err != nil {
		s.log.Error("Error deleting product: ", err)
		return nil, err
	}

	s.log.Info("Product deleted: ", req.Id)
	return &pb.DeleteProductResponse{
		Id: req.Id,
	}, nil
}
