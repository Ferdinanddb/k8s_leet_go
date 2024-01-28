package api

import (
	"fmt"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	apiutils "k8s_leet_code/api/utils"
	"k8s_leet_code/helper"
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

	jobResponse := apiutils.CreateK8sJob(codeRequestToExec.Language, codeRequestToExec.Content)

	resultTime := time.Now().UTC().UnixMilli()

	timeDiff := (resultTime - instanciationTime) / 1000

	c.String(http.StatusOK, fmt.Sprintf("User is %s\nResult is %s\nCode executed is:\n\n%s\nTime between instanciation and result is %d seconds\n", user.Username, jobResponse, codeRequestToExec.Content, timeDiff))
}