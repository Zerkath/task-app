package main

import (
	"github.com/gin-gonic/gin"
    "task_service/rest"
    "task_service/db"
    "log"
)

func main() {

	route := gin.Default()

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
	route.GET("/task", rest.GetTasks)

    // disable trusted proxies
    route.SetTrustedProxies([]string{})

	route.Run(":8080")
}
