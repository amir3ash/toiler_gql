package graph

import (
	// "context"
	"toiler-graphql/cache"
	"toiler-graphql/database"
	"toiler-graphql/dataloaders"
	// "toiler-graphql/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository  database.Repository
	Dataloaders dataloaders.Retriever
	Cache       cache.Cache
}

// func (r *Resolver) GetActivity(ctx context.Context, id int64) (database.GanttActivity, error)

// func (r* Resolver) GetActivityComments(ctx context.Context, activityID int64) ([]GanttComment, error)

// func (r* Resolver) GetAssignedToUser(ctx context.Context, userID int32) ([]GanttAssigned, error)

// func (r* Resolver) GetProjectActivities(ctx context.Context, projectID int64) ([]GanttActivity, error)

// func (r* Resolver) GetProjectAssignees(ctx context.Context, projectID int64) ([]GanttAssigned, error)

// func (r* Resolver) GetProjectEmployees(ctx context.Context, projectID int64) ([]UserUser, error)

// func (r* Resolver) GetStates(ctx context.Context, projectID int64) ([]GanttState, error)

// func (r* Resolver) ListProjects(ctx context.Context, arg ListProjectsParams) ([]GanttProject, error)
