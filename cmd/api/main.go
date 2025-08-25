package main

import (
	"clusterix-code/internal/api/server"
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/constants"
	"clusterix-code/internal/consumers"
	workspaceConsumers "clusterix-code/internal/consumers/workspace"
	"clusterix-code/internal/data/db"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/helpers"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/mongo"
	"clusterix-code/internal/utils/rabbitmq"
	"os"
)

func main() {
	helpers.LoadEnv()
	logger.Init(os.Getenv("APP_ENV"))
	defer logger.Sync()

	c := di.NewContainer(0)

	di.Register(c, config.Provider)
	di.Register(c, db.Provider)
	di.Register(c, rabbitmq.Provider)
	di.Register(c, mongo.Provider)
	di.Register(c, repositories.Provider)
	di.Register(c, api_clients.Provider)
	di.Register(c, services.Provider)
	di.Register(c, server.Provider)
	di.Register(c, consumers.Provider)

	c.Bootstrap()

	consumerManager := di.Make[*consumers.ConsumerManager](c)
	consumerManager.RegisterConsumer(constants.WORKSPACE_LOG_HANDLER_QUEUE, workspaceConsumers.NewWorkspaceLogHandlerConsumer(c))
	go consumerManager.Start()

	server := di.Make[*server.Server](c)
	server.Run()
}
