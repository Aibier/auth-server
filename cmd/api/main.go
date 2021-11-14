package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"JWT-auth/internal/middleware"
	"JWT-auth/internal/routes"
)

func main() {

	app := gin.New()

	// Logging to a file.
	f, err := os.Create("log/" + time.Now().String())
	if err != nil {
		log.WithError(err).Printf("failed to create log file %s", err)
	}
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)
	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.Use(middleware.CORS())
	app.NoRoute(middleware.NoRouteHandler())

	app.POST("/login", routes.Login)
	port := os.Getenv("PORT")
	// Elastic Beanstalk forwards requests to port 5000
	if port == "" {
		port = "5000"
	}
	err = app.Run(":" + port)
	if err != nil {
		log.WithError(err).Error("failed to start the server")
	}
}
