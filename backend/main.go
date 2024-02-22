package main

import (
	"context"
	"log"
	"net/http"
	"task-service/repository"
	"task-service/rest"
	"task-service/types"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Should be a list of allowed origins
			return true
		},
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
	conn, err := pgxpool.New(ctx, "host=localhost user=postgres password=postgres dbname=tasks sslmode=disable")

	if err != nil {
		log.Fatalln("Error connecting to database: ", err)
	}
	defer conn.Close()

	types.Conn = conn
	types.Repository = repository.New(conn)

	router.GET("/", rest.Ping)

	router.POST("/task", rest.NewTask)
	router.GET("/task", rest.GetTasks)

	router.GET("/task/:id", rest.GetTaskById)
	router.DELETE("/task/:id", rest.RemoveTask)

	router.GET("/task/listen", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Error upgrading connection: ", err)
			return
		}
		defer conn.Close()

		// idx := pgtype.UUID{}

		var data types.ListenList

		go func() {
			for {
				log.Println("Waiting for new listen list")
				var newList types.ListenList
				err := conn.ReadJSON(&newList) // read waits for new message, and only returns if message is received

				// Verify that list was updated, or new message was received between reads
				if err != nil {
					println("Error reading from socket: ", err)
					// connection most likely killed
					return
				}

				log.Println("Received listen list: ", newList)
				data = newList
			}
		}()

		for {
			ids := []pgtype.UUID{}
			for _, id := range data {
				idx := pgtype.UUID{}
				err := idx.Scan(id.String())
				if err != nil {
					log.Printf("Error scanning id: %s %e", id, err)
					continue
				}
				ids = append(ids, idx)
			}

			tasks, err := types.Repository.GetTasksByIds(ctx, ids)
			if err != nil {
				log.Printf("Error getting tasks from db with ids %v %e/n", data, err)
				continue
			}
			if len(tasks) == 0 {
				conn.WriteJSON([]types.Task{})
			} else {
				err = conn.WriteJSON(tasks)

				if err != nil {
					log.Printf("Error writing to socket: %e", err)
					// maybe connection was killed?
					// should exit the loop and kill the other goroutine
					return
				}
			}
            time.Sleep(5 * time.Second)
		}
	})

	// disable trusted proxies
	router.SetTrustedProxies([]string{})

	go queueHandler(ctx)
	// TODO: Should be env variable or default
	router.Run(":8080")

}

type ErrorMessage struct {
	message string
	mType   string
}

// Daemon type of process that will query a batch of tasks and process them
// should be called in main as a go routine
// note will pretend to handle tasks
// tasks will be read, and randomly changed to completed, failed
// if failed, will be requeued and restart count is increased
func queueHandler(ctx context.Context) {

	args := repository.GetUncompletedTasksParams{
		Limit:  10,
		Offset: 0,
	}

	for {
		tasks, err := types.Repository.GetUncompletedTasks(ctx, args)
		if err != nil {
			log.Println("Error getting tasks: ", err)
		}

		for _, task := range tasks {
			readableUUID := uuid.UUID(task.ID.Bytes).String()
			log.Println("Processing task: ", readableUUID)

			random := time.Now().Nanosecond() % 2

			restarts := task.Restarts
			if task.Status == "failed" {
				restarts.Int32++
			}

			params := repository.UpdateTaskParams{
				ID:          task.ID,
				Status:      "running",
				Restarts:    restarts,
				CompletedAt: task.CompletedAt,
			}

			_, err = types.Repository.UpdateTask(ctx, params)
			if err != nil {
				log.Println("Error updating task status: ", err)
			}
			time.Sleep(3 * time.Second) // simulate processing time for task

			if random == 0 {
				log.Println("Task completed: ", readableUUID)
				params.Status = "completed"
				params.CompletedAt.Scan(time.Now())
				_, err = types.Repository.UpdateTask(ctx, params)
				if err != nil {
					log.Println("Error updating task status: ", err)
				}
			} else {
				log.Println("Task failed: ", readableUUID)
				params.Status = "failed"
				_, err = types.Repository.UpdateTask(ctx, params)
				if err != nil {
					log.Println("Error updating task status: ", err)
				}
			}

			time.Sleep(3 * time.Second) // simulate delay between starting next task
		}
		time.Sleep(30 * time.Second) // process batches every 30 seconds
	}
}
