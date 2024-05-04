package controllers

import (
	"cats-social/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterUser(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" validate:"required,customEmail"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Invalid request body",
		})
		return
	}

	var existingEmail string
	err := database.DB.QueryRow(`SELECT email FROM "user" WHERE email = $1`, requestBody.Email).Scan(&existingEmail)

	if err == nil {
		// Email already exists, return an error
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  "Failed",
			"message": "Email already exists",
		})
		return
	}
	if err == sql.ErrNoRows {
		// Email does not exist, continue with the registration process
		register, _ := database.DB.Exec(`INSERT INTO "user"
		("name", email, "password")
		VALUES($1, $2, $3);`, requestBody.Name, requestBody.Email, requestBody.Password)

		rowsAffected, _ := register.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "Failed",
				"message": "User failed registered",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "Succeed",
			"message": "User successfully registered",
		})

		return
	} else {
		// Handle other errors
		fmt.Println("Error:", err)
		return
	}
}

func Login(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email" validate:"required,customEmail"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Invalid request body",
		})
		return
	}

	var id int
	query := `SELECT id FROM "user" WHERE email = $1 AND "password" = $2`
	err := database.DB.QueryRow(query, requestBody.Email, requestBody.Password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were found, return an error indicating user not found
			// return nil, fmt.Errorf("user not found")
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Failed",
				"message": "User not found",
			})
			return
		}
		// For other errors, return the error
		// return nil, err
	}

	expTime := time.Now().Add(10 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expTime.Unix(),
		"sub": id,
	})

	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenStr, 3600*24, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "Succeed",
		"message": "Login Success",
	})
}
