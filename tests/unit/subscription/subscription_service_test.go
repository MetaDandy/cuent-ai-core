//go:build unit

package subscription_test

import (
	"testing"
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/core/subscription"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
)

// Test de paginación - lógica del servicio
func TestPaginationLogic(t *testing.T) {
	tests := []struct {
		name          string
		total         int64
		limit         uint
		expectedPages uint
	}{
		{"5 items, limit 10", 5, 10, 1},
		{"5 items, limit 2", 5, 2, 3}, // (5 + 2 - 1) / 2 = 3
		{"10 items, limit 10", 10, 10, 1},
		{"11 items, limit 10", 11, 10, 2}, // (11 + 10 - 1) / 10 = 2
		{"0 items, limit 10", 0, 10, 0},
		{"1 item, limit 1", 1, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := uint((tt.total + int64(tt.limit) - 1) / int64(tt.limit))
			if pages != tt.expectedPages {
				t.Errorf("expected pages %d, got %d", tt.expectedPages, pages)
			}
		})
	}
}

// Test de DTO conversion
func TestSubscriptionToDTO(t *testing.T) {
	subID := uuid.New()
	now := time.Now()
	mockSub := &model.Subscription{
		ID:         subID,
		Name:       "Premium",
		Cuentokens: 10000,
		Duration:   now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	dto := subscription.SubscriptionToDTO(mockSub)

	if dto.ID != subID.String() {
		t.Errorf("expected ID %s, got %s", subID.String(), dto.ID)
	}
	if dto.Name != "Premium" {
		t.Errorf("expected Name Premium, got %s", dto.Name)
	}
	if dto.Cuentokens != 10000 {
		t.Errorf("expected Cuentokens 10000, got %d", dto.Cuentokens)
	}
}

// Test de conversión de lista a DTOs
func TestSubscriptionToListDTO(t *testing.T) {
	subs := []model.Subscription{
		{
			ID:         uuid.New(),
			Name:       "Free",
			Cuentokens: 1000,
			Duration:   time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			Name:       "Premium",
			Cuentokens: 10000,
			Duration:   time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	dtos := subscription.SubscriptionToListDTO(subs)

	if len(dtos) != len(subs) {
		t.Errorf("expected %d DTOs, got %d", len(subs), len(dtos))
	}

	if dtos[0].Name != "Free" {
		t.Errorf("expected first DTO name Free, got %s", dtos[0].Name)
	}

	if dtos[1].Name != "Premium" {
		t.Errorf("expected second DTO name Premium, got %s", dtos[1].Name)
	}
}

