// internal/repositories/workspace_log_repository.go

package repositories

import (
	"clusterix-code/internal/data/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type WorkspaceLogRepository struct {
	collection *mongo.Collection
}

func NewWorkspaceLogRepository(db *mongo.Database) *WorkspaceLogRepository {
	return &WorkspaceLogRepository{
		collection: db.Collection("workspace_logs"),
	}
}

func (r *WorkspaceLogRepository) GetLatestLogs(ctx context.Context, workspaceID uint64) ([]*models.WorkspaceLog, error) {
	filter := bson.M{"workspaceid": workspaceID}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdat", Value: 1}}).
		SetLimit(200)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.WorkspaceLog
	for cursor.Next(ctx) {
		var log models.WorkspaceLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *WorkspaceLogRepository) Create(ctx context.Context, log *models.WorkspaceLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	return err
}
