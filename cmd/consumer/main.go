package main

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/constants"
	"clusterix-code/internal/consumers"
	userConsumers "clusterix-code/internal/consumers/user"
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
	di.Register(c, consumers.Provider)

	c.Bootstrap()

	consumerManager := di.Make[*consumers.ConsumerManager](c)
	consumerManager.RegisterConsumer(constants.AUTH_USER_CREATED_QUEUE, userConsumers.NewAuthUserCreatedConsumer(c))
	consumerManager.RegisterConsumer(constants.AUTH_USER_UPDATED_QUEUE, userConsumers.NewAuthUserUpdatedConsumer(c))
	consumerManager.RegisterConsumer(constants.AUTH_USER_DELETED_QUEUE, userConsumers.NewAuthUserDeletedConsumer(c))
	consumerManager.Start()
}
