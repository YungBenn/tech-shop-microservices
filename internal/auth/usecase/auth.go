package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/jwt"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/redis"
	"github.com/YungBenn/tech-shop-microservices/internal/utils"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	log  *logrus.Logger
	rdb  redis.TokenRepository
	repo repository.AuthRepository
}

func NewAuthServiceServer(log *logrus.Logger, rdb redis.TokenRepository, repo repository.AuthRepository) pb.AuthServiceServer {
	return &AuthServiceServer{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		log:                            log,
		rdb:                            rdb,
		repo:                           repo,
	}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Error hashing password: ", err)
		return nil, err
	}

	arg := entity.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: time.Time{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	user, err := s.repo.SaveUser(arg)
	if err != nil {
		s.log.Error("Error saving user:	", err)
		return nil, err
	}

	s.log.Info("User saved:	", user.ID)
	return &pb.RegisterResponse{
		Status:  http.StatusCreated,
		Message: "User saved",
		User:    &pb.RegisterRequest{},
	}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	record, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		s.log.Error("Error finding user: ", err)
		return nil, err
	}

	err = utils.CheckPassword(req.Password, record.Password)
	if err != nil {
		s.log.Error("Error checking password: ", err)
		return nil, err
	}

	tokenString, err := jwt.GenerateJWT(record.ID, record.Email)
	if err != nil {
		s.log.Error("Error generating JWT: ", err)
		return nil, err
	}

	tokenExpiry := time.Now().Add(1 * time.Hour).Unix()

	token := redis.Token{
		Token:  tokenString,
		Expiry: tokenExpiry,
	}

	err = s.rdb.SetToken(record.ID, token)
	if err != nil {
		s.log.Error("Error setting token: ", err)
		return nil, err
	}

	s.log.Info("User logged in: ", record.ID)
	return &pb.LoginResponse{
		Status:  http.StatusOK,
		Message: "User logged in",
		Token:   tokenString,
	}, nil
}
