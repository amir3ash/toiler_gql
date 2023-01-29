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

	ProjectTasks(projectId int64) ([]database.GanttTask, bool)

	TaskActivities(taskId int64) ([]database.GanttActivity, bool)

	ProjectStates(projectId int64) ([]database.GanttState, bool)

	SetProject(o *database.GanttProject)

	SetTask(o *database.GanttTask)

	SetActivity(o *database.GanttActivity)

	SetAssigned(o *database.GanttAssigned)

	SetState(o *database.GanttState)

	SetUser(o *model.UserUser)

	SetProjectTasks([]database.GanttTask)

	SetTaskActivities([]database.GanttActivity)

	SetProjectStates([]database.GanttState)
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
	taskListType
	activityListType
	stateListType
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

func (l *ganttLRU) ProjectTasks(projectId int64) ([]database.GanttTask, bool) {
	ids, ok := l.objects.Get(objectKey{taskListType, projectId})
	if !ok {
		return nil, false
	}

	taskIds, ok := ids.([]int64)
	if !ok {
		fmt.Println("can't convert ids to list ")
		return nil, false
	}

	tasks := make([]database.GanttTask, len(taskIds))

	for i, id := range taskIds {
		task, found := l.Task(id)
		if !found {
			
		}

		tasks[i] = *task
	}
	return tasks, ok
}

func (l *ganttLRU) TaskActivities(taskId int64) ([]database.GanttActivity, bool) {
	return nil, false
}

func (l *ganttLRU) ProjectStates(projectId int64) ([]database.GanttState, bool) {
	return nil, false
}

func (l *ganttLRU) SetProject(o *database.GanttProject) {
	if o == nil {
		return
	}

	fmt.Println("set p", o.ID)
	l.objects.Add(objectKey{projectType, o.ID}, o)
}

func (l *ganttLRU) SetTask(o *database.GanttTask) {
	if o == nil {
		return
	}
	fmt.Println("set task", o.ID)

	l.objects.Add(objectKey{taskType, o.ID}, o)
}

func (l *ganttLRU) SetActivity(o *database.GanttActivity) {
	if o == nil {
		return
	}
	fmt.Println("set activity", o.ID)

	l.objects.Add(objectKey{activityType, o.ID}, o)
}

func (l *ganttLRU) SetAssigned(o *database.GanttAssigned) {
}

func (l *ganttLRU) SetState(o *database.GanttState) {
	if o == nil {
		return
	}
	fmt.Println("set state", o.ID)

	l.objects.Add(objectKey{stateType, o.ID}, o)
}

func (l *ganttLRU) SetUser(o *model.UserUser) {
	if o == nil {
		return
	}
	fmt.Println("set user", o.ID)

	l.objects.Add(objectKey{userType, int64(o.ID)}, o)
}

func (l *ganttLRU) SetProjectTasks(tasks []database.GanttTask) {
	if len(tasks) == 0 {
		return
	}

	ids := make([]int64, len(tasks))

	for i, v := range tasks {
		l.SetTask(&v)
		ids[i] = v.ID
	}

	l.objects.Add(objectKey{taskListType, tasks[0].ProjectID}, ids)
}

func (l *ganttLRU) SetTaskActivities(activities []database.GanttActivity) {
	if len(activities) == 0 {
		return
	}

	ids := make([]int64, len(activities))

	for i, v := range activities {
		l.SetActivity(&v)
		ids[i] = v.ID
	}

	l.objects.Add(objectKey{taskListType, activities[0].TaskID}, ids)
}

func (l *ganttLRU) SetProjectStates(states []database.GanttState) {
	if len(states) == 0 {
		return
	}

	ids := make([]int64, len(states))

	for i, v := range states {
		l.SetState(&v)
		ids[i] = v.ID
	}

	l.objects.Add(objectKey{taskListType, states[0].ProjectID}, ids)
}

func (l *ganttLRU) removeProject(id int64) {
	l.objects.Remove(objectKey{projectType, id})
}

func (l *ganttLRU) removeTask(id int64) {
	l.objects.Remove(objectKey{taskType, id})
}

func (l *ganttLRU) removeActivity(id int64) {
	l.objects.Remove(objectKey{activityType, id})
}

func (l *ganttLRU) removeAssigned(id int64) {
}

func (l *ganttLRU) removeState(id int64) {
	l.objects.Remove(objectKey{stateType, id})
}

func (l *ganttLRU) removeUser(id int32) {
	l.objects.Remove(objectKey{stateType, int64(id)})
}

func addTaskToTaskList(task *database.GanttTask) bool
