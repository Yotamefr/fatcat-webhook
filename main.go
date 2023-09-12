package main

import (
	"fatcat_webhook/m/v2/routes"
	"fatcat_webhook/m/v2/utils"
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("env", ".env", "The path for the environment file")
	flag.Parse()
	godotenv.Load(*envFile)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/grafana", routes.GrafanaHandler)

	router.Run(utils.Getenv("FATCAT_WEBHOOK_HOST", "0.0.0.0") + ":" + utils.Getenv("FATCAT_WEBHOOK_PORT", "8080"))
}
