package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/raojinlin/clogs/docker"
)

type Container struct {
	Name    string            `json:"name"`
	Id      string            `json:"id"`
	Labels  map[string]string `json:"labels"`
	Created int64             `json:"created"`
	Status  string            `json:"status"`
	State   string            `json:"state"`
}

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
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Header("Access-Control-Expose-Header", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		if c.Request.Method == "HEAD" {
			c.AbortWithStatus(200)
		}
	})
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
			log.Print(logFile, logOptions)
			logsOut, err := docker.Logs(container, logFile, logOptions)
			if err != nil {
				c.AbortWithError(500, err)
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

		containerGroup.GET("list", func(ctx *gin.Context) {
			cli, err := docker.NewCli()
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}

			containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}

			var result []Container
			for _, c := range containers {
				result = append(result, Container{
					Name:    c.Names[0],
					Id:      c.ID,
					Created: c.Created,
					Labels:  c.Labels,
					Status:  c.Status,
					State:   c.State,
				})
			}
			ctx.JSON(200, result)
		})
	}

	router.LoadHTMLGlob("build/*.html")

	router.GET("/", func(c *gin.Context) {
		c.Header("content-type", "text/html")
		c.HTML(200, "index.html", gin.H{})
	})

	router.Static("/build", "./build")
	router.StaticFile("/favicon.ico", "./build/favicon.ico")
	router.StaticFile("/index.html", "./build/index.html")
	router.StaticFile("/robot.txt", "./build/robot.txt")
	router.StaticFile("/manifest.json", "./build/manifest.json")
	router.StaticFile("/asset-manifest.json", "./build/asset-manifest.json")
	router.Static("/static", "./build/static")

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

	// router.Use(cors.New(cors.Config{
	// 	AllowMethods:  []string{"PUT", "PATCH", "POST", "GET"},
	// 	AllowWildcard: true,
	// 	AllowHeaders:  []string{"*"},
	// 	ExposeHeaders: []string{"Content-Length", "Content-Type"},
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return true
	// 	},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	router.Run(fmt.Sprintf(":%d", port))
}
