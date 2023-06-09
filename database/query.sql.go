// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package database

import (
	"context"
)

const getActivity = `-- name: GetActivity :one
SELECT id, name, description, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, dependency_id, state_id, task_id FROM ` + "`" + `gantt_activity` + "`" + `
WHERE id = ?
`

func (q *Queries) GetActivity(ctx context.Context, id int64) (GanttActivity, error) {
	row := q.db.QueryRowContext(ctx, getActivity, id)
	var i GanttActivity
	err := row.Scan(
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
	)
	return i, err
}

const getActivityAssignees = `-- name: GetActivityAssignees :many
SELECT a.id, a.activity_id, a.user_id FROM ` + "`" + `gantt_assigned` + "`" + ` a
WHERE a.activity_id = ?
`

func (q *Queries) GetActivityAssignees(ctx context.Context, activityID int64) ([]GanttAssigned, error) {
	rows, err := q.db.QueryContext(ctx, getActivityAssignees, activityID)
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

const getActivityComments = `-- name: GetActivityComments :many
SELECT id, created_at, updated_at, text, activity_id, author_id FROM ` + "`" + `gantt_comment` + "`" + `
WHERE activity_id = ?
`

func (q *Queries) GetActivityComments(ctx context.Context, activityID int64) ([]GanttComment, error) {
	rows, err := q.db.QueryContext(ctx, getActivityComments, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttComment
	for rows.Next() {
		var i GanttComment
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Text,
			&i.ActivityID,
			&i.AuthorID,
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

const getAssigned = `-- name: GetAssigned :one
SELECT id, activity_id, user_id FROM ` + "`" + `gantt_assigned` + "`" + `
WHERE id = ?
`

func (q *Queries) GetAssigned(ctx context.Context, id int64) (GanttAssigned, error) {
	row := q.db.QueryRowContext(ctx, getAssigned, id)
	var i GanttAssigned
	err := row.Scan(&i.ID, &i.ActivityID, &i.UserID)
	return i, err
}

const getAssignedToUser = `-- name: GetAssignedToUser :many
SELECT id, activity_id, user_id FROM ` + "`" + `gantt_assigned` + "`" + `
WHERE user_id = ?
`

func (q *Queries) GetAssignedToUser(ctx context.Context, userID int32) ([]GanttAssigned, error) {
	rows, err := q.db.QueryContext(ctx, getAssignedToUser, userID)
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

const getProject = `-- name: GetProject :one
SELECT id, name, planned_start_date, planned_end_date, actual_start_date, actual_end_date, description, project_manager_id FROM ` + "`" + `gantt_project` + "`" + `
WHERE id = ?
`

func (q *Queries) GetProject(ctx context.Context, id int64) (GanttProject, error) {
	row := q.db.QueryRowContext(ctx, getProject, id)
	var i GanttProject
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PlannedStartDate,
		&i.PlannedEndDate,
		&i.ActualStartDate,
		&i.ActualEndDate,
		&i.Description,
		&i.ProjectManagerID,
	)
	return i, err
}

const getProjectActivities = `-- name: GetProjectActivities :many
SELECT ac.id, ac.name, ac.description, ac.planned_start_date, ac.planned_end_date, ac.planned_budget, ac.actual_start_date, ac.actual_end_date, ac.actual_budget, ac.dependency_id, ac.state_id, ac.task_id FROM ` + "`" + `gantt_activity` + "`" + ` ac
JOIN ` + "`" + `gantt_task` + "`" + ` t ON ac.task_id = t.id
WHERE t.project_id = ?
`

func (q *Queries) GetProjectActivities(ctx context.Context, projectID int64) ([]GanttActivity, error) {
	rows, err := q.db.QueryContext(ctx, getProjectActivities, projectID)
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

const getProjectAssignees = `-- name: GetProjectAssignees :many
SELECT a.id, a.activity_id, a.user_id FROM ` + "`" + `gantt_assigned` + "`" + ` a
JOIN ` + "`" + `gantt_activity` + "`" + ` ac ON a.activity_id = ac.id
JOIN ` + "`" + `gantt_task` + "`" + ` t ON ac.task_id = t.id
WHERE t.project_id = ?
`

func (q *Queries) GetProjectAssignees(ctx context.Context, projectID int64) ([]GanttAssigned, error) {
	rows, err := q.db.QueryContext(ctx, getProjectAssignees, projectID)
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

const getProjectEmployees = `-- name: GetProjectEmployees :many
SELECT u.id, u.username, u.first_name, u.last_name, u.avatar FROM ` + "`" + `user_user` + "`" + ` u
JOIN ` + "`" + `gantt_teammember` + "`" + ` tm ON tm.user_id = u.id
JOIN ` + "`" + `gantt_team` + "`" + ` t ON t.id = tm.team_id
WHERE t.project_id = ?
`

func (q *Queries) GetProjectEmployees(ctx context.Context, projectID int64) ([]UserUser, error) {
	rows, err := q.db.QueryContext(ctx, getProjectEmployees, projectID)
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

const getProjectRoles = `-- name: GetProjectRoles :many
SELECT id, name, project_id FROM ` + "`" + `gantt_role` + "`" + `
WHERE project_id = ?
`

func (q *Queries) GetProjectRoles(ctx context.Context, projectID int64) ([]GanttRole, error) {
	rows, err := q.db.QueryContext(ctx, getProjectRoles, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttRole
	for rows.Next() {
		var i GanttRole
		if err := rows.Scan(&i.ID, &i.Name, &i.ProjectID); err != nil {
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

const getProjectStates = `-- name: GetProjectStates :many
SELECT id, name, project_id FROM ` + "`" + `gantt_state` + "`" + `
WHERE project_id = ?
`

func (q *Queries) GetProjectStates(ctx context.Context, projectID int64) ([]GanttState, error) {
	rows, err := q.db.QueryContext(ctx, getProjectStates, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttState
	for rows.Next() {
		var i GanttState
		if err := rows.Scan(&i.ID, &i.Name, &i.ProjectID); err != nil {
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

const getProjectTasks = `-- name: GetProjectTasks :many
SELECT id, name, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, description, project_id FROM ` + "`" + `gantt_task` + "`" + `
WHERE project_id = ?
`

func (q *Queries) GetProjectTasks(ctx context.Context, projectID int64) ([]GanttTask, error) {
	rows, err := q.db.QueryContext(ctx, getProjectTasks, projectID)
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

const getProjectTeammembers = `-- name: GetProjectTeammembers :many
SELECT tm.id, tm.role_id, tm.team_id, tm.user_id FROM ` + "`" + `gantt_teammember` + "`" + ` tm
JOIN ` + "`" + `gantt_team` + "`" + ` t ON t.id = tm.team_id
WHERE t.project_id = ?
`

func (q *Queries) GetProjectTeammembers(ctx context.Context, projectID int64) ([]GanttTeammember, error) {
	rows, err := q.db.QueryContext(ctx, getProjectTeammembers, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttTeammember
	for rows.Next() {
		var i GanttTeammember
		if err := rows.Scan(
			&i.ID,
			&i.RoleID,
			&i.TeamID,
			&i.UserID,
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

const getProjectTeams = `-- name: GetProjectTeams :many
SELECT id, name, project_id FROM ` + "`" + `gantt_team` + "`" + `
WHERE project_id = ?
`

func (q *Queries) GetProjectTeams(ctx context.Context, projectID int64) ([]GanttTeam, error) {
	rows, err := q.db.QueryContext(ctx, getProjectTeams, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttTeam
	for rows.Next() {
		var i GanttTeam
		if err := rows.Scan(&i.ID, &i.Name, &i.ProjectID); err != nil {
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

const getRole = `-- name: GetRole :one
SELECT id, name, project_id FROM ` + "`" + `gantt_role` + "`" + `
WHERE id = ?
`

func (q *Queries) GetRole(ctx context.Context, id int64) (GanttRole, error) {
	row := q.db.QueryRowContext(ctx, getRole, id)
	var i GanttRole
	err := row.Scan(&i.ID, &i.Name, &i.ProjectID)
	return i, err
}

const getState = `-- name: GetState :one
SELECT id, name, project_id FROM ` + "`" + `gantt_state` + "`" + `
WHERE id = ?
`

func (q *Queries) GetState(ctx context.Context, id int64) (GanttState, error) {
	row := q.db.QueryRowContext(ctx, getState, id)
	var i GanttState
	err := row.Scan(&i.ID, &i.Name, &i.ProjectID)
	return i, err
}

const getTask = `-- name: GetTask :one
SELECT id, name, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, description, project_id FROM ` + "`" + `gantt_task` + "`" + `
WHERE id = ?
`

func (q *Queries) GetTask(ctx context.Context, id int64) (GanttTask, error) {
	row := q.db.QueryRowContext(ctx, getTask, id)
	var i GanttTask
	err := row.Scan(
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
	)
	return i, err
}

const getTaskActivities = `-- name: GetTaskActivities :many
SELECT id, name, description, planned_start_date, planned_end_date, planned_budget, actual_start_date, actual_end_date, actual_budget, dependency_id, state_id, task_id FROM ` + "`" + `gantt_activity` + "`" + `
WHERE task_id = ?
`

func (q *Queries) GetTaskActivities(ctx context.Context, taskID int64) ([]GanttActivity, error) {
	rows, err := q.db.QueryContext(ctx, getTaskActivities, taskID)
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

const getTeam = `-- name: GetTeam :one
SELECT id, name, project_id FROM ` + "`" + `gantt_team` + "`" + `
WHERE id = ?
`

func (q *Queries) GetTeam(ctx context.Context, id int64) (GanttTeam, error) {
	row := q.db.QueryRowContext(ctx, getTeam, id)
	var i GanttTeam
	err := row.Scan(&i.ID, &i.Name, &i.ProjectID)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, first_name, last_name, avatar FROM ` + "`" + `user_user` + "`" + `
WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int32) (UserUser, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i UserUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Avatar,
	)
	return i, err
}

const listProjects = `-- name: ListProjects :many
SELECT DISTINCT p.id, p.name, p.planned_start_date, p.planned_end_date, p.actual_start_date, p.actual_end_date, p.description, p.project_manager_id FROM ` + "`" + `gantt_project` + "`" + ` p
LEFT JOIN ` + "`" + `gantt_team` + "`" + ` t ON p.id = t.project_id
LEFT JOIN ` + "`" + `gantt_teammember` + "`" + ` tm ON t.id = tm.team_id
WHERE p.project_manager_id = ? OR tm.user_id = ?
`

type ListProjectsParams struct {
	ProjectManagerID int32
	UserID           int32
}

func (q *Queries) ListProjects(ctx context.Context, arg ListProjectsParams) ([]GanttProject, error) {
	rows, err := q.db.QueryContext(ctx, listProjects, arg.ProjectManagerID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GanttProject
	for rows.Next() {
		var i GanttProject
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PlannedStartDate,
			&i.PlannedEndDate,
			&i.ActualStartDate,
			&i.ActualEndDate,
			&i.Description,
			&i.ProjectManagerID,
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
