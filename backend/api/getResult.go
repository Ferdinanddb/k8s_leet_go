package api

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"k8s_leet_code_backend/helper"
)

func GetUserCodeReqHistory(context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user.CodeRequests})
}
