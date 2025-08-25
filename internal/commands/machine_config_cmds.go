package commands

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/data/db"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/helpers"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/mongo"
	"clusterix-code/internal/utils/rabbitmq"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var ImportMachineConfigsCmd = &cobra.Command{
	Use:   "import-machine-configs",
	Short: "import machine configs into the database",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LoadEnv()
		logger.Init(os.Getenv("APP_ENV"))
		defer logger.Sync()

		c := di.NewContainer(0)

		di.Register(c, config.Provider)
		di.Register(c, db.Provider)
		di.Register(c, mongo.Provider)
		di.Register(c, rabbitmq.Provider)
		di.Register(c, repositories.Provider)
		di.Register(c, services.Provider)
		di.Register(c, api_clients.Provider)

		c.Bootstrap()

		service := di.Make[*services.Services](c)
		ctx := cmd.Context()
		if err := service.MachineConfig.Import(ctx); err != nil {
			log.Fatalf("Failed to import machine configs: %v", err)
		}

		log.Println("Successfully imported machine configs")
	},
}
