package seed

import (
	"log"
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedSubscriptions inserta los planes básicos solo
// si la tabla está vacía.
func SeedSubscriptions(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Subscription{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("⚠️ Las suscripciones ya existen.")
		return nil // ya hay datos; no repetir
	}

	// por claridad, Duration representa “vigencia del plan” (días)
	// usa time.Date con año 0 para evitar confusiones de zona horaria
	makeDuration := func(days int) time.Time {
		return time.Date(0, 1, days, 0, 0, 0, 0, time.UTC)
	}

	plans := []model.Subscription{
		{
			ID:         uuid.New(),
			Name:       "Free",
			Cuentokens: "1000",
			Duration:   makeDuration(30),
		},
		{
			ID:         uuid.New(),
			Name:       "Standard", // plan base
			Cuentokens: "5000",
			Duration:   makeDuration(30),
		},
		{
			ID:         uuid.New(),
			Name:       "Pro",
			Cuentokens: "25000",
			Duration:   makeDuration(30),
		},
		{
			ID:         uuid.New(),
			Name:       "Enterprise",
			Cuentokens: "100000",
			Duration:   makeDuration(30),
		},
	}
	if err := db.Create(&plans).Error; err != nil {
		log.Fatalf("❌ Error creando las suscripciones: %v", err)
	}

	log.Printf("✅ Suscripciones creadas correctamente.")
	return nil
}
