package commands

import (
	"clusterix-code/internal/config"
	"clusterix-code/internal/data/db"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/helpers"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/mongo"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var MigrateUpCmd = &cobra.Command{
	Use:   "migrate:up",
	Short: "Apply defined migrations on the DB",
	Run: func(cmd *cobra.Command, args []string) {
		defer logger.Sync()
		postgresDb := cmdBootstrap()
		err := db.RunMigrateUp(postgresDb)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

var MigrateDownCmd = &cobra.Command{
	Use:   "migrate:down [migration-rollback-count]",
	Short: "Rollback already run migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrationDownCount := 1
		if len(args) > 0 {
			var err error
			migrationDownCount, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
		}
		fmt.Println("Migration down count: ", migrationDownCount)
		defer logger.Sync()
		postgresDb := cmdBootstrap()
		err := db.RunMigrateDown(postgresDb, migrationDownCount)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func cmdBootstrap() *gorm.DB {
	helpers.LoadEnv()
	logger.Init(os.Getenv("APP_ENV"))

	c := di.NewContainer(0)

	di.Register(c, config.Provider)
	di.Register(c, db.Provider)
	di.Register(c, mongo.Provider)

	c.Bootstrap()

	return di.Make[*gorm.DB](c)
}
