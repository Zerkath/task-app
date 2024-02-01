package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
    "github.com/google/uuid"
	"task_service/types"
)

func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Service up",
	})
}

func NewTask(c *gin.Context) {

	c.JSON(http.StatusOK, types.Task{
		Id:     uuid.MustParse("63dc6aa8-3db1-4ceb-b578-3944b4947f1a"),
		Status: "running",
	})
}

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

func GetTasks(c *gin.Context) {

	list := []types.Task{
		{Id: uuid.MustParse("63dc6aa8-3db1-4ceb-b578-3944b4947f1a"), Status: "running"},
		{Id: uuid.MustParse("e60f4a26-e40a-49b5-ab7e-7b6c2396791a"), Status: "running"},
		{Id: uuid.MustParse("613f7102-4f6a-4d7c-8aeb-5df8a250e1ff"), Status: "running"},
	}
	c.JSON(http.StatusOK, list)
}
