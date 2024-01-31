package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"sync"

	"net/http"

	"github.com/gin-gonic/gin"

	apiutils "k8s_leet_code/api/utils"
	"k8s_leet_code/helper"
	"k8s_leet_code/model"
	"k8s_leet_code/redis"
	// "k8s_leet_code/database"

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

	go func () {
		defer wg.Done()
		// Insert in DB
		entryDB := model.UserCodeRequest{
			UserID: user.ID,
			InstanciationTS: instanciationTime,
			RequestUUID: uniqueID,
			CodeContent: codeRequestToExec.Content,
			WorkerStatus: sql.NullString{},
			OutputResult: sql.NullString{},
		}

		_, savedEntryDBErr := entryDB.Save()
		if savedEntryDBErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": savedEntryDBErr.Error()})
			return
		}
	}()

	go func () {
		defer wg.Done()
		// Push in Redis
		codeReqTS := ExtentedCodeRequest{
			CodeReq: codeRequestToExec,
			ReqTS: instanciationTime.UnixMilli(),
			UniqueID: uniqueID,
		}
		payload, err := json.Marshal(codeReqTS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error when marshalling the payload", "err": err})
		}
		if errRPush := redis.RedisClient.RPush(c, "queue:new-code-request", payload).Err(); errRPush != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error when pushing the message to the queue.", "err": errRPush})
		}
	}()

	wg.Wait()



	jobResponse := apiutils.CreateK8sJob(codeRequestToExec.Language, codeRequestToExec.Content)

	resultTime := time.Now().UTC().UnixMilli()

	timeDiff := (resultTime - instanciationTime.UnixMilli()) / 1000

	c.String(http.StatusOK, fmt.Sprintf("UUID is %s\n\nUser is %s\nResult is %s\nCode executed is:\n\n%s\nTime between instanciation and result is %d seconds\n", uniqueID, user.Username, jobResponse, codeRequestToExec.Content, timeDiff))
}