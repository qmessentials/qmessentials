package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDummyHashCost(t *testing.T) {
	cost, err := bcrypt.Cost(dummyHash)
	if err != nil || cost != bcryptCost {
		t.Fatalf("dummy bcrypt cost: %d, %v", cost, err)
	}
}
