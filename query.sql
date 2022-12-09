-- name: GetActivity :one
SELECT * FROM `gantt_activity`
WHERE id = ?;

-- name: GetProjectActivities :many
SELECT ac.* FROM `gantt_activity` ac
JOIN `gantt_task` t ON ac.task_id = t.id
WHERE t.project_id = ?;

-- name: GetProjectAssignees :many
SELECT a.* FROM `gantt_assigned` a
JOIN `gantt_activity` ac ON a.activity_id = ac.id
JOIN `gantt_task` t ON ac.task_id = t.id
WHERE t.project_id = ?;

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

-- name: ListProjects :many
SELECT p.* FROM `gantt_project` p
JOIN `gantt_team` t ON p.id = t.project_id
JOIN `gantt_teammember` tm ON t.id = tm.team_id
WHERE p.project_manager_id = ? OR tm.user_id = ?;

-- name: GetStates :many
SELECT * FROM `gantt_state`
WHERE project_id = ?;

