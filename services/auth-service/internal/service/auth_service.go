package service

import (
	"context"
	"errors"
	"time"

	"github.com/MiltonJ23/Daedalus/AuthService/internal/domain"
	"github.com/MiltonJ23/Daedalus/AuthService/pkg/crypto"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInternal           = errors.New("internal server error")
)

type AuthService struct {
	repo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	// at first, let's hash the password
	hash, hashingErr := crypto.HashPassword(req.Password)
	if hashingErr != nil {
		return nil, ErrInternal
	}

	now := time.Now().UTC()

	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		if err.Error() == "user already exists" {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrInternal
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.User, error) {
	user, fetchUserErr := s.repo.FindByEmail(ctx, req.Email)
	if fetchUserErr != nil {
		return nil, ErrInternal
	}

	if user == nil {
		crypto.DummyCheck(req.Password)
		return nil, ErrInvalidCredentials
	}

	if !crypto.CheckPasswordHash(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
