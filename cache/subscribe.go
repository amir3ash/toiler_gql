package cache

import (
	"encoding/json"
	"errors"
	"log"
	"toiler-graphql/database"
	"toiler-graphql/graph/model"

	"github.com/go-redis/redis"
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
			return errors.New("can't convert Messages's obj to GanttProject")
		}
		db.lru.SetProject(&obj)

	case "task":
		obj, ok := m.Obj.(database.GanttTask)
		if !ok {
			return errors.New("can't convert Messages's obj to GanttTask")
		}
		db.lru.SetTask(&obj)

	case "activity":
		obj, ok := m.Obj.(database.GanttActivity)
		if !ok {
			return errors.New("can't convert Messages's obj to GanttActivity")
		}
		db.lru.SetActivity(&obj)

	case "state":
		obj, ok := m.Obj.(database.GanttState)
		if !ok {
			return errors.New("can't convert Messages's obj to GanttState")
		}
		db.lru.SetState(&obj)

	case "user":
		obj, ok := m.Obj.(model.UserUser)
		if !ok {
			return errors.New("can't convert Messages's obj to UserUser")
		}
		db.lru.SetUser(&obj)
	}
	return nil
}

func (db *redisDB) deleteCache(m *redisMessage) error {
	switch m.Type {
	case "project":
		db.lru.removeProject(m.Id)

	case "activity":
		db.lru.removeActivity(m.Id)

	case "task":
		db.lru.removeTask(m.Id)

	case "state":
		db.lru.removeState(m.Id)

	case "user":
		db.lru.removeUser(int32(m.Id))
	}

	return nil
}
