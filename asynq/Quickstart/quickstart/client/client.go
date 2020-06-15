package main

import (
	"github.com/hibiken/asynq"
	"log"
	"time"
)

var redis = asynq.RedisClientOpt{
	Addr: "localhost:6379",
	// Omit if no password is required
	Password: "mypassword",
	// Use a dedicated db number for asynq.
	// By default, Redis offers 16 databases (0..15)
	DB: 0,
}

// Task represents a task to be performed.
type Task struct {
	// Type indicates the type of a task to be performed.
	Type string
	// Payload holds data needed to perform the task.
}

func main() {

	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(r)

	// Create a task with typename and payload.
	t1 := asynq.NewTask("email:welcome", map[string]interface{}{"user_id": 42})

	t2 := asynq.NewTask("email:reminder", map[string]interface{}{"user_id": 42})

	// Process the task immediately.
	err := client.Enqueue(t1)
	if err != nil {
		log.Fatal(err)
	}
	// Process the task 24 hours later.
	err = client.EnqueueIn(24*time.Hour, t2)
	if err != nil {
		log.Fatal(err)
	}
}
