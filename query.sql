-- name: GetTask :one
SELECT * FROM `gantt_task`
WHERE id = ?;

-- name: GetProjectTasks :many
SELECT * FROM `gantt_task`
WHERE project_id = ?;

-- name: GetActivity :one
SELECT * FROM `gantt_activity`
WHERE id = ?;

-- name: GetTaskActivities :many
SELECT * FROM `gantt_activity`
WHERE task_id = ?;

-- name: GetProjectActivities :many
SELECT ac.* FROM `gantt_activity` ac
JOIN `gantt_task` t ON ac.task_id = t.id
WHERE t.project_id = ?;

-- name: GetProjectAssignees :many
SELECT a.* FROM `gantt_assigned` a
JOIN `gantt_activity` ac ON a.activity_id = ac.id
JOIN `gantt_task` t ON ac.task_id = t.id
WHERE t.project_id = ?;

-- name: GetAssigned :one
SELECT * FROM `gantt_assigned`
WHERE id = ?;

-- name: GetActivityAssignees :many
SELECT a.* FROM `gantt_assigned` a
WHERE a.activity_id = ?;

-- name: GetAssignedToUser :many
SELECT * FROM `gantt_assigned`
WHERE user_id = ?;

-- name: GetActivityComments :many
SELECT * FROM `gantt_comment`
WHERE activity_id = ?;

-- name: GetProjectEmployees :many
SELECT u.* FROM `user_user` u
JOIN `gantt_teammember` tm ON tm.user_id = u.id
JOIN `gantt_team` t ON t.id = tm.team_id
WHERE t.project_id = ?;

-- name: GetProjectTeammembers :many
SELECT tm.* FROM `gantt_teammember` tm
JOIN `gantt_team` t ON t.id = tm.team_id
WHERE t.project_id = ?;

-- name: GetProject :one
SELECT * FROM `gantt_project`
WHERE id = ?;

-- name: ListProjects :many
SELECT DISTINCT p.* FROM `gantt_project` p
LEFT JOIN `gantt_team` t ON p.id = t.project_id
LEFT JOIN `gantt_teammember` tm ON t.id = tm.team_id
WHERE p.project_manager_id = ? OR tm.user_id = ?;

-- name: GetProjectStates :many
SELECT * FROM `gantt_state`
WHERE project_id = ?;

-- name: GetState :one
SELECT * FROM `gantt_state`
WHERE id = ?;

-- name: GetRole :one
SELECT * FROM `gantt_role`
WHERE id = ?;

-- name: GetProjectRoles :many
SELECT * FROM `gantt_role`
WHERE project_id = ?;

-- name: GetTeam :one
SELECT * FROM `gantt_team`
WHERE id = ?;

-- name: GetProjectTeams :many
SELECT * FROM `gantt_team`
WHERE project_id = ?;


-- name: GetUser :one
SELECT * FROM `user_user`
WHERE id = ?;


