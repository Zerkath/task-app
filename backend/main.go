package main

import (
	"context"
	"log"
	"task-service/repository"
	"task-service/rest"
    "task-service/types"
	"time"

    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

func main() {

    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

	router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return true // TODO: Should be a list of allowed origins
        },
        MaxAge: 12 * time.Hour,
    }))


    ctx := context.Background()
    conn, err := pgx.Connect(ctx, "host=localhost user=postgres password=postgres dbname=tasks sslmode=disable")

    if err != nil {
        log.Fatalln("Error connecting to database: ", err)
    }
    defer conn.Close(ctx)

    types.Conn = conn
    types.Repository = repository.New(conn)
    
	router.GET("/", rest.Ping)
	router.POST("/task", rest.NewTask)
	router.GET("/task/:id", rest.GetTaskById)
    router.DELETE("/task/:id", rest.RemoveTask)
	router.GET("/task/listen/:id", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println("Error upgrading connection: ", err)
            return
        }
        defer conn.Close()

        for {
            // TODO:
            // Should read DB for task status
            // if not completed, should send message
            // once done send final message and close connection
            // Also this could be moved elsewhere
            conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
            time.Sleep(1 * time.Second)
        }
    })

	router.GET("/task", rest.GetTasks)

    // disable trusted proxies
    router.SetTrustedProxies([]string{})

    // TODO: Should be env variable or default
	router.Run(":8080")
}
