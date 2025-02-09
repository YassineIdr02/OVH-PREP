package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/YassineIdr02/ovh-prep/E2E-Tests/routes"
	"github.com/YassineIdr02/ovh-prep/E2E-Tests/storage"
)

func main() {

	storage.InitDB()

	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	routes.SetupRoutes(router)

	log.Println("Server is running on :8080")
	router.Run(":8080")
}
