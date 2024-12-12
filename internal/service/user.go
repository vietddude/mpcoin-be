package service

import (
	"context"
	"database/sql"
	"fmt"
	"mpc/internal/db/redis"
	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/cache"
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"mpc/pkg/utils"

	stderrors "errors"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo   *repository.UserRepository
	walletRepo *repository.WalletRepository
	cache      *cache.Cache
}

func NewUserService(userRepo *repository.UserRepository, walletRepo *repository.WalletRepository, redisClient *redis.Client) *UserService {
	return &UserService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		cache:      cache.NewCache(redisClient, "user"),
	}
}

// GetUser get user
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (model.GetUserResponse, error) {
	return cache.FetchOrStore(ctx, s.cache, "user:"+userID.String(), func() (model.GetUserResponse, error) {
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			logger.Error("Service:GetUser", err)
			return model.GetUserResponse{}, fmt.Errorf("getting user: %w", err)
		}
		if user.ID == uuid.Nil {
			return model.GetUserResponse{}, errors.ErrUserNotFound
		}

		wallets, err := s.walletRepo.GetWalletsByUserID(ctx, userID)
		if err != nil {
			logger.Error("Service:GetUser", err)
			return model.GetUserResponse{}, fmt.Errorf("getting wallets: %w", err)
		}
		if len(wallets) == 0 {
			return model.GetUserResponse{}, errors.ErrWalletNotFound
		}

		return model.GetUserResponse{
			User:   utils.ToUserResponse(user),
			Wallet: utils.ToWalletResponse(wallets[0]),
		}, nil
	})
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.ErrUserNotFound
		}
		logger.Error("Service:GetUserByEmail", err)
		return model.User{}, fmt.Errorf("getting user by email: %w", err)
	}
	if user.ID == uuid.Nil {
		return model.User{}, errors.ErrUserNotFound
	}
	return user, nil
}

// CreateUser create user
func (s *UserService) CreateUser(ctx context.Context, email, hashedPassword string) (model.User, error) {
	// Check if user already exists
	existingUser, err := s.GetUserByEmail(ctx, email)
	if err != nil && !stderrors.Is(err, errors.ErrUserNotFound) {
		logger.Error("Service:CreateUser", err)
		return model.User{}, fmt.Errorf("checking existing user: %w", err)
	}
	if existingUser.ID != uuid.Nil {
		return model.User{}, errors.ErrEmailAlreadyInUse
	}

	// Create new user
	return s.userRepo.CreateUser(ctx, email, hashedPassword)
}
