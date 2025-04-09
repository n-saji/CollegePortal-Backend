package utils

import (
	"CollegeAdministration/config"
	"CollegeAdministration/daos"
	"CollegeAdministration/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsId = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex // Mutex to protect the clients map
	broadcast = make(chan string)
)

func InitiateWebSockets() {

	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Error writing message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func HandleConnections(c *gin.Context) {
	w := c.Writer
	r := c.Request
	params := c.Params
	id := params.ByName("id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Protect the clients map with a mutex
	clientsMu.Lock()
	clients[conn] = true
	clientsId[id] = conn
	clientsMu.Unlock()
	log.Println("New client connected", id)

	// Listen for messages from the client (optional)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
	}

	// Remove the client when it disconnects
	clientsMu.Lock()
	delete(clients, conn)
	delete(clientsId, id)
	clientsMu.Unlock()
	log.Println("Client disconnected")
}

func SendMessageToClientId(id string, msg models.Messages) bool {
	clientsMu.Lock()
	jsonMsg, _ := json.Marshal(msg)
	conn, ok := clientsId[id]
	if ok {
		conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
		clientsMu.Unlock()
		return true
	}
	clientsMu.Unlock()
	return false
}

func SendMessageToConnectedClients(msg, author, title string) {
	clientsMu.Lock()
	MSG := &models.Messages{}
	MSG.ID = uuid.New()
	MSG.Messages = msg
	MSG.Author = author
	MSG.Title = title
	MSG.CreatedAt = time.Now().Unix()
	jsonMsg, _ := json.Marshal(MSG)

	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
	}
	clientsMu.Unlock()
}

// message to send , account type to send to, account to skip
func StoreMessages(title, message, account_type, sender_id string) {

	dbConn := config.DBInit()
	db := daos.New(dbConn)
	author := &models.Account{}

	ids, err := db.GetAccountIDsByType(account_type)
	if err != nil {
		fmt.Println("error while fetching instructor ids" + err.Error())
		return
	}
	if sender_id != "" {
		skip_account_uuid, err := uuid.Parse(sender_id)
		if err != nil {
			fmt.Println("error while parsing skip account" + err.Error())
			return
		}
		author, err = db.GetAccountNameById(skip_account_uuid)
		if err != nil {
			fmt.Println("error while fetching author" + err.Error())
			return
		}
	}
	msg := &models.Messages{}
	msg.Messages = message
	msg.IsRead = false
	if author != nil {
		msg.Author = author.Name
	} else {
		msg.Author = "nil"
	}
	msg.Title = title
	msg.CreatedAt = time.Now().Unix()
	for _, id := range ids {
		if sender_id != "" && id.Id.String() == sender_id {
			continue
		}
		msg.ID = uuid.New()
		msg.AccountID = id.Id
		err := db.InsertIntoMessages(msg)
		if err != nil {
			fmt.Println("error while inserting message" + err.Error())
			return
		}
	}
	// clientsMu.Lock()
	// jsonMsg, _ := json.Marshal(msg)
	// for acntId, conn := range clientsId {
	// 	if acntId == sender_id {
	// 		continue
	// 	}
	// 	conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
	// 	// db.UpdateMessageStatusForAccountId(uuid.MustParse(acntId))
	// }
	// clientsMu.Unlock()
	defer config.CloseDB(dbConn)
}
