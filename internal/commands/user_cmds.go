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
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"time"
)

var SyncUsersCmd = &cobra.Command{
	Use:   "sync-users",
	Short: "Sync users from auth service",
	Run: func(cmd *cobra.Command, args []string) {
		SyncUsers(cmd, args)
	},
}

func SyncUsers(cmd *cobra.Command, args []string) {
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

	services := di.Make[*services.Services](c)
	apiClients := di.Make[*api_clients.APIClients](c)
	authApiClient := apiClients.Auth
	context := cmd.Context()

	currentPage := 1
	totalPage := 1
	var batch []uint64
	for currentPage <= totalPage {
		fmt.Println("Current Page: ", currentPage)
		time.Sleep(3 * time.Second)
		users, err := authApiClient.GetUsers(context, currentPage)
		if err != nil {
			logger.Error("Failed to get users", err)
			return
		}

		totalPage = users.LastPage
		for _, user := range users.Data {
			userModel, err := services.User.SyncUser(context, &user)
			if err != nil {
				logger.Error("Failed to sync a user: ", err, zap.Any("user", user))
				continue
			}

			batch = append(batch, userModel.ID)
		}
		currentPage++
	}

	fmt.Println("Users synced successfully")
}
