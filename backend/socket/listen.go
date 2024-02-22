package socket

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"task-service/types"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Should be a list of allowed origins
		return true
	},
}

func Listen(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}
	defer conn.Close()

	var data types.ListenList
	// this channel is used for immediate updates
	// we can run two goroutines, one to listen for new messages
	// one just sends to the available list in 'data'
	var channel = make(chan types.ListenList)

	go func() {
		for {
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
			channel <- newList
		}
	}()

	// this goroutine will listen for new messages in the channel
	// for immediate updates, improving user experience
	go func() {
		for {
			list := <-channel
			err := pushToSocket(c, conn, list)
			if err != nil {
				return
			}
		}
	}()

	// basic loop to check the given list occasionally
	// and push to the socket
	for {
		err := pushToSocket(c, conn, data)
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func pushToSocket(ctx context.Context, conn *websocket.Conn, data types.ListenList) error {
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
		return nil
	}

	if len(tasks) == 0 {
		conn.WriteJSON([]types.Task{})
	} else {
		err = conn.WriteJSON(tasks)

		if err != nil {
			log.Printf("Error writing to socket: %e", err)
			// maybe connection was killed?
			// should exit the loop and kill the other goroutine
			return err
		}
	}
	return nil
}
