package main

import (
	"context"
	"log"
	"github.com/hibiken/asynq"
)

func main() {
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: "localhost:6379"},
        asynq.Config{Concurrency: 10},
    )

    mux := asynq.NewServeMux()
    mux.HandleFunc("email:welcome", sendWelcomeEmail)
    mux.HandleFunc("email:reminder", sendReminderEmail)

    if err := srv.Run(mux); err != nil {
        log.Fatal(err)
    }
}

func sendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
    var p EmailTaskPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err
    }
    log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
    return nil
}

func sendReminderEmail(ctx context.Context, t *asynq.Task) error {
    var p EmailTaskPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err
    }
    log.Printf(" [*] Send Reminder Email to User %d", p.UserID)
    return nil
}