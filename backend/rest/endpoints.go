package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"task_service/types"
)

// Simple health check, should return 200 if all dependencies are up
// Currently not implemented
// TODO: Implement health check
// TODO: Add logging
// TODO: Test db connection
func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Service up",
	})
}

// TODO: Should create task in the DB
// Also a coroutine that updates the status of the task
func NewTask(c *gin.Context) {

	c.JSON(http.StatusOK, types.Task{
		Id:     uuid.MustParse("63dc6aa8-3db1-4ceb-b578-3944b4947f1a"),
		Status: "running",
	})
}

// TODO: Should read task from DB
// Return the result as is
func GetTaskById(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("Error parsing id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id",
		})
		return
	}

	log.Println("Received id: ", id)

	c.JSON(http.StatusOK, types.Task{Id: id, Status: "running"})
}

// TODO: Should read a page of results from the DB
// Also provide some filtering options
// Sorting by date, status, etc
// Filtering by status, etc
// Pagination
func GetTasks(c *gin.Context) {

	list := []types.Task{
		{Id: uuid.MustParse("63dc6aa8-3db1-4ceb-b578-3944b4947f1a"), Status: "running"},
		{Id: uuid.MustParse("e60f4a26-e40a-49b5-ab7e-7b6c2396791a"), Status: "running"},
		{Id: uuid.MustParse("613f7102-4f6a-4d7c-8aeb-5df8a250e1ff"), Status: "running"},
	}
	c.JSON(http.StatusOK, list)
}
