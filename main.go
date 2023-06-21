package main

import (
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raojinlin/clogs/docker"
	"time"
)

type SSEWriter struct {
	ctx *gin.Context
}

func (s *SSEWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	s.ctx.SSEvent("message", p)
	s.ctx.Writer.Flush()
	return
}

func main() {
	var port = 8082
	flag.IntVar(&port, "port", port, "Listen port")
	flag.Parse()
	router := gin.Default()
	containerGroup := router.Group("/api/container")
	{
		containerGroup.GET("logs/:container", func(c *gin.Context) {
			container := c.Param("container")
			if container == "" {
				c.AbortWithStatus(400)
				return
			}

			tail := c.Query("tail")
			if tail == "" {
				tail = "10"
			}
			logFile := c.Query("logFile")
			if logFile == "" {
				logFile = "stdout"
			}

			logOptions := &types.ContainerLogsOptions{
				Tail:       tail,
				ShowStderr: c.Query("showStderr") != "",
				ShowStdout: c.Query("showStdout") != "",
				Follow:     c.Query("follow") != "",
			}
			logsOut, err := docker.Logs(container, logFile, logOptions)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}

			sseWriter := &SSEWriter{ctx: c}
			c.Header("content-type", "text/event-stream")
			c.Header("X-Accel-Buffering", "no")
			_, err = stdcopy.StdCopy(sseWriter, sseWriter, logsOut)
			if err != nil {
				logsOut.Close()
			}
		})
	}

	router.LoadHTMLGlob("template/*")
	router.GET("/logs/:container", func(c *gin.Context) {
		c.Header("content-type", "text/html")
		logFile := c.Query("logFile")
		if logFile == "" {
			logFile = "stdout"
		}

		tail := c.Query("tail")
		if tail == "" {
			tail = "1500"
		}

		c.HTML(200, "log.html", gin.H{
			"Container": c.Param("container"),
			"LogFile":   logFile,
			"Tail":      tail,
		})
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
