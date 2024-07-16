package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/token"
	"github.com/YungBenn/tech-shop-microservices/internal/utils"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	log  *logrus.Logger
	rdb  token.TokenRepository
	repo repository.AuthRepository
}

func NewAuthServiceServer(
	log *logrus.Logger,
	rdb token.TokenRepository,
	repo repository.AuthRepository,
) pb.AuthServiceServer {
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
		return nil, status.Errorf(codes.Internal, "Error hashing password: %v", err)
	}

	dateOFBirth, err := time.Parse(time.DateOnly, req.DateOfBirth)
	if err != nil {
		s.log.Error("Error parsing date of birth: ", err)
		return nil, status.Errorf(codes.Internal, "Error parsing date of birth: %v", err)
	}

	arg := entity.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: dateOFBirth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	user, err := s.repo.SaveUser(arg)
	if err != nil {
		s.log.Error("Error saving user:	", err)
		return nil, status.Errorf(codes.Internal, "Error saving user: %v", err)
	}

	s.log.Info("User saved:	", user.ID)
	return &pb.RegisterResponse{
		Status:  http.StatusCreated,
		Message: "Register successful",
		User:    utils.ConvertUser(arg),
	}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	record, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		s.log.Error("Error finding user: ", err)
		return nil, status.Errorf(codes.NotFound, "Error finding user: %v", err)
	}

	s.log.Info(record.Password)
	err = utils.CheckPassword(req.Password, record.Password)
	if err != nil {
		s.log.Error("Error checking password: ", err)
		return nil, status.Errorf(codes.Internal, "Error checking password: %v", err)
	}

	tokenString, claim, err := token.GenerateJWT(record)
	if err != nil {
		s.log.Error("Error generating JWT: ", err)
		return nil, status.Errorf(codes.Internal, "Error generating JWT: %v", err)
	}

	token := entity.Token{
		Token:  tokenString,
		Expiry: claim.ExpiresAt.Unix(),
	}

	err = s.rdb.SetToken(record.ID, token)
	if err != nil {
		s.log.Error("Error setting token: ", err)
		return nil, status.Errorf(codes.Internal, "Error setting token: %v", err)
	}

	s.log.Info("User logged in: ", record.ID)
	return &pb.LoginResponse{
		Status:  http.StatusOK,
		Message: "User logged in",
		Token:   tokenString,
	}, nil
}
