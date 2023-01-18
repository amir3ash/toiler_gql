package dataloaders

//go:generate go run github.com/vektah/dataloaden ActivityLoader int64 []toiler-graphql/database.GanttActivity
//go:generate go run github.com/vektah/dataloaden StateLoader int64 *toiler-graphql/database.GanttState
//go:generate go run github.com/vektah/dataloaden UserLoader int32 *toiler-graphql/graph/model.UserUser
//go:generate go run github.com/vektah/dataloaden AssignedSliceLoader int64 []toiler-graphql/database.GanttAssigned
//go:generate go run github.com/vektah/dataloaden TaskSliceLoader int64 []toiler-graphql/database.GanttTask

import (
	"context"
	"time"
	"toiler-graphql/database"
	"toiler-graphql/graph/model"
)

type contextKey string

const key = contextKey("con")

type Loaders struct {
	TasksByProjectID      *TaskSliceLoader
	ActivitiesByTaskID    *ActivityLoader
	StateByID             *StateLoader
	UserByID              *UserLoader
	AssignedsByActivityID *AssignedSliceLoader
}

func newLoaders(ctx context.Context, repo database.Repository) *Loaders {
	return &Loaders{
		TasksByProjectID:      newTasksByProjectID(ctx, repo),
		ActivitiesByTaskID:    newActivityByTaskID(ctx, repo),
		StateByID:             newStateByID(ctx, repo),
		UserByID:              newUserByID(ctx, repo),
		AssignedsByActivityID: newAssignedsByActivityID(ctx, repo),
	}
}

type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newActivityByTaskID(ctx context.Context, repo database.Repository) *ActivityLoader {
	return NewActivityLoader(ActivityLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(taskIDs []int64) ([][]database.GanttActivity, []error) {
			// db query
			res, err := repo.ListActivitiesByTaskIDs(ctx, taskIDs)
			if err != nil {
				return nil, []error{err}
			}

			result := findSliceOfSlice(taskIDs, res, func(t database.GanttActivity) int64 { return t.TaskID })

			return result, nil
		},
	})
}

func newTasksByProjectID(ctx context.Context, repo database.Repository) *TaskSliceLoader {
	return NewTaskSliceLoader(TaskSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(projectIDs []int64) ([][]database.GanttTask, []error) {
			// db query
			res, err := repo.ListTasksByProjectIDs(ctx, projectIDs)
			if err != nil {
				return nil, []error{err}
			}

			result := findSliceOfSlice(projectIDs, res, func(t database.GanttTask) int64 { return t.ProjectID })

			return result, nil
		},
	})
}

func newStateByID(ctx context.Context, repo database.Repository) *StateLoader {
	return NewStateLoader(StateLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(stateIDs []int64) ([]*database.GanttState, []error) {
			// db query
			res, err := repo.ListStateByIDS(ctx, stateIDs)
			if err != nil {
				return nil, []error{err}
			}

			result := findSlice(stateIDs, res, func(t database.GanttState) int64 { return t.ID })

			return result, nil
		},
	})
}

func newUserByID(ctx context.Context, repo database.Repository) *UserLoader {
	return NewUserLoader(UserLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(userIDs []int32) ([]*model.UserUser, []error) {
			// db query
			res, err := repo.ListUsersIDS(ctx, userIDs)
			if err != nil {
				return nil, []error{err}
			}

			users := model.NormalizeUsersAvatar(res)
			result := findSlice(userIDs, users, func(t model.UserUser) int32 { return t.ID })

			return result, nil
		},
	})
}

func newAssignedsByActivityID(ctx context.Context, repo database.Repository) *AssignedSliceLoader {
	return NewAssignedSliceLoader(AssignedSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(activityIDs []int64) ([][]database.GanttAssigned, []error) {
			// db query
			res, err := repo.ListAssignedsByActivityIDs(ctx, activityIDs)
			if err != nil {
				return nil, []error{err}
			}

			result := findSliceOfSlice(activityIDs, res, func(t database.GanttAssigned) int64 { return t.ActivityID })

			return result, nil
		},
	})
}

/*
given ids and objects
returns list of ids to
*/
func findSliceOfSlice[K int32 | int64, T any](ids []K, objects []T, key func(T) K) [][]T {
	// map
	groupByParentID := make(map[K][]T, len(ids))
	for _, r := range objects {
		parentId := key(r)
		groupByParentID[parentId] = append(groupByParentID[parentId], r)
	}
	// order
	result := make([][]T, len(ids))
	for i, id := range ids {
		result[i] = groupByParentID[id]
	}
	return result
}

func findSlice[K int32 | int64, T any](ids []K, objects []T, key func(T) K) []*T {
	// map
	groupByParentID := make(map[K]T, len(ids))
	for _, r := range objects {
		parentId := key(r)
		groupByParentID[parentId] = r
	}
	// order
	result := make([]*T, len(ids))
	for i, id := range ids {
		if a, ok := groupByParentID[id]; ok {
			result[i] = &a
		} else {
			result[i] = nil
		}
	}
	return result
}
