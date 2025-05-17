package seed

import (
	"log"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) {
	if err := SeedUser(db); err != nil {
		log.Fatalf("Error al seedear usuarios: %v", err)
	}
	if err := SeedSubscriptions(db); err != nil {
		log.Fatalf("Error al seedear suscripciones: %v", err)
	}
}
