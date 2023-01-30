package cache

import (
	"encoding/json"
	"errors"
	"log"
	"toiler-graphql/database"
	"toiler-graphql/graph/model"

	"github.com/go-redis/redis"
)

var (
	errCantConvertToProject  = errors.New("can't convert Messages's obj to GanttProject")
	errCantconvertToTask     = errors.New("can't convert Messages's obj to GanttTask")
	errCantConvertToActivity = errors.New("can't convert Messages's obj to GanttActivity")
	errCantConvertToState    = errors.New("can't convert Messages's obj to GanttState")
	errCantConvertToUser     = errors.New("can't convert Messages's obj to UserUser")
)

type redisMessage struct {
	Event string
	Type  string
	Id    int64
	Obj   interface{}
}

type redisDB struct {
	Client *redis.Client
	lru    *ganttLRU
}

func NewRedisDB(redisAddr string, password string, database int, lru *ganttLRU) (*redisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       database,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &redisDB{client, lru}, nil
}

func (db *redisDB) ConsumeEvents() {
	go db.consume()
}

func (db *redisDB) consume() {
	subscriber := db.Client.Subscribe("changes")

	m := redisMessage{}

	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
			panic(err)
		}

		switch m.Event {
		case "added", "updated":
			err = db.updateCache(&m)

		case "deleted":
			err = db.deleteCache(&m)

		default:
			log.Default().Println("event not detected:", m.Event)
		}

		if err != nil {
			log.Default().Printf("error in consuming event: %v\n", err)
		}
	}
}

func (db *redisDB) updateCache(m *redisMessage) error {
	switch m.Type {
	case "project":
		obj, ok := m.Obj.(database.GanttProject)
		if !ok {
			return errCantConvertToProject
		}
		db.lru.SetProject(&obj)

	case "task":
		obj, ok := m.Obj.(database.GanttTask)
		if !ok {
			return errCantconvertToTask
		}
		db.lru.SetTask(&obj)

	case "activity":
		obj, ok := m.Obj.(database.GanttActivity)
		if !ok {
			return errCantConvertToActivity
		}
		db.lru.SetActivity(&obj)

	case "state":
		obj, ok := m.Obj.(database.GanttState)
		if !ok {
			return errCantConvertToState
		}
		db.lru.SetState(&obj)

	case "user":
		obj, ok := m.Obj.(model.UserUser)
		if !ok {
			return errCantConvertToUser
		}
		db.lru.SetUser(&obj)
	}
	return nil
}

func (db *redisDB) deleteCache(m *redisMessage) error {
	switch m.Type {
	case "project":
		db.lru.RemoveProject(m.Id)
		db.lru.RemoveProjectStates(m.Id)
		db.lru.RemoveProjectTasks(m.Id)

	case "activity":
		db.lru.RemoveActivity(m.Id)

		obj, ok := m.Obj.(database.GanttActivity)
		if !ok {
			return errCantConvertToActivity
		}

		activities, found := db.lru.TaskActivities(obj.TaskID)
		if !found {
			return nil
		}
		new := listObjWithoutIndex(activities, func(i int) bool {return activities[i].ID == m.Id})
		db.lru.SetTaskActivities(new)

	case "task":
		db.lru.RemoveTask(m.Id)
		db.lru.RemoveTaskActivities(m.Id)

		obj, ok := m.Obj.(database.GanttTask)
		if !ok {
			return errCantconvertToTask
		}

		tasks, found := db.lru.ProjectTasks(obj.ProjectID)
		if !found {
			return nil
		}
		new := listObjWithoutIndex(tasks, func(i int) bool {return tasks[i].ID == m.Id})
		db.lru.SetProjectTasks(new)

	case "state":
		db.lru.RemoveState(m.Id)

		obj, ok := m.Obj.(database.GanttState)
		if !ok {
			return errCantConvertToState
		}

		states, found := db.lru.ProjectStates(obj.ProjectID)
		if !found {
			return nil
		}
		new := listObjWithoutIndex(states, func(i int) bool {return states[i].ID == m.Id})
		db.lru.SetProjectStates(new)

	case "user":
		db.lru.RemoveUser(int32(m.Id))
	}

	return nil
}

/* if index found, returns new slice without the element at the index,
 else retruns objects
*/
func listObjWithoutIndex[T any](objects []T, f func(index int) bool) []T {
	index := -1
	for i := range objects {
		if f(i){
			index = i
			break
		}
	}
	
	if index == -1 {
		return objects
	}

	return append(objects[:index], objects[index+1:]...)
}
