package mongo

import (
	"context"
	"fmt"
	"time"

	"clusterix-code/internal/config"
	"clusterix-code/internal/utils/di"
	internalLogger "clusterix-code/internal/utils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Provider(c *di.Container) (*mongo.Database, error) {
	cfg := di.Make[*config.Config](c)
	mongoCfg := cfg.MongoDB

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s",
		mongoCfg.Username,
		mongoCfg.Password,
		mongoCfg.Host,
		mongoCfg.Port,
		mongoCfg.AuthSource,
	)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	internalLogger.Info("MongoDB connection established successfully")

	db := client.Database(mongoCfg.Database)
	return db, nil
}
