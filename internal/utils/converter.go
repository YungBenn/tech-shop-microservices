package utils

import (
	"time"

	userEntity "github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	userPb "github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	productEntity "github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	productPb "github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	cartEntity "github.com/YungBenn/tech-shop-microservices/internal/cart/entity"
	cartPb "github.com/YungBenn/tech-shop-microservices/internal/cart/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user userEntity.User) *userPb.User {
	return &userPb.User{
		Email:       user.Email,
		Password:    user.Password,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth.Format(time.DateOnly),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}

func ConvertProduct(product productEntity.Product) *productPb.Product {
	return &productPb.Product{
		Id:          product.ID.Hex(),
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Image:       product.Image,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ConvertCartProduct(cart cartEntity.Products) *cartPb.Product {
	return &cartPb.Product{
		Id:       cart.ProductID,
		Quantity: cart.Quantity,
	}
}
