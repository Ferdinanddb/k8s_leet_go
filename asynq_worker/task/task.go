package task


import (
    "context"
    // "fmt"
	"encoding/json"
	"log"
	"time"
   
    "github.com/hibiken/asynq"
	taskutils "k8s_leet_code_asynq_worker/task/utils"
)

// A list of task types.
const (
    TypeRunCode  = "run_code"
	// TypeRunCode  = "queue:new-code-request"
	
)

type RunCodePayload struct {
    UserID uint
	Language     string
    Content  string
	ReqTS int64
	UniqueID string
}

func HandleRunCodePythonTask(ctx context.Context, t *asynq.Task) error {
    var p RunCodePayload  
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err
    }

	instanciationTime := time.Now().UTC()

	jobResponse := taskutils.CreateK8sJob(p.Language, p.Content, p.UniqueID, p.UserID)

	resultTime := time.Now().UTC().UnixMilli()
	timeDiff := (resultTime - instanciationTime.UnixMilli()) / 1000
	log.Printf("Time taken to process the request : %d seconds", timeDiff)

	log.Printf(" [*] Received a task with uniqueID %s from %d to run a %s script which content is :\n\n%s\n\nResult is : %s ", p.UniqueID, p.UserID, p.Language, p.Content, jobResponse)
    return nil
}
