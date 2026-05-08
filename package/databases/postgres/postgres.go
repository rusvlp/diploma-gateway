package postgres

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func buildPostgresURL(host, port, user, password, dbname string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func RunMigrations(migrationsPath, host, port, user, password, dbname string) error {
	log.Printf("[migrations] path=%s  db=%s@%s:%s/%s", migrationsPath, user, host, port, dbname)

	if !folderExists(migrationsPath) {
		return fmt.Errorf("migrations folder not found: %s", migrationsPath)
	}

	dsn := buildPostgresURL(host, port, user, password, dbname)
	source := "file://" + migrationsPath

	const maxAttempts = 10
	var m *migrate.Migrate
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		m, err = migrate.New(source, dsn)
		if err == nil {
			break
		}
		wait := time.Duration(attempt) * 2 * time.Second
		log.Printf("[migrations] attempt %d/%d failed to connect: %v — retrying in %s", attempt, maxAttempts, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		return fmt.Errorf("[migrations] could not connect after %d attempts: %w", maxAttempts, err)
	}
	defer m.Close()

	log.Println("[migrations] applying...")

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("[migrations] already up to date")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[migrations] failed: %w", err)
	}

	log.Println("[migrations] done")
	return nil
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
