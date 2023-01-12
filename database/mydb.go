package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository interface {
	GetProject(ctx context.Context, id int64) (GanttProject, error)

	GetTask(ctx context.Context, id int64) (GanttTask, error)

	GetProjectTasks(ctx context.Context, id int64) ([]GanttTask, error)

	GetActivity(ctx context.Context, id int64) (GanttActivity, error)

	GetTaskActivities(ctx context.Context, id int64) ([]GanttActivity, error)

	GetProjectActivities(ctx context.Context, projectID int64) ([]GanttActivity, error)

	GetProjectAssignees(ctx context.Context, projectID int64) ([]GanttAssigned, error)

	GetAssigned(ctx context.Context, id int64) (GanttAssigned, error)

	GetActivityAssignees(ctx context.Context, activityID int64) ([]GanttAssigned, error)

	GetAssignedToUser(ctx context.Context, userID int32) ([]GanttAssigned, error)

	GetActivityComments(ctx context.Context, activityID int64) ([]GanttComment, error)

	GetProjectEmployees(ctx context.Context, projectID int64) ([]UserUser, error)

	GetProjectTeammembers(ctx context.Context, projectID int64) ([]GanttTeammember, error)

	GetProjectStates(ctx context.Context, projectID int64) ([]GanttState, error)

	ListProjects(ctx context.Context, arg ListProjectsParams) ([]GanttProject, error)

	GetState(ctx context.Context, id int64) (GanttState, error)

	GetRole(ctx context.Context, id int64) (GanttRole, error)

	GetProjectRoles(ctx context.Context, projectID int64) ([]GanttRole, error)

	GetTeam(ctx context.Context, int int64) (GanttTeam, error)

	GetProjectTeams(ctx context.Context, projectID int64) ([]GanttTeam, error)

	GetUser(ctx context.Context, id int32) (UserUser, error)

	// - - -- - - -- -- - - -- - -- - -- -- - -- batching

	ListActivitiesByTaskIDs(ctx context.Context, taskIDs []int64) ([]GanttActivity, error)

	ListTasksByProjectIDs(ctx context.Context, projectIDs []int64) ([]GanttTask, error)

	ListStateByIDS(ctx context.Context, activityIDs []int64) ([]GanttState, error)

	ListUsersIDS(ctx context.Context, ids []int32) ([]UserUser, error)

	ListAssignedsByActivityIDs(ctx context.Context, activityIDs []int64) ([]GanttAssigned, error)
}

type repoSvc struct {
	*Queries
	db *sql.DB
}

// NewRepository returns an implementation of the Repository interface.
func NewRepository(db *sql.DB) Repository {
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(20)

	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

// Open opens a database specified by the data source name.
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}

func genSqlListWithIds[T int | int32 | int64](ids []T) string {
	len := len(ids)

	if len <= 0 {
		return "()"
	}

	if len == 1 {
		return fmt.Sprintf("(%d)", ids[0])
	}
	var sb strings.Builder

	sb.WriteRune('(')
	for i := 0; i < len-1; i++ {
		sb.WriteString(strconv.FormatInt(int64(ids[i]), 10))
		sb.WriteRune(',')
	}

	sb.WriteString(strconv.FormatInt(int64(ids[len-1]), 10))
	sb.WriteRune(')')

	return sb.String()
}

func (rep *repoSvc) ListActivitiesByTaskIDs(ctx context.Context, taskIDs []int64) ([]GanttActivity, error) {
	query := fmt.Sprintf(`-- name: GetTaskActivitie	s_batch :many
	SELECT id, name, description, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, dependency_id, state_id, task_id FROM `+"`"+`gantt_activity`+"`"+`
	WHERE task_id IN %s
	`, genSqlListWithIds(taskIDs))

	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttActivity
	for rows.Next() {
		var i GanttActivity
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.PlannedStartDate,
			&i.PlannedEndDate,
			&i.PlannedBudget,
			&i.ActualStartDate,
			&i.ActualEndDate,
			&i.ActualBudget,
			&i.DependencyID,
			&i.StateID,
			&i.TaskID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *repoSvc) ListStateByIDS(ctx context.Context, activityIDs []int64) ([]GanttState, error) {
	query := fmt.Sprintf(`-- name: GetState_batch :Many
	SELECT id, name, project_id FROM gantt_state
	WHERE id IN %s
	`, genSqlListWithIds(activityIDs))

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttState
	for rows.Next() {
		var i GanttState
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ProjectID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *repoSvc) ListUsersIDS(ctx context.Context, userIDs []int32) ([]UserUser, error) {
	query := fmt.Sprintf(`-- name: GetUser_batch :many
	SELECT id, username, first_name, last_name, avatar FROM user_user
	WHERE id IN %s
	`, genSqlListWithIds(userIDs))

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserUser
	for rows.Next() {
		var i UserUser
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.FirstName,
			&i.LastName,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *repoSvc) ListAssignedsByActivityIDs(ctx context.Context, activityIDs []int64) ([]GanttAssigned, error) {
	query := fmt.Sprintf(`-- name: GetActivityAssignees_batch :many
	SELECT id, activity_id, user_id FROM gantt_assigned
	WHERE activity_id in %s
	`, genSqlListWithIds(activityIDs))

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttAssigned
	for rows.Next() {
		var i GanttAssigned
		if err := rows.Scan(&i.ID, &i.ActivityID, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *repoSvc) ListTasksByProjectIDs(ctx context.Context, projectIDs []int64) ([]GanttTask, error) {
	query := fmt.Sprintf(`-- name: GetProjectTasks_batch :many
	SELECT id, name, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, description, project_id FROM ` + "`" + `gantt_task` + "`" + `
	WHERE project_id IN %s
	`, genSqlListWithIds(projectIDs))
	
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttTask
	for rows.Next() {
		var i GanttTask
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PlannedStartDate,
			&i.PlannedEndDate,
			&i.PlannedBudget,
			&i.ActualStartDate,
			&i.ActualEndDate,
			&i.ActualBudget,
			&i.Description,
			&i.ProjectID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
