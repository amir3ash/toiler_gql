// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package database

import (
	"database/sql"
	"time"
)

type AuthGroup struct {
	ID   int32
	Name string
}

type AuthGroupPermission struct {
	ID           int32
	GroupID      int32
	PermissionID int32
}

type AuthPermission struct {
	ID            int32
	Name          string
	ContentTypeID int32
	Codename      string
}

type ChangesDetail struct {
	ID            int32
	ObjectID      string
	Fields        interface{}
	ContentTypeID int32
	DateCreated   time.Time
	Username      string
}

type DjangoAdminLog struct {
	ID            int32
	ActionTime    time.Time
	ObjectID      sql.NullString
	ObjectRepr    string
	ActionFlag    int32
	ChangeMessage string
	ContentTypeID sql.NullInt32
	UserID        int32
}

type DjangoContentType struct {
	ID       int32
	AppLabel string
	Model    string
}

type DjangoMigration struct {
	ID      int32
	App     string
	Name    string
	Applied time.Time
}

type DjangoSession struct {
	SessionKey  string
	SessionData string
	ExpireDate  time.Time
}

type GanttActivity struct {
	ID               int64
	Name             string
	Description      string
	PlannedStartDate time.Time
	PlannedEndDate   time.Time
	PlannedBudget    sql.NullString
	ActualStartDate  sql.NullTime
	ActualEndDate    sql.NullTime
	ActualBudget     sql.NullString
	DependencyID     sql.NullInt64
	StateID          sql.NullInt64
	TaskID           int64
}

type GanttAssigned struct {
	ID         int64
	ActivityID int64
	UserID     int32
}

type GanttComment struct {
	ID         int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Text       string
	ActivityID int64
	AuthorID   int32
}

type GanttProject struct {
	ID               int64
	Name             string
	PlannedStartDate time.Time
	PlannedEndDate   time.Time
	ActualStartDate  sql.NullTime
	ActualEndDate    sql.NullTime
	Description      string
	ProjectManagerID int32
}

type GanttRole struct {
	ID        int64
	Name      string
	ProjectID int64
}

type GanttState struct {
	ID        int64
	Name      string
	ProjectID int64
}

type GanttTask struct {
	ID               int64
	Name             string
	PlannedStartDate time.Time
	PlannedEndDate   time.Time
	PlannedBudget    string
	ActualStartDate  sql.NullTime
	ActualEndDate    sql.NullTime
	ActualBudget     string
	Description      string
	ProjectID        int64
}

type GanttTeam struct {
	ID        int64
	Name      string
	ProjectID int64
}

type GanttTeammember struct {
	ID     int64
	RoleID int64
	TeamID int64
	UserID int32
}

type ReversionRevision struct {
	ID          int32
	DateCreated time.Time
	Comment     string
	UserID      sql.NullInt32
}

type ReversionVersion struct {
	ID             int32
	ObjectID       string
	Format         string
	SerializedData string
	ObjectRepr     string
	ContentTypeID  int32
	RevisionID     int32
	Db             string
}

type UserTempuser struct {
	ID         int32
	Password   string
	LastLogin  sql.NullTime
	LastName   string
	IsActive   bool
	DateJoined time.Time
	Username   string
	Email      string
	FirstName  string
}

type UserUser struct {
	ID        int32
	Username  string
	FirstName string
	LastName  string
	Avatar    sql.NullString
}

type UserUserGroup struct {
	ID      int32
	UserID  int32
	GroupID int32
}

type UserUserUserPermission struct {
	ID           int32
	UserID       int32
	PermissionID int32
}
