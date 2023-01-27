package cache

import (
	"fmt"
	"toiler-graphql/database"
	"toiler-graphql/graph/model"

	lru "github.com/hashicorp/golang-lru"
)

type Cache interface {
	Project(id int64) (*database.GanttProject, bool)

	Task(id int64) (*database.GanttTask, bool)

	Activity(id int64) (*database.GanttActivity, bool)

	// Assigned(ctx context.Context, id int64) (database.GanttAssigned, bool)

	State(id int64) (*database.GanttState, bool)

	User(id int32) (*model.UserUser, bool)

	SetProject(o *database.GanttProject)

	SetTask(o *database.GanttTask)

	SetActivity(o *database.GanttActivity)

	SetAssigned(o *database.GanttAssigned)

	SetState(o *database.GanttState)

	SetUser(o *model.UserUser)
}

func NewLRU() (*GanttLRU, error) {
	projectCache, err := lru.New2Q(200)
	if err != nil {
		return nil, err
	}

	userCache, err := lru.New2Q(200)
	if err != nil {
		return nil, err
	}

	stateCache, err := lru.New2Q(200)
	if err != nil {
		return nil, err
	}

	return &GanttLRU{
		projects: projectCache,
		users:    userCache,
		states:   stateCache,
	}, nil
}

type GanttLRU struct {
	projects *lru.TwoQueueCache
	users    *lru.TwoQueueCache
	states   *lru.TwoQueueCache
}

func (l *GanttLRU) Project(id int64) (*database.GanttProject, bool) {
	obj, ok := l.projects.Get(id)
	if !ok {
		return nil, false
	}

	project, ok := obj.(*database.GanttProject)
	return project, ok
}

func (l *GanttLRU) Task(id int64) (*database.GanttTask, bool) {
	return nil, false
}

func (l *GanttLRU) Activity(id int64) (*database.GanttActivity, bool) {
	return nil, false
}

// func (l *GanttLRU) Assigned(ctx context.Context, id int64) (*database.GanttAssigned, bool)

func (l *GanttLRU) State(id int64) (*database.GanttState, bool) {
	obj, ok := l.states.Get(id)
	if !ok {
		return nil, false
	}
	state, ok := obj.(*database.GanttState)
	return state, ok
}

func (l *GanttLRU) User(id int32) (*model.UserUser, bool) {
	obj, ok := l.users.Get(id)
	if !ok {
		return nil, false
	}
	user, ok := obj.(*model.UserUser)
	return user, ok
}

func (l *GanttLRU) SetProject(o *database.GanttProject) {
	if o == nil {
		return
	}
	fmt.Println("set p", o.ID)
	l.projects.Add(o.ID, o)
}

func (l *GanttLRU) SetTask(o *database.GanttTask) {
}

func (l *GanttLRU) SetActivity(o *database.GanttActivity) {
}

func (l *GanttLRU) SetAssigned(o *database.GanttAssigned) {
}

func (l *GanttLRU) SetState(o *database.GanttState) {
	if o == nil {
		return
	}
	fmt.Println("set state", o.ID)

	l.states.Add(o.ID, o)
}

func (l *GanttLRU) SetUser(o *model.UserUser) {
	if o == nil {
		return
	}
	fmt.Println("set user", o.ID)

	l.users.Add(o.ID, o)
}
