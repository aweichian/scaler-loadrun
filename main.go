package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"scaler-loadrun/handler"
	"time"
)

func main() {
	// Logging to a file.
	env := os.Getenv("envID")
	var baseDir string
	if env == "dev" {
		baseDir = "./logs"
	} else {
		baseDir = fmt.Sprintf("/service/logs/app/%s/devops/webhook/", env)
	}
	ymd := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("http-%s.log", ymd)
	logFileName := filepath.Join(baseDir, fileName)

	f, _ := os.Create(logFileName)
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		// health check url
		v1.GET("/health", handler.Health)
		v1.GET("/task", handler.Task)
	}

	// custome http
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	//router.Run(":12345")

}
