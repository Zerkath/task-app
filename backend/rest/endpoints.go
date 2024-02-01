package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task_service/types"
)

func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Service up",
	})
}

func NewTask(c *gin.Context) {

	c.JSON(http.StatusOK, types.Task{
		ID:     "1182881-123123-123123-123123",
		Status: "running",
	})
}

func GetTaskById(c *gin.Context) {

	id := c.Param("id")
	log.Println("Received id: ", id)

	c.JSON(http.StatusOK, types.Task{ID: id, Status: "running"})
}

func GetTasks(c *gin.Context) {

	list := []types.Task{
		{ID: "1182881-444444-123123-123123", Status: "running"},
		{ID: "1182882-123123-123123-123123", Status: "running"},
		{ID: "1182881-123123-555555-123123", Status: "running"},
	}
	c.JSON(http.StatusOK, list)
}
