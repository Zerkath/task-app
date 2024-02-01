package main

import (
	"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "task_service/rest"
    "task_service/db"
    "time"
    "log"
)

func main() {

    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

	route := gin.Default()

    // TODO: Should find a better mechanism to share the connection between endpoints
    // Currently the way is not clear, as the signature for the endpoints is expected to be of
    // func(c *gin.Context)
    // and the signature for the db connection is of
    // func(db *sqlx.DB) error
    // Potentially signature could be combination of both
    // func(c *gin.Context, db *sqlx.DB) error, but unclear if gin supports this
    db := db.ConnectDB(
        "localhost", 
        5432, 
        "postgres", 
        "postgres", 
        "postgres",
    )
    defer db.Close()

    x := db.MustExec("SELECT 1")
    log.Printf("Result: %v\n", x)

	route.GET("/", rest.Ping)

	route.POST("/task", rest.NewTask)
	route.GET("/task/:id", rest.GetTaskById)
	route.GET("/task/listen/:id", func(c *gin.Context) {
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

	route.GET("/task", rest.GetTasks)

    // disable trusted proxies
    route.SetTrustedProxies([]string{})

    // TODO: Should be env variable or default
	route.Run(":8080")
}
