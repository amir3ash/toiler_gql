package cache

import (
	"reflect"
	"testing"
	"toiler-graphql/database"
	"toiler-graphql/graph/model"
)

func TestEmptyLru(t *testing.T) {
	lru, err := NewLRU(64)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	cache := lru
	// cache.
	proj, ok := cache.Project(1)
	if proj != nil || ok {
		t.Error("expected empty project")
	}

	task, ok := cache.Task(1)
	if task != nil || ok {
		t.Error("expected empty task")
	}

	act, ok := cache.Activity(1)
	if act != nil || ok {
		t.Error("expected empty activity")
	}

	assigned, ok := cache.Assigned(1)
	if assigned != nil || ok {
		t.Error("expected empty assigned")
	}

	state, ok := cache.State(1)
	if state != nil || ok {
		t.Error("expected empty state")
	}

	user, ok := cache.User(1)
	if user != nil || ok {
		t.Error("expected empty user")
	}

	projectTasks, ok := cache.ProjectTasks(1)
	if projectTasks != nil || ok {
		t.Error("expected empty projectTasks")
	}

	taskActivities, ok := cache.TaskActivities(1)
	if taskActivities != nil || ok {
		t.Error("expected empty taskActivities")
	}

	projectStates, ok := cache.ProjectStates(1)
	if projectStates != nil || ok {
		t.Error("expected empty projectStates")
	}

	activityAssigneds, ok := cache.ActivityAssigneds(1)
	if activityAssigneds != nil || ok {
		t.Error("expected empty activityAssigneds")
	}

}

func equalSlices[T any](o1 []T, o2 []T) bool {
	return reflect.DeepEqual(o1, o2)
}

func TestLru(t *testing.T) {
	lru, err := NewLRU(64)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	cache := Cache(lru)

	project1 := &database.GanttProject{ID: 1}
	project2 := &database.GanttProject{ID: 2, Name: "second"}
	task1 := &database.GanttTask{ID: 1, ProjectID: 1}
	activity2 := &database.GanttActivity{ID: 2}
	assigned1 := &database.GanttAssigned{ID: 1}
	state1 := &database.GanttState{ID: 1, ProjectID: 1}
	user1 := &model.UserUser{UserUser: database.UserUser{ID: 1}}
	project1Tasks := []database.GanttTask{*task1}
	task1Activities := []database.GanttActivity{*activity2}
	project1States := []database.GanttState{*state1}

	cache.SetProject(project1)
	cache.SetProject(&database.GanttProject{ID: 2})
	cache.SetProject(project2) // update it
	cache.SetTask(task1)
	cache.SetTask(&database.GanttTask{ID: 2})
	cache.SetActivity(&database.GanttActivity{ID: 1})
	cache.SetActivity(activity2)
	cache.SetActivity(&database.GanttActivity{ID: 3})
	cache.SetAssigned(assigned1)
	cache.SetAssigned(&database.GanttAssigned{ID: 2})
	cache.SetState(state1)
	cache.SetUser(user1)
	cache.SetProjectTasks(1, project1Tasks)
	cache.SetTaskActivities(1, task1Activities)
	cache.SetProjectStates(1, project1States)

	if o, _ := cache.Project(1); o != project1 {
		t.Errorf("projects not equal: got %v expected: %v", o, project1)
	}
	if o, _ := cache.Project(2); o != project2 {
		t.Errorf("projects not equal: got %v expected: %v", o, project2)
	}

	if o, _ := cache.Task(1); o != task1 {
		t.Errorf("tasks not equal: got %v expected: %v", o, task1)
	}

	if o, _ := cache.Activity(2); o != activity2 {
		t.Errorf("activies not equal: got %v expected: %v", o, activity2)
	}

	if o, _ := cache.Assigned(1); o != assigned1 {
		t.Errorf("assigneds not equal: got %v expected: %v", o, assigned1)
	}

	if o, _ := cache.State(1); o != state1 {
		t.Errorf("states not equal: got %v expected: %v", o, state1)
	}

	if o, _ := cache.User(1); o != user1 {
		t.Errorf("users not equal: got %v expected: %v", o, user1)
	}

	// removing single objects
	cache.RemoveProject(1)
	if _, ok := cache.Project(1); ok {
		t.Error("project1 not removed")
	}
	if _, ok := cache.Task(1); !ok {
		t.Error("task1 removed when removing the project")
	}
	if _, ok := cache.State(1); !ok {
		t.Error("state1 removed when removing the project")
	}
	if _, ok := cache.Project(2); !ok {
		t.Error("project2 removed")
	}

	cache.RemoveTask(1)
	if _, ok := cache.Task(1); ok {
		t.Error("task1 not removed")
	}
	if _, ok := cache.Activity(1); !ok {
		t.Error("actvity1 removed when removing the project")
	}
	if _, ok := cache.Task(2); !ok {
		t.Error("task2 removed")
	}

	cache.RemoveActivity(1)
	if _, ok := cache.Activity(1); ok {
		t.Error("activity1 not removed")
	}
	if _, ok := cache.Activity(2); !ok {
		t.Error("activity2 removed")
	}

	cache.RemoveAssigned(1)
	if _, ok := cache.Assigned(1); ok {
		t.Error("assigned1 not removed")
	}

	cache.RemoveState(1)
	if _, ok := cache.State(1); ok {
		t.Error("state1 not removed")
	}

	cache.RemoveUser(1)
	if _, ok := cache.User(1); ok {
		t.Error("user1 not removed")
	}

	// setting object slices
	if o, _ := cache.ProjectTasks(1); !equalSlices(o, project1Tasks) {
		t.Errorf("projectTasks not equal: got %v expected: %v", o, project1Tasks)
	}

	if o, _ := cache.TaskActivities(1); !equalSlices(o, task1Activities) {
		t.Errorf("taskActivities not equal: got %v expected: %v", o, task1Activities)
	}

	if o, _ := cache.ProjectStates(1); !equalSlices(o, project1States) {
		t.Errorf("projectStates not equal: got %v expected: %v", o, project1States)
	}

	// if o, _ := cache.ActivityAssigneds(1); !testEqualslice(o, a) {
	// 	t.Errorf("projectStates not equal: got %v expected: %v", o, project1States)
	// }

	cache.RemoveProjectTasks(1)
	if _, ok := cache.ProjectTasks(1); ok {
		t.Error("project1Tasks not removed")
	}

	cache.RemoveTaskActivities(1)
	if _, ok := cache.TaskActivities(1); ok {
		t.Error("task1Activities not removed")
	}

	cache.RemoveProjectStates(1)
	if _, ok := cache.ProjectStates(1); ok {
		t.Error("project1States not removed")
	}

	// cache.RemoveActivityAssigneds(1)
	// if _, ok := cache.State(1); ok {
	// 	t.Error("state1 not removed")
	// }
}
