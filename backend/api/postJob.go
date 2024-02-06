package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"

	// apiutils "k8s_leet_code_backend/api/utils"
	"k8s_leet_code_backend/helper"
	"k8s_leet_code_backend/model"
	"k8s_leet_code_backend/asynq_client"

	// "k8s_leet_code_backend/database"

	"github.com/google/uuid"
)

func PostK8sJob(c *gin.Context) {
	var codeRequestToExec CodeRequest
	if parsingErr := c.BindJSON(&codeRequestToExec); parsingErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parsingErr.Error()})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instanciationTime := time.Now().UTC()
	uniqueID := uuid.New().String()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		// Insert in DB
		entryDB := model.UserCodeRequest{
			UserID:          user.ID,
			InstanciationTS: instanciationTime,
			RequestUUID:     uniqueID,
			CodeContent:     codeRequestToExec.Content,
			WorkerStatus:    sql.NullString{},
			OutputResult:    sql.NullString{},
		}

		_, savedEntryDBErr := entryDB.Save()
		if savedEntryDBErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": savedEntryDBErr.Error()})
			return
		}
	}()

	go func() {
		defer wg.Done()
		// Push in Redis
		codeReqTS := ExtentedCodeRequest{
			UserID: user.ID,
			Language: codeRequestToExec.Language,
			Content: codeRequestToExec.Content,
			ReqTS:    instanciationTime.UnixMilli(),
			UniqueID: uniqueID,
		}
		payload, err := json.Marshal(codeReqTS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error when marshalling the payload", "err": err})
		}
		// if errRPush := redis.RedisClient.RPush(c, "queue:new-code-request", payload).Err(); errRPush != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Error when pushing the message to the queue.", "err": errRPush})
		// }
		asynqTask := asynq.NewTask("run_code", payload)
		asynqTaskInfo, asynqTaskErr := asynq_client.AsynqRedisClient.Enqueue(asynqTask)
		if asynqTaskErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error when pushing the message to the queue.", "err": asynqTaskErr})
		}
		log.Printf(" [*] Successfully enqueued task: %+v\nThe Payload is %v", asynqTaskInfo, payload)
	}()

	wg.Wait()

	// jobResponse := apiutils.CreateK8sJob(codeRequestToExec.Language, codeRequestToExec.Content)

	resultTime := time.Now().UTC().UnixMilli()

	timeDiff := (resultTime - instanciationTime.UnixMilli()) / 1000

	// c.String(http.StatusOK, fmt.Sprintf("UUID is %s\n\nUser is %s\nResult is %s\nCode executed is:\n\n%s\nTime between instanciation and result is %d seconds\n", uniqueID, user.Username, jobResponse, codeRequestToExec.Content, timeDiff))
	c.String(http.StatusOK, fmt.Sprintf("UUID is %s\n\nUser is %s\n\nCode executed is:\n\n%s\nTime between instanciation and result is %d seconds\n", uniqueID, user.Username, codeRequestToExec.Content, timeDiff))

}
