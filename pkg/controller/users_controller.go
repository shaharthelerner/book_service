package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/service/pkg/data/request"
)

func GetUserActivity(c *gin.Context) {
	var req request.UserActivityRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	actions, err := FetchUserActivity(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, actions)
}

// Remove
//func getIdQueryParam(c *gin.Context) string {
//	return c.Query("id")
//}

// Remove
//func getUsernameQueryParam(c *gin.Context) string {
//	return c.Query("username")
//}
