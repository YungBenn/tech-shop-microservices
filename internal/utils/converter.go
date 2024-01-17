package utils

import (
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user entity.User) *pb.User {
	return &pb.User{
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