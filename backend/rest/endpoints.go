package rest

import (
	"log"
	"net/http"
	"task-service/repository"
	"task-service/types"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Simple health check
func Ping(c *gin.Context) {

	err := types.Conn.Ping(c)
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

	task, err := types.Repository.NewTask(c)

	if err != nil {
		log.Println("Error inserting task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error inserting task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTaskById(c *gin.Context) {

	idx := pgtype.UUID{}
	err := idx.Scan(c.Param("id"))
	if err != nil {
		log.Println("Error parsing id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id",
		})
		return
	}

	task, err := types.Repository.GetTaskById(c, idx)
	if err == pgx.ErrNoRows {
		log.Println("Error reading task: ", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No task found",
		})
		return
	} else if err != nil {
		log.Println("Error reading task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {

	args := repository.GetTasksParams{}

	queryParams := c.Request.URL.Query()

	args.Limit = int32(types.ToInt(queryParams.Get("limit"), 10))
	args.Offset = int32(types.ToInt(queryParams.Get("offset"), 0))
	args.Status = repository.NullTaskStatus{}
	status := queryParams.Get("status")
	if len(status) > 0 {
		args.Status.Scan(status)
	} else {
		args.Status.Valid = false
	}

	list, err := types.Repository.GetTasks(c, args)
	if err == pgx.ErrNoRows || len(list) == 0 {
		c.JSON(http.StatusOK, make([]repository.Task, 0))
		return
	} else if err != nil {
		log.Println("Error reading tasks: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading tasks",
		})
		return
	}

	c.JSON(http.StatusOK, list)
}
