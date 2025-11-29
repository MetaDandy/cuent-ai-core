//go:build containers

package setup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresTestDB envuelve un testcontainer de PostgreSQL y ofrece helpers
type PostgresTestDB struct {
	Container *postgres.PostgresContainer
	DSN       string
	DB        *gorm.DB
}

// SetupPostgresContainer inicia un contenedor PostgreSQL limpio para tests
func SetupPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	return container, nil
}

// SetupTestDB inicia un contenedor PostgreSQL y retorna una instancia de PostgresTestDB
func SetupTestDB(ctx context.Context) (*PostgresTestDB, error) {
	container, err := SetupPostgresContainer(ctx)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		defer func() {
			if errTerm := container.Terminate(ctx); errTerm != nil {
				log.Printf("failed to terminate container: %v", errTerm)
			}
		}()
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		defer func() {
			if errTerm := container.Terminate(ctx); errTerm != nil {
				log.Printf("failed to terminate container: %v", errTerm)
			}
		}()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresTestDB{
		Container: container,
		DSN:       dsn,
		DB:        db,
	}, nil
}

// Close termina el contenedor
func (p *PostgresTestDB) Close(ctx context.Context) error {
	return p.Container.Terminate(ctx)
}

// Migrate ejecuta las migraciones en la BD de test
// Importa config del proyecto para reutilizar las migraciones
func (p *PostgresTestDB) Migrate(ctx context.Context) error {
	// Las migraciones se ejecutarán directamente en los tests
	// usando la función Migrate del paquete config
	return nil
}

// SeedFixtures aplica fixtures iniciales a la BD
func (p *PostgresTestDB) SeedFixtures(ctx context.Context, seeders ...func(*gorm.DB) error) error {
	for _, seeder := range seeders {
		if err := seeder(p.DB); err != nil {
			return fmt.Errorf("failed to seed fixtures: %w", err)
		}
	}
	return nil
}

// CleanDB limpia todas las tablas (útil entre tests)
func (p *PostgresTestDB) CleanDB(ctx context.Context) error {
	tables := []string{
		"user_subscribeds",
		"payments",
		"assets",
		"scripts",
		"projects",
		"users",
		"generate_jobs",
		"subscriptions",
	}

	for _, table := range tables {
		if err := p.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	return nil
}

// ExecuteQuery ejecuta una query raw y retorna los resultados
func (p *PostgresTestDB) ExecuteQuery(query string, args ...interface{}) *gorm.DB {
	return p.DB.Raw(query, args...)
}
