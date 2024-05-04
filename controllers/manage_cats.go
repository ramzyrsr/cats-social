package controllers

import (
	"cats-social/database"
	"cats-social/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

func CreateCat(c *gin.Context) {
	cookieValue, err := c.Cookie("Authorization")
	if err != nil {
		// Cookie not found or error occurred
		// Handle the error or return an error response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You don't have access",
		})
		return
	}

	token, err := jwt.Parse(cookieValue, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// This block of code is extracting the user ID from the JWT token. Here's a breakdown of what it
	// does:
	var id int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			id = int(sub)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
	}

	var requestBody struct {
		Name        string   `json:"name" binding:"required,min=1,max=30"`
		Race        string   `json:"race" binding:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
		Sex         string   `json:"sex" binding:"required,oneof=male female"`
		AgeInMonth  int      `json:"ageInMonth" binding:"required,min=1,max=120082"`
		Description string   `json:"description" binding:"required,min=1,max=200"`
		ImageUrls   []string `json:"imageUrls" binding:"required,min=1,dive,required,url"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": validationErrors[0].Translate(nil),
		})
		return
	}

	// Create the cat object
	cat := models.Cats{
		Name:        requestBody.Name,
		Race:        requestBody.Race,
		Sex:         requestBody.Sex,
		Age:         requestBody.AgeInMonth,
		Description: requestBody.Description,
		ImageUrl:    strings.Join(requestBody.ImageUrls, ","),
	}

	// Save the cat to the database or perform other actions
	query := `INSERT INTO cat (name, race, sex, age, description, image, owner_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id`

	// Execute the SQL query with the cat data
	var catID int
	err = database.DB.QueryRow(query, cat.Name, cat.Race, cat.Sex, cat.Age, cat.Description, cat.ImageUrl, id).Scan(&catID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "failed to insert cat into database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Success",
		"message": "Cat created successfully",
		"data": gin.H{
			"id":         catID,
			"created_at": time.Now(),
		},
	})
}

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
