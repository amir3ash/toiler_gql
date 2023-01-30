package cache

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type redisMessage struct {
	Event    string
	Type     string
	Id       int64
	ParentId int64 `json:"parent"`
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
			log.Default().Printf("err while recieving message: %v\n", err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
			log.Default().Printf("err while unmarshalling message: %v\n", err)
		}

		switch m.Event {
		case "added", "updated", "deleted":
			err = db.deleteCache(&m)

		default:
			log.Default().Println("event not detected:", m.Event)
		}

		if err != nil {
			log.Default().Printf("error in consuming event: %v\n", err)
		}
	}
}

func (db *redisDB) deleteCache(m *redisMessage) error {
	switch m.Type {
	case "project":
		db.lru.RemoveProject(m.Id)
		db.lru.RemoveProjectStates(m.Id)
		db.lru.RemoveProjectTasks(m.Id)

	case "activity":
		db.lru.RemoveActivity(m.Id)
		db.lru.RemoveTaskActivities(m.ParentId)

	case "task":
		db.lru.RemoveTask(m.Id)
		db.lru.RemoveTaskActivities(m.Id)
		db.lru.RemoveProjectTasks(m.ParentId)

	case "state":
		db.lru.RemoveState(m.Id)
		db.lru.RemoveProjectStates(m.ParentId)

	case "user":
		db.lru.RemoveUser(int32(m.Id))

	case "assigned":
		db.lru.RemoveAssigned(m.Id)
		db.lru.RemoveActivityAssigneds(m.ParentId)

	default:
		return fmt.Errorf("event has unknown type: %s", m.Type)
	}

	return nil
}
