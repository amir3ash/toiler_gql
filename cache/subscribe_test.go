package cache

import (
	"reflect"
	"testing"
)

type FakeCache struct {
	calledFuncs map[string]int64
}

func newFakeCache() *FakeCache {
	m := map[string]int64{}
	c := FakeCache{m}
	return &c
}

func (l *FakeCache) RemoveProject(id int64) {
	l.calledFuncs["RemoveProject"] = id
}

func (l *FakeCache) RemoveTask(id int64) {
	l.calledFuncs["RemoveTask"] = id

}

func (l *FakeCache) RemoveActivity(id int64) {
	l.calledFuncs["RemoveActivity"] = id
}

func (l *FakeCache) RemoveAssigned(id int64) {
	l.calledFuncs["RemoveAssigned"] = id
}

func (l *FakeCache) RemoveState(id int64) {
	l.calledFuncs["RemoveState"] = id
}

func (l *FakeCache) RemoveUser(id int32) {
	l.calledFuncs["RemoveUser"] = int64(id)
}

func (l *FakeCache) RemoveProjectTasks(projectId int64) {
	l.calledFuncs["RemoveProjectTasks"] = projectId
}

func (l *FakeCache) RemoveTaskActivities(taskId int64) {
	l.calledFuncs["RemoveTaskActivities"] = taskId
}

func (l *FakeCache) RemoveProjectStates(projectId int64) {
	l.calledFuncs["RemoveProjectStates"] = projectId
}

func (l *FakeCache) RemoveActivityAssigneds(activityId int64) {
	l.calledFuncs["RemoveActivityAssigneds"] = activityId
}

func TestDeleteCache(t *testing.T) {
	fakeCache := newFakeCache()
	db := redisDB{Client: nil, lru: fakeCache}

	testCases := []struct {
		msg                   redisMessage
		name                  string
		expectedFunctionCalls map[string]int64
		expectedError         bool
	}{
		{redisMessage{}, "without type", map[string]int64{}, true},
		{redisMessage{Type: "foo"}, "with wrong type", map[string]int64{}, true},

		{redisMessage{Type: "project", Id: 1}, "project", map[string]int64{
			"RemoveProject":       1,
			"RemoveProjectStates": 1,
			"RemoveProjectTasks":  1,
		}, false},
		{redisMessage{Type: "task", Id: 1, ParentId: 2}, "task", map[string]int64{
			"RemoveTask":           1,
			"RemoveTaskActivities": 1,
			"RemoveProjectTasks":   2,
		}, false},
		{redisMessage{Type: "activity", Id: 1, ParentId: 2}, "activity", map[string]int64{
			"RemoveActivity":          1,
			"RemoveTaskActivities":    2,
			"RemoveActivityAssigneds": 1,
		}, false},
		{redisMessage{Type: "state", Id: 1, ParentId: 2}, "state", map[string]int64{
			"RemoveState":         1,
			"RemoveProjectStates": 2,
		}, false},
		{redisMessage{Type: "assigned", Id: 1, ParentId: 2}, "assigned", map[string]int64{
			"RemoveAssigned":          1,
			"RemoveActivityAssigneds": 2,
		}, false},
		{redisMessage{Type: "user", Id: 1, ParentId: 2}, "assigned", map[string]int64{
			"RemoveUser": 1,
		}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.deleteCache(&tc.msg)
			if err == nil && tc.expectedError {
				t.Errorf("expected error but error not occord")

			} else if err != nil && !tc.expectedError {
				t.Errorf("err occord: %v", err)
			}

			if !reflect.DeepEqual(fakeCache.calledFuncs, tc.expectedFunctionCalls) {
				t.Errorf("called funcs not equal, got %v expected: %v", fakeCache.calledFuncs, tc.expectedFunctionCalls)
			}
			// purge the map
			for funcName := range fakeCache.calledFuncs {
				delete(fakeCache.calledFuncs, funcName)
			}
		})
	}
}
