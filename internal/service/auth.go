package service

import (
	"context"
	"mpc/internal/model"
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"mpc/pkg/token"
	"mpc/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	tokenManager  *token.TokenManager
	walletService *WalletService
	userService   *UserService
}

func NewAuthService(userService *UserService, walletService *WalletService, tokenManager *token.TokenManager) *AuthService {
	return &AuthService{
		userService:   userService,
		walletService: walletService,
		tokenManager:  tokenManager,
	}
}

// Login login
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (model.AuthResponse, error) {
	// Get user by email
	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("Service:Login", err)
		return model.AuthResponse{}, err
	}
	if user.ID == uuid.Nil {
		return model.AuthResponse{}, errors.ErrUserNotFound
	}

	// Check if password is correct
	if !comparePassword(user.PasswordHash, req.Password) {
		return model.AuthResponse{}, errors.ErrInvalidPassword
	}

	// Get primary wallet
	wallet, err := s.getPrimaryWallet(ctx, user.ID)
	if err != nil {
		logger.Error("Service:Login", err)
		return model.AuthResponse{}, err
	}
	if wallet.ID == uuid.Nil {
		return model.AuthResponse{}, errors.ErrWalletNotFound
	}

	// Generate tokens
	token, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		logger.Error("Service:Login", err)
		return model.AuthResponse{}, err
	}

	return s.createAuthResponse(user, wallet, token), nil
}

// Signup signup
func (s *AuthService) Signup(ctx context.Context, req *model.SignupRequest) (model.AuthResponse, error) {
	// Check if email already exists
	if _, err := s.userService.GetUserByEmail(ctx, req.Email); err == nil {
		return model.AuthResponse{}, errors.ErrEmailAlreadyInUse
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		logger.Error("Service:Signup", err)
		return model.AuthResponse{}, err
	}

	// Create user and wallet
	user, err := s.userService.CreateUser(ctx, req.Email, hashedPassword)
	if err != nil {
		logger.Error("Service:Signup", err)
		return model.AuthResponse{}, err
	}

	wallet, err := s.walletService.CreateWallet(ctx, user.ID)
	if err != nil {
		logger.Error("Service:Signup", err)
		return model.AuthResponse{}, err
	}

	// Generate tokens
	token, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		logger.Error("Service:Signup", err)
		return model.AuthResponse{}, err
	}

	return s.createAuthResponse(user, wallet, token), nil
}

// Refresh refresh
func (s *AuthService) Refresh(ctx context.Context, req *model.Refresh) (model.RefreshResponse, error) {
	// Verify refresh token
	userID, err := s.tokenManager.VerifyToken(ctx, req.RefreshToken, token.TokenTypeRefresh)
	if err != nil {
		logger.Error("Service:Refresh", err)
		return model.RefreshResponse{}, err
	}

	// Generate new tokens
	token, err := s.generateTokenPair(ctx, userID)
	if err != nil {
		logger.Error("Service:Refresh", err)
		return model.RefreshResponse{}, err
	}

	return model.RefreshResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// getPrimaryWallet get primary wallet
func (s *AuthService) getPrimaryWallet(ctx context.Context, userID uuid.UUID) (model.Wallet, error) {
	wallet, err := s.walletService.GetWalletByUserID(ctx, userID)
	if err != nil {
		logger.Error("Service:GetPrimaryWallet", err)
		return model.Wallet{}, err
	}
	if wallet.ID == uuid.Nil {
		return model.Wallet{}, errors.ErrWalletNotFound
	}
	return wallet, nil
}

// createAuthResponse create auth response
func (s *AuthService) createAuthResponse(user model.User, wallet model.Wallet, token token.Token) model.AuthResponse {
	return model.AuthResponse{
		User:         utils.ToUserResponse(user),
		Wallet:       utils.ToWalletResponse(wallet),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

// generateTokenPair generate token pair
func (s *AuthService) generateTokenPair(ctx context.Context, userID uuid.UUID) (token.Token, error) {
	return s.tokenManager.GenerateTokenPair(ctx, userID)
}

// hashPassword hash password
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// comparePassword compare password
func comparePassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
