package main

import (
	"github.com/gin-gonic/gin"
    "task_service/rest"
)

func main() {

	route := gin.Default()

	route.GET("/", rest.Ping)
	route.POST("/task", rest.NewTask)
	route.GET("/task/:id", rest.GetTaskById)
	route.GET("/task", rest.GetTasks)

    // disable trusted proxies
    route.SetTrustedProxies([]string{})

	route.Run(":8080")
}
