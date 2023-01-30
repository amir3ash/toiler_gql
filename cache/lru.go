package cache

import (
	"toiler-graphql/database"
	"toiler-graphql/graph/model"

	lru "github.com/hashicorp/golang-lru"
)

type Cache interface {
	Project(id int64) (*database.GanttProject, bool)

	Task(id int64) (*database.GanttTask, bool)

	Activity(id int64) (*database.GanttActivity, bool)

	Assigned(id int64) (*database.GanttAssigned, bool)

	State(id int64) (*database.GanttState, bool)

	User(id int32) (*model.UserUser, bool)

	ProjectTasks(projectId int64) ([]database.GanttTask, bool)

	TaskActivities(taskId int64) ([]database.GanttActivity, bool)

	ProjectStates(projectId int64) ([]database.GanttState, bool)

	ActivityAssigneds(activityId int64) ([]database.GanttAssigned, bool)

	SetProject(o *database.GanttProject)

	SetTask(o *database.GanttTask)

	SetActivity(o *database.GanttActivity)

	SetAssigned(o *database.GanttAssigned)

	SetState(o *database.GanttState)

	SetUser(o *model.UserUser)

	SetProjectTasks([]database.GanttTask)

	SetTaskActivities([]database.GanttActivity)

	SetProjectStates([]database.GanttState)

	SetActivityAssigneds([]database.GanttAssigned)

	RemoveProject(id int64)

	RemoveTask(id int64)

	RemoveActivity(id int64)

	RemoveAssigned(id int64)

	RemoveState(id int64)

	RemoveUser(id int32)

	RemoveProjectTasks(projectId int64)

	RemoveTaskActivities(taskId int64)

	RemoveProjectStates(projectId int64)

	RemoveActivityAssigneds(activityId int64)
}

func NewLRU(cacheSize int) (*ganttLRU, error) {
	cache, err := lru.New2Q(cacheSize)
	if err != nil {
		return nil, err
	}

	return &ganttLRU{
		objects: cache,
	}, nil
}

type objectType int8

type objectKey struct {
	Type objectType
	Id   int64
}

const (
	defaultType objectType = iota
	projectType
	taskType
	activityType
	userType
	stateType
	assignedType
	taskListType
	activityListType
	stateListType
	assignedListType
)

type ganttLRU struct {
	objects *lru.TwoQueueCache
}

func (l *ganttLRU) Project(id int64) (*database.GanttProject, bool) {
	obj, ok := l.objects.Get(objectKey{projectType, id})
	if !ok {
		return nil, false
	}

	project, ok := obj.(*database.GanttProject)
	return project, ok
}

func (l *ganttLRU) Task(id int64) (*database.GanttTask, bool) {
	obj, ok := l.objects.Get(objectKey{taskType, id})
	if !ok {
		return nil, false
	}
	task, ok := obj.(*database.GanttTask)
	return task, ok
}

func (l *ganttLRU) Activity(id int64) (*database.GanttActivity, bool) {
	obj, ok := l.objects.Get(objectKey{activityType, id})
	if !ok {
		return nil, false
	}
	activity, ok := obj.(*database.GanttActivity)
	return activity, ok
}

// func (l *GanttLRU) Assigned(ctx context.Context, id int64) (*database.GanttAssigned, bool)

func (l *ganttLRU) State(id int64) (*database.GanttState, bool) {
	obj, ok := l.objects.Get(objectKey{stateType, id})
	if !ok {
		return nil, false
	}
	state, ok := obj.(*database.GanttState)
	return state, ok
}

func (l *ganttLRU) User(id int32) (*model.UserUser, bool) {
	obj, ok := l.objects.Get(objectKey{userType, int64(id)})
	if !ok {
		return nil, false
	}
	user, ok := obj.(*model.UserUser)
	return user, ok
}

func (l *ganttLRU) Assigned(id int64) (*database.GanttAssigned, bool) {
	obj, ok := l.objects.Get(objectKey{assignedType, id})
	if !ok {
		return nil, false
	}
	assigned, ok := obj.(*database.GanttAssigned)
	return assigned, ok
}

func (l *ganttLRU) ProjectTasks(projectId int64) ([]database.GanttTask, bool) {
	objects, ok := l.objects.Get(objectKey{taskListType, projectId})
	if !ok {
		return nil, false
	}

	tasks, ok := objects.([]database.GanttTask)
	return tasks, ok
}

func (l *ganttLRU) TaskActivities(taskId int64) ([]database.GanttActivity, bool) {
	objects, ok := l.objects.Get(objectKey{activityListType, taskId})
	if !ok {
		return nil, false
	}

	activities, ok := objects.([]database.GanttActivity)
	return activities, ok
}

func (l *ganttLRU) ProjectStates(projectId int64) ([]database.GanttState, bool) {
	objects, ok := l.objects.Get(objectKey{stateListType, projectId})
	if !ok {
		return nil, false
	}

	states, ok := objects.([]database.GanttState)
	return states, ok
}

func (l *ganttLRU) ActivityAssigneds(activityId int64) ([]database.GanttAssigned, bool) {
	objects, ok := l.objects.Get(objectKey{assignedListType, activityId})
	if !ok {
		return nil, false
	}

	assigneds, ok := objects.([]database.GanttAssigned)
	return assigneds, ok
}

func (l *ganttLRU) SetProject(o *database.GanttProject) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{projectType, o.ID}, o)
}

func (l *ganttLRU) SetTask(o *database.GanttTask) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{taskType, o.ID}, o)
}

func (l *ganttLRU) SetActivity(o *database.GanttActivity) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{activityType, o.ID}, o)
}

func (l *ganttLRU) SetAssigned(o *database.GanttAssigned) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{assignedType, o.ID}, o)
}

func (l *ganttLRU) SetState(o *database.GanttState) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{stateType, o.ID}, o)
}

func (l *ganttLRU) SetUser(o *model.UserUser) {
	if o == nil {
		return
	}

	l.objects.Add(objectKey{userType, int64(o.ID)}, o)
}

func (l *ganttLRU) SetProjectTasks(tasks []database.GanttTask) {
	if len(tasks) == 0 {
		return
	}

	l.objects.Add(objectKey{taskListType, tasks[0].ProjectID}, tasks)
}

func (l *ganttLRU) SetTaskActivities(activities []database.GanttActivity) {
	if len(activities) == 0 {
		return
	}

	l.objects.Add(objectKey{activityListType, activities[0].TaskID}, activities)
}

func (l *ganttLRU) SetProjectStates(states []database.GanttState) {
	if len(states) == 0 {
		return
	}

	l.objects.Add(objectKey{stateListType, states[0].ProjectID}, states)
}

func (l *ganttLRU) SetActivityAssigneds(assigneds []database.GanttAssigned) {
	if len(assigneds) == 0 {
		return
	}

	l.objects.Add(objectKey{assignedListType, assigneds[0].ActivityID}, assigneds)
}

func (l *ganttLRU) RemoveProject(id int64) {
	l.objects.Remove(objectKey{projectType, id})
}

func (l *ganttLRU) RemoveTask(id int64) {
	l.objects.Remove(objectKey{taskType, id})
}

func (l *ganttLRU) RemoveActivity(id int64) {
	l.objects.Remove(objectKey{activityType, id})
}

func (l *ganttLRU) RemoveAssigned(id int64) {
	l.objects.Remove(objectKey{assignedType, id})
}

func (l *ganttLRU) RemoveState(id int64) {
	l.objects.Remove(objectKey{stateType, id})
}

func (l *ganttLRU) RemoveUser(id int32) {
	l.objects.Remove(objectKey{userType, int64(id)})
}

func (l *ganttLRU) RemoveProjectTasks(projectId int64) {
	l.objects.Remove(objectKey{taskListType, projectId})
}

func (l *ganttLRU) RemoveTaskActivities(taskId int64) {
	l.objects.Remove(objectKey{activityListType, taskId})
}

func (l *ganttLRU) RemoveProjectStates(projectId int64) {
	l.objects.Remove(objectKey{stateListType, projectId})
}

func (l *ganttLRU) RemoveActivityAssigneds(activityId int64) {
	l.objects.Remove(objectKey{assignedListType, activityId})
}
