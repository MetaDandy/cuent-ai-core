package seed

import (
	"log"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminUser asegura que exista un usuario administrador.
func SeedUser(db *gorm.DB) error {
	const (
		adminName     = "Administrador"
		adminEmail    = "admin@gmail.com"
		adminPassword = "changeme123" // ⇢ luego cámbiala en producción
	)

	// 1) ¿Ya existe?
	var existing model.User
	err := db.Where("email = ?", adminEmail).First(&existing).Error
	if err == nil {
		log.Printf("⚠️ Usuario %q (%s) ya existe; skip.", existing.Name, existing.Email)
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		log.Fatalf("❌ Error buscando usuario admin: %v", err)
	}

	// 2) Generar hash de la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Error generando hash de contraseña: %v", err)
	}

	// 3) Construir registro
	user := model.User{
		ID:       uuid.New(),
		Name:     adminName,
		Email:    adminEmail,
		Password: string(hash),
	}

	// 4) Crear en una transacción
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// Ejemplo: si quisieras asociar roles aquí, hazlo dentro del mismo tx.
		return nil
	}); err != nil {
		log.Fatalf("❌ Error creando usuario admin: %v", err)
	}

	log.Printf("✅ Usuario admin %q creado con email %s.", adminName, adminEmail)
	return nil
}
