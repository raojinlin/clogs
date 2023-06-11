package main

import (
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/raojinlin/container-go/docker"
	"io"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsWriter struct {
	wsConn *websocket.Conn
}

func (ws *WsWriter) Write(p []byte) (n int, err error) {
	err = ws.wsConn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
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

			wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			defer wsConn.Close()

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
				wsConn.Close()
				return
			}

			go func(logsOut io.ReadCloser) {
				for {
					_, msg, err := wsConn.ReadMessage()
					if err != nil {
						wsConn.Close()
						logsOut.Close()
						return
					}
					fmt.Println("recv: ", msg, err)
				}
			}(logsOut)
			wsWriter := &WsWriter{wsConn: wsConn}
			stdcopy.StdCopy(wsWriter, wsWriter, logsOut)
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
