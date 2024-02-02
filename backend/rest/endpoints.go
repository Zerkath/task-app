package rest

import (
	"log"
	// "math/rand"
	"net/http"
	"task-service/db"
	"task-service/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Simple health check
func Ping(c *gin.Context) {

	err := db.TestConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service up",
	})
}

func NewTask(c *gin.Context) {

	vUuid := uuid.New()
	r, e := db.CONNECTION.Exec("INSERT INTO task (id, status) VALUES ($1, $2)", vUuid, "running")

	if e != nil {
		log.Println("Error inserting task: ", e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error inserting task",
		})
		return
	}

	log.Println("Result: ", r)

	c.JSON(http.StatusOK, gin.H{
		"id":     vUuid.String(),
		"status": "running",
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

	rows, err := db.CONNECTION.Query("SELECT status, completed_at FROM task WHERE id = $1", id.String())
	if err != nil {
		log.Println("Error reading task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading task",
		})
		return
	}

	if !rows.Next() {
		log.Println("No task found with id: ", id)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}

	task := types.Task{Id: id, Status: "running"}

    // in microseconds
	var completedAt int64

	rows.Scan(&task.Status, &completedAt)

	if completedAt > 0 {
		log.Println("Task completed at: ", completedAt)
		task.CompletedAt = int32(completedAt / 1000 / 1000) // convert to seconds, for unix timestamp
	}

	c.JSON(http.StatusOK, task)
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
