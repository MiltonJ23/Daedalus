package service

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MiltonJ23/Daedalus/AuthService/internal/domain"
	"github.com/MiltonJ23/Daedalus/AuthService/pkg/crypto"
)

type mockUserRepository struct {
	users map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil // Comportement normal si non trouvé
	}
	return user, nil
}

// --- Tests Unitaires ---

func TestAuthService_Register(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo)
	ctx := context.Background()

	req := domain.RegisterRequest{
		Email:    "test@daedalus.cm",
		Password: "SecurePassword123!",
	}

	// Test 1: Création réussie
	user, err := authService.Register(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}

	// Test 2: Doublon (Doit renvoyer ErrUserAlreadyExists)
	_, err = authService.Register(ctx, req)
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

// Test crucial pour vérifier l'exigence FR-AUTH-06
func TestAuthService_Login_ConstantTiming(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo)
	ctx := context.Background()

	validEmail := "exist@daedalus.cm"
	invalidEmail := "notfound@daedalus.cm"
	password := "CommonPassword123!"

	// Insertion préalable d'un utilisateur valide
	hash, _ := crypto.HashPassword(password)
	_ = repo.Create(ctx, &domain.User{
		ID:           "123",
		Email:        validEmail,
		PasswordHash: hash,
	})

	// Mesure 1 : Login avec un email existant, mais mauvais mot de passe
	reqExist := domain.LoginRequest{Email: validEmail, Password: "WrongPassword!"}
	start1 := time.Now()
	_, err1 := authService.Login(ctx, reqExist)
	durationExist := time.Since(start1)

	if err1 != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err1)
	}

	// Mesure 2 : Login avec un email INEXISTANT (Tentative d'énumération)
	reqNotExist := domain.LoginRequest{Email: invalidEmail, Password: "WrongPassword!"}
	start2 := time.Now()
	_, err2 := authService.Login(ctx, reqNotExist)
	durationNotExist := time.Since(start2)

	if err2 != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err2)
	}

	// Évaluation de l'écart type (Tolerance max de 30ms due aux variations du CPU)
	diff := math.Abs(float64(durationExist.Milliseconds() - durationNotExist.Milliseconds()))

	if diff > 30 {
		t.Errorf("TIMING VULNERABILITY DETECTED! \nExist user took: %v \nNot exist user took: %v \nDifference: %v ms (Exceeds 30ms tolerance)", durationExist, durationNotExist, diff)
	}
}
