package main

import (
	"cats-social/router"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// Setup routes
	r := router.SetupRouter()

	// Run the server
	r.Run(":8080")
}
