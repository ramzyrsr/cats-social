package controllers

import (
	"cats-social/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RejectCat(c *gin.Context) {
	var requestBody struct {
		MatchID string `json:"matchId"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Invalid request body",
		})
		return
	}

	matchId, err := strconv.Atoi(requestBody.MatchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Match ID is not valid",
		})
		return
	}

	result, err := database.DB.Exec("UPDATE matched_cat SET status = 'Rejected' WHERE id = $1", matchId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Failed to reject data",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "Cat connection not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Succeed",
		"message": "Success reject cat",
	})
}
