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

    currPage := types.ToInt(queryParams.Get("page"), 0)
    pageSize := types.ToInt(queryParams.Get("size"), 10)
    if pageSize < 1 {
        pageSize = 1 // prevent division by zero
    }
	args.Limit = int32(pageSize)
	args.Offset = int32(currPage) * args.Limit
	args.Status = repository.NullTaskStatus{}
    taskCountArgs := repository.NullTaskStatus{}
	status := queryParams.Get("status")
	if len(status) > 0 {
		args.Status.Scan(status)
        taskCountArgs.Scan(status)
	} else {
		args.Status.Valid = false
        taskCountArgs.Valid = false
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

    count, err := types.Repository.GetTaskCount(c, taskCountArgs)
    if err != nil {
        log.Println("Error reading task count: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Error reading task count",
        })
        return
    }

    page := types.Page {
        Count: count,
        Page: currPage,
        Data: list,
    }

	c.JSON(http.StatusOK, page)
}

func RemoveTask(c *gin.Context) {

    idx := pgtype.UUID{}
    err := idx.Scan(c.Param("id"))
    if err != nil {
        log.Println("Error parsing id: ", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid id",
        })
        return
    }

    err = types.Repository.DeleteTask(c, idx)
    if err != nil {
        log.Println("Error removing task: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Error removing task",
        })
        return
    }

    c.JSON(http.StatusNoContent, gin.H{
        "message": "Task removed",
    })
}
