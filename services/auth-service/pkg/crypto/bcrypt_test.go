package crypto_test

import (
	"testing"

	"github.com/MiltonJ23/Daedalus/AuthService/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptImplementation(t *testing.T) {
	password := "Daedalus2026!"

	// Test de génération
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Vérification du cost factor (Doit être 12 selon NFR-SEC-02)
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("Failed to extract cost: %v", err)
	}
	if cost != crypto.BcryptCost {
		t.Errorf("Expected cost %d, got %d", crypto.BcryptCost, cost)
	}

	// Test de vérification correcte
	if !crypto.CheckPasswordHash(hash, password) {
		t.Error("CheckPasswordHash failed for correct password")
	}

	// Test de vérification incorrecte
	if crypto.CheckPasswordHash("WrongPassword!", hash) {
		t.Error("CheckPasswordHash succeeded for wrong password")
	}
}
