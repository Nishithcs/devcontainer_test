package main

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/data/db"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/jobs"
	"clusterix-code/internal/services"
	"clusterix-code/internal/tasks"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/helpers"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/mongo"
	"clusterix-code/internal/utils/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func main() {
	helpers.LoadEnv()
	logger.Init(os.Getenv("APP_ENV"))
	defer logger.Sync()

	c := di.NewContainer(0)

	di.Register(c, config.Provider)
	di.Register(c, db.Provider)
	di.Register(c, mongo.Provider)
	di.Register(c, rabbitmq.Provider)
	di.Register(c, repositories.Provider)
	di.Register(c, api_clients.Provider)
	di.Register(c, services.Provider)
	c.Bootstrap()

	services := di.Make[*services.Services](c)

	//go devpod.StartPersistentSupervisor(context.Background(), repos.Workspace, workspaceSvc.UpdateWorkspaceStatus)

	server := jobs.NewAsynqServer()
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TaskStartWorkspace, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.StartWorkspacePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		log.Printf("🛠 Processing workspace starting job for workspace ID %d", p.WorkspaceID)
		return jobs.HandleStartWorkspaceTask(ctx, t, services.Workspace, services.Publisher, services.WorkspaceConfig)
	})

	mux.HandleFunc(tasks.TaskStopWorkspace, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.StopWorkspacePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		log.Printf("🛠 Processing workspace stopping job for workspace ID %d", p.WorkspaceID)
		return jobs.HandleStopWorkspaceTask(ctx, t, services.Workspace, services.Publisher)
	})

	mux.HandleFunc(tasks.TaskRestartWorkspace, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.RestartWorkspacePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		log.Printf("🛠 Processing workspace restarting job for workspace ID %d", p.WorkspaceID)
		return jobs.HandleRestartWorkspaceTask(ctx, t, services.Workspace, services.Publisher)
	})

	mux.HandleFunc(tasks.TaskRebuildWorkspace, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.RebuildWorkspacePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		log.Printf("🛠 Processing workspace rebuilding job for workspace ID %d", p.WorkspaceID)
		return jobs.HandleRebuildWorkspaceTask(ctx, t, services.Workspace, services.Publisher)
	})

	mux.HandleFunc(tasks.TaskTerminateWorkspace, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.StartWorkspacePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		log.Printf("🛠 Processing workspace terminattiog job for workspace ID %d", p.WorkspaceID)
		return jobs.HandleTerminateWorkspaceTask(ctx, t, services.Workspace, services.Publisher)
	})

	log.Println("🚀 Worker starting to process jobs...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("❌ Could not start worker server: %v", err)
	}
}
