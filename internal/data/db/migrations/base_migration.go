package migrations

import (
	"fmt"
	"os"
	"strings"
	"time"

	"clusterix-code/internal/utils/helpers"

	"gorm.io/gorm"
)

const migrationStructStub = `
package migrations

type :structName struct {
	BaseMigration
	Name string
}

func (m *:structName) UpSql() string {
	return "ENTER_SQL_HERE"
}

func (m *:structName) DownSql() string {
	return "ENTER_SQL_HERE"
}

func (m *:structName) GetName() string {
	// don't change this after the migration is applied
	return ":migrationName"
}
`

func GenerateMigrationFile(structName string) error {
	structName = strings.ToUpper(structName[:1]) + structName[1:]
	now := time.Now()
	migrationName := fmt.Sprintf("%d_%02d_%02d_%d_%s",
		now.Year(),
		now.Month(),
		now.Day(),
		now.UnixMilli(),
		helpers.PascalToSnakeCase(structName),
	)

	stub := strings.ReplaceAll(migrationStructStub, ":structName", structName)
	stub = strings.ReplaceAll(stub, ":migrationName", migrationName)
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationsDir := currDir + "/internal/common/db/migrations"
	fileName := migrationsDir + "/" + migrationName + ".go"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(stub)
	if err != nil {
		return err
	}
	fmt.Printf("Migration file created at %s\n", fileName)
	return nil
}

type Migrant interface {
	UpSql() string
	DownSql() string
	GetName() string
}

type BaseMigration struct {
	db      *gorm.DB
	migrant Migrant
}

func NewBaseMigration(db *gorm.DB, migrant Migrant) *BaseMigration {
	return &BaseMigration{
		db:      db,
		migrant: migrant,
	}
}

func (b *BaseMigration) execQuery(sqlQuery string) error {
	err := b.db.Exec(sqlQuery).Error

	if err != nil {
		return err
	}

	return nil
}

func (b *BaseMigration) MigrateUp() error {
	query := b.migrant.UpSql()
	return b.execQuery(query)
}

func (b *BaseMigration) MigrateDown() error {
	query := b.migrant.DownSql()
	return b.execQuery(query)
}
