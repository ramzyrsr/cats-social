package controllers

import (
	"cats-social/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteCats(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "ID is not valid",
		})
		return
	}

	result, err := database.DB.Exec("DELETE FROM cat WHERE id = $1", catID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Failed to delete data",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "Cat not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Succeed",
		"message": "Success delete cat",
	})
}
