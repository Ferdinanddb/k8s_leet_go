package api

import (
	"fmt"
	"time"
	"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"

	apiutils "k8s_leet_code/api/utils"
	"k8s_leet_code/helper"
	"k8s_leet_code/redis"

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

	instanciationTime := time.Now().UTC().UnixMilli()

	uniqueID := uuid.New().String()


	codeReqTS := ExtentedCodeRequest{
		codeRequestToExec,
		instanciationTime,
		uniqueID,
	}
	payload, err := json.Marshal(codeReqTS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error when marshalling the Payload", "err": err})
	}
	if errRPush := redis.RedisClient.RPush(c, "queue:new-code-request", payload).Err(); errRPush != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error when pushing the message to the queue.", "err": errRPush})
	}


	jobResponse := apiutils.CreateK8sJob(codeRequestToExec.Language, codeRequestToExec.Content)

	resultTime := time.Now().UTC().UnixMilli()

	timeDiff := (resultTime - instanciationTime) / 1000

	c.String(http.StatusOK, fmt.Sprintf("UUID is %s\n\nUser is %s\nResult is %s\nCode executed is:\n\n%s\nTime between instanciation and result is %d seconds\n", codeReqTS.UniqueID, user.Username, jobResponse, codeRequestToExec.Content, timeDiff))
}