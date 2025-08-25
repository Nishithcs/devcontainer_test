package db

import (
	"clusterix-code/internal/data/db/migrations"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Register migrations here in the order in which they should be run.
// It is worth noting that two tables migrate automatically:
// 1. the migrations table when migrate up is run for the first time.
// 2. the casbin_rule table when the access control service is initialized.
// The casbin_rule table is created by the casbin gorm library and is used to store policies.
// The migrations specified here only run when the migrate up command is run
var migrationsRegistry = []migrations.Migrant{
	&migrations.CreateUsersTable{},
	&migrations.CreateProvidersTable{},
	&migrations.CreateMachineConfigsTable{},
	&migrations.CreateRepositoriesTable{},
	&migrations.CreateGitPersonalAccessTokensTable{},
	&migrations.CreateWorkspacesTable{},
	&migrations.AddCreatedByAndIconToProviders{},
	&migrations.RemoveProviderIdFromRepositories{},
	&migrations.AddProviderIdToWorkspaces{},
	&migrations.CreateWorkspaceStatusEventsTable{},
	&migrations.AddFingerprintToWorkspaces{},
	&migrations.CreateWorkspaceConfigsTable{},
	&migrations.AddWorkspaceConfigIdToWorkspaces{},
}

func findMigrationForRollback(name string) (migrations.Migrant, error) {
	migrants := migrationsRegistry
	for i := len(migrants) - 1; i >= 0; i-- {
		migrant := migrants[i]
		if migrant.GetName() == name {
			return migrant, nil
		}
	}
	return nil, fmt.Errorf("migration %s not found", name)
}

func RunMigrateMake(structName string) error {
	err := migrations.GenerateMigrationFile(structName)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	return nil
}

func RunMigrateUp(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	migrants := migrationsRegistry
	appliedCount := 0
	for _, migrant := range migrants {
		migrationName := migrant.GetName()
		var existing Migration

		if err := db.Where("name = ?", migrationName).First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			baseMigration := migrations.NewBaseMigration(db, migrant)

			if err := baseMigration.MigrateUp(); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migrationName, err)
			}

			if err := db.Create(&Migration{
				Name:      migrationName,
				AppliedAt: time.Now(),
			}).Error; err != nil {
				return fmt.Errorf("failed to record migrations %s: %w", migrationName, err)
			}
			appliedCount++
		} else {
			fmt.Printf("\nMigration %s already applied at %s\n\n", migrationName, existing.AppliedAt)
		}
	}
	if appliedCount == 0 {
		fmt.Println("\nNo new migrations to apply")
	} else {
		fmt.Printf("\n%d migration(s) applied successfully\n", appliedCount)
	}
	return nil
}

func RunMigrateDown(db *gorm.DB, migrationDownCount int) error {
	if !db.Migrator().HasTable(&Migration{}) {
		return fmt.Errorf("migrations table does not exist")
	}
	// check if the migrations table is empty, if so, return an error
	var count int64
	if err := db.Model(&Migration{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count migrations: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("no migrations to rollback")
	}
	// get the last migrationDownCount migrations
	var dbMigrations []Migration
	if err := db.Order("applied_at desc").Limit(migrationDownCount).Find(&dbMigrations).Error; err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}
	// alert the user before proceeding, ask for input
	fmt.Printf("The following migrations will be rolled back:\n\n")
	for _, migration := range dbMigrations {
		fmt.Printf("  - %s\n", migration.Name)
	}

	fmt.Printf("\nAre you sure you want to proceed? (yes/no): ")
	var input string

	if _, err := fmt.Scanln(&input); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	if input != "yes" {
		return fmt.Errorf("rollback aborted")
	}

	// get the migration struct for each migration and run the migration down
	for _, migration := range dbMigrations {
		migrant, err := findMigrationForRollback(migration.Name)
		if err != nil {
			fmt.Printf("migration %s not found in migrations registry\n\n", migration.Name)
			continue
		}
		baseMigration := migrations.NewBaseMigration(db, migrant)
		if err := baseMigration.MigrateDown(); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
		}
		if err := db.Delete(&migration).Error; err != nil {
			return fmt.Errorf("failed to delete migration %s: %w", migration.Name, err)
		}
		fmt.Printf("\nmigration %s rolled back successfully\n\n", migration.Name)
	}
	return nil
}
