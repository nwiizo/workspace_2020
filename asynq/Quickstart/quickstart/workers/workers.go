package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
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

func handler(ctx context.Context, t *asynq.Task) error {
	switch t.Type {
	case "email:welcome":
		id, err := t.Payload.GetInt("user_id")
		if err != nil {
			return err
		}
		fmt.Printf("Send Welcome Email to User %d\n", id)

	case "email:reminder":
		id, err := t.Payload.GetInt("user_id")
		if err != nil {
			return err
		}
		fmt.Printf("Send Reminder Email to User %d\n", id)

	default:
		return fmt.Errorf("unexpected task type: %s", t.Type)
	}
	return nil
}
func main() {
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	srv := asynq.NewServer(r, asynq.Config{
		Concurrency: 10,
	})
	// Use asynq.HandlerFunc adapter for a handler function
	if err := srv.Run(asynq.HandlerFunc(handler)); err != nil {
		log.Fatal(err)
	}
}
