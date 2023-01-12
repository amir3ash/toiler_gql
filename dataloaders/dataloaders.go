package dataloaders

//go:generate go run github.com/vektah/dataloaden ActivityLoader int64 []toiler-graphql/database.GanttActivity
//go:generate go run github.com/vektah/dataloaden StateLoader int64 *toiler-graphql/database.GanttState
//go:generate go run github.com/vektah/dataloaden UserLoader int64 *toiler-graphql/database.UserUser
//go:generate go run github.com/vektah/dataloaden AssignedSliceLoader int64 []toiler-graphql/database.GanttAssigned

import (
	"context"
	"time"
	"toiler-graphql/database"
)

type contextKey string

const key = contextKey("con")

type Loaders struct {
	ActivitiesByTaskID    *ActivityLoader
	StateByActivityID     *StateLoader
	UserByAssignedID      *UserLoader
	AssignedsByActivityID *AssignedSliceLoader
}

func newLoaders(ctx context.Context, repo database.Repository) *Loaders {
	return &Loaders{
		ActivitiesByTaskID:    newActivityByTaskID(ctx, repo),
		StateByActivityID:     newStateByActivityID(ctx, repo),
		UserByAssignedID:      newUserByAssignedID(ctx, repo),
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

func newStateByActivityID(ctx context.Context, repo database.Repository) *StateLoader {
	return NewStateLoader(StateLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(activityIDs []int64) ([]*database.GanttState, []error) {
			// db query
			res, err := repo.ListStateByActivityIDS(ctx, activityIDs)
			if err != nil {
				return nil, []error{err}
			}

			statesAndParents := findSlice(activityIDs, res, func(t database.StateAndParent[int64]) int64 { return t.ParentId })

			result := make([]*database.GanttState, len(statesAndParents))
			for i, s := range statesAndParents {
				if s != nil {
					result[i] = &database.GanttState{ID: s.ID, Name: s.Name, ProjectID: s.ProjectID}
				}
			}
			return result, nil
		},
	})
}

func newUserByAssignedID(ctx context.Context, repo database.Repository) *UserLoader {
	return NewUserLoader(UserLoaderConfig{
		MaxBatch: 100,
		Wait:     400 * time.Microsecond,
		Fetch: func(assigneeIDs []int64) ([]*database.UserUser, []error) {
			// db query
			res, err := repo.ListUsersByAssignedIDS(ctx, assigneeIDs)
			if err != nil {
				return nil, []error{err}
			}

			usersAndParents := findSlice(assigneeIDs, res, func(t database.UserAndParent[int64]) int64 { return t.ParentId })

			result := make([]*database.UserUser, len(usersAndParents))
			for i, s := range usersAndParents {
				if s != nil {
					result[i] = &database.UserUser{
						ID:        s.ID,
						Username:  s.Username,
						FirstName: s.FirstName,
						LastName:  s.LastName,
						Avatar:    s.Avatar,
					}
				}
			}
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
