package main

import (
    "database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gorilla/websocket"

    "social-network/module/database"
	"social-network/module/structure"
	"social-network/module/pages/mainpage"
	"social-network/module/pages/cabinetpage"
	"social-network/module/pages/privatemessagepage"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan structure.Message)
var onlineUsers = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer func() {
		conn.Close()
		// Remove the client from the map of connected clients
		delete(clients, conn)
		// Delete the user from the onlineUsers map
		for uuid, c := range onlineUsers {
			if c == conn {
				delete(onlineUsers, uuid)
				break
			}
		}
	}()

	// Add the new client to the map of connected clients
	clients[conn] = true

    for {
        _, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			// Remove the client from the map of connected clients if there's an error
			delete(clients, conn)

			break
		}

		// Unmarshal the received JSON message into a generic map[string]interface{}
		var types map[string]interface{}
		if err := json.Unmarshal(msg, &types); err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

        switch types["type"] {
			case "login":
				var message structure.LoginData
				err = json.Unmarshal(msg, &message)
				loginData:= structure.LoginData{
					Email:     	message.Email, 
					Password:  	message.Password,
				}
				userID := database.LoginUser(loginData)
				if userID != "" {
					// Add the user to the online users map
					onlineUsers[userID] = conn
				}
				sendMessage(conn, userID)
			case "register":
				var message structure.RegistrationData
				err = json.Unmarshal(msg, &message)
				regData := structure.RegistrationData{
					Nickname:	message.Nickname,
					FirstName:	message.FirstName,
					LastName:  	message.LastName,
					Age:      	message.Age, 
					Gender:    	message.Gender,
					Email:     	message.Email, 
					Password:  	message.Password,
				}
				userID := database.RegisterUser(regData)
				if userID != "" {
					onlineUsers[userID] = conn
				}
				sendMessage(conn, userID)
			case "log_out":
				var message structure.Message
				err = json.Unmarshal(msg, &message)
				delete(onlineUsers, message.UUID)
			case "addUserToConnection":
				var message structure.Message
				err = json.Unmarshal(msg, &message)
				pageData := structure.Message{
					UUID: message.UUID,
				}
				if pageData.UUID != "" {
					onlineUsers[pageData.UUID] = conn
				}
			case "post":
				var message structure.Post
				err = json.Unmarshal(msg, &message)
				post := structure.Post{
					UserID:		message.UserID,
					Title:		message.Title,
					Content:  	message.Content,
					Category:   message.Category,
				}
				log.Println(post)
				mainpage.PostCreation(post)
			case "comment":
				var message structure.Comment
				err = json.Unmarshal(msg, &message)
				comment := structure.Comment{
					UserID:		message.UserID,
					Content:	message.Content,
					PostID:  	message.PostID,
				}
				log.Println(comment)
				mainpage.CommentCreation(comment)
			case "posts":
				sendPosts(conn,mainpage.GetPostsFromDatabase())
			case "change":
				var message structure.UserData
				err = json.Unmarshal(msg, &message)
				userData:= structure.UserData{
					UserID:		message.UserID,
					Nickname:	message.Nickname,
					FirstName:	message.FirstName,
					LastName:  	message.LastName,
					Age:      	message.Age, 
					Gender:    	message.Gender,
					Email:     	message.Email, 
				}
				cabinetpage.UpdateUserData(userData)
			case "userlist":
				var message structure.Message
				err = json.Unmarshal(msg, &message)
				pageData := structure.Message{
					UUID: message.UUID,
				}
				// Get the list of users from the database
				users := privatemessage.GetUsers(pageData.UUID)

				// Iterate through the list of users and check if they are online
				for i := range users {
					fmt.Println(users[i].UUID)
					if _, ok := onlineUsers[users[i].UUID]; ok {
						users[i].Activity = "online"
					} else {
						users[i].Activity = "offline"
					}
				}
				
				sendUsersData(conn, users)
			case "user_uuid":
				var message structure.Message
				err = json.Unmarshal(msg, &message)
				sendMessage(conn, message.UUID)
			case "showmessage":
				var message structure.UserMessageData
				err = json.Unmarshal(msg, &message)
				// Convert the local RegistrationData struct to module.RegistrationData
				retriveMessageData := structure.UserMessageData{
					UserID:   message.UserID,
					Nickname: message.Nickname,
					Offset:   message.Offset,
				}
				messages := privatemessage.ProvideMessages(retriveMessageData, 10)
				privateMessages := make([]structure.PrivateMessages, len(messages))
				for i, message := range messages {
					var userID string
					// Open a connection to the SQLite database
					db, err := sql.Open("sqlite3", "./forum.db")
					if err != nil {
						log.Println(err)
					}
					defer db.Close()
					row := db.QueryRow("SELECT user_id FROM users WHERE nickname = ?", message.SenderId)
					row.Scan(&userID)
					privateMessages[i] = structure.PrivateMessages{
						Content:  message.Content,
						Time:     message.Time,
						SenderId: userID,
					}
				}
				sendPrivateMessages(conn, privateMessages)
			case "sendPrivateMessage":
				var message structure.PrivateMesssageSend
				err = json.Unmarshal(msg, &message)
				messageData := structure.PrivateMesssageSend{
					UserID:   message.UserID,
					Nickname: message.Nickname,
					Content:  message.Content,
				}
				privatemessage.SendMessageToUser(messageData)
			case "cabinet":
				var message structure.Message
				err = json.Unmarshal(msg, &message)
				pageData:=structure.Message{
					UUID: message.UUID,
				}
				db, err := sql.Open("sqlite3", "./forum.db")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()
				var nickname string
				var age int
				var gender string
				var first_name string
				var last_name string
				var email string
				// Retrieve the hashed password from the database based on the provided email
				row := db.QueryRow("SELECT nickname, age, gender, first_name, last_name, email FROM users WHERE user_id = ?", pageData.UUID)
				err = row.Scan(&nickname, &age, &gender, &first_name, &last_name, &email)
				if err != nil {
					log.Println("Error retrieving user UUID:", err)
				}
				userData:= structure.UserData{
						Nickname:  nickname,
						FirstName: first_name,
						LastName:  last_name,
						Age:       age,
						Gender:    gender,
						Email:     email,
				}
				sendUserData(conn, userData)
		}
    }
}

// Function to send a message back to the sender
func sendMessage(conn *websocket.Conn, message string) {
	err := conn.WriteJSON(structure.Message{UUID: message})
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func sendUserData(conn *websocket.Conn, message structure.UserData) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func sendUsersData(conn *websocket.Conn, message []structure.UserData) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}


// Function to send a message back to the sender
func sendPosts(conn *websocket.Conn, message []structure.Post) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func sendPrivateMessages(conn *websocket.Conn, message []structure.PrivateMessages) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

// Function to broadcast a message to all connected clients except the sender
func broadcastMessage(sender *websocket.Conn, message string) {
	for conn := range clients {
		if conn != sender {
			err := conn.WriteJSON(structure.Message{UUID: message})
			if err != nil {
				log.Println("Failed to broadcast message:", err)
			}
		}
	}
}

func main() {
    db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.CreateTables(db)

	http.HandleFunc("/websocket", handleWebSocketConnection)

	// Serve the frontend files
	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	fmt.Println("WebSocket server started. Listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}