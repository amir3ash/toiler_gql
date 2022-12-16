package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// type Repository interface {
// 	GetActivity(ctx context.Context, id int64) (GanttActivity, error)

// 	GetActivityComments(ctx context.Context, activityID int64) ([]GanttComment, error)

// 	GetAssignedToUser(ctx context.Context, userID int32) ([]GanttAssigned, error)

// 	GetProjectActivities(ctx context.Context, projectID int64) ([]GanttActivity, error)

// 	GetProjectAssignees(ctx context.Context, projectID int64) ([]GanttAssigned, error)

// 	GetProjectEmployees(ctx context.Context, projectID int64) ([]UserUser, error)

// 	GetStates(ctx context.Context, projectID int64) ([]GanttState, error)

// 	ListProjects(ctx context.Context, arg ListProjectsParams) ([]GanttProject, error)
// }

// type repoSvc struct {
// 	*Queries
// 	db *sql.DB
// }

// NewRepository returns an implementation of the Repository interface.
func NewRepository(db *sql.DB) *Queries {
	return &Queries{
		// Queries: New(db),
		db: db,
	}
}

// Open opens a database specified by the data source name.
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}

