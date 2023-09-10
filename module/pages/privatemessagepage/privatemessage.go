package privatemessage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"social-network/module/structure"
)


func GetUsers(user_id string) []structure.UserData {
	var users []structure.UserData
	var userID string
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", user_id)
	err = row.Scan(&userID)
	if err != nil {
		log.Fatal("Retrive UUID to check",err)
	}
	hasRecord := false
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_messages WHERE sender_id = ? OR receiver_id = ?)", userID, userID).Scan(&hasRecord)
	if err != nil {
		log.Fatal("Check if record exists ",err)
	}
	if hasRecord {
		rows, err := db.Query("SELECT id, nickname, user_id FROM users WHERE user_id != ? ORDER BY (SELECT MAX(timestamp) FROM chat_messages WHERE sender_id = users.id OR receiver_id = users.id) DESC", user_id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var user structure.UserData
			err := rows.Scan(&user.UserID, &user.Nickname, &user.UUID)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
	} else {
		rows, err := db.Query("SELECT id, nickname, user_id FROM users WHERE user_id != ? ORDER BY nickname ASC", user_id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var user structure.UserData
			err := rows.Scan(&user.UserID, &user.Nickname, &user.UUID)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
	}
	return users
}

func ProvideMessages(userData structure.UserMessageData) ([]structure.PrivateMessages){
	var userID string
	var otherUserID string
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", userData.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}

	row = db.QueryRow("SELECT id FROM users WHERE nickname = ?", userData.Nickname)
	err = row.Scan(&otherUserID)
	if err != nil {
		log.Println("Error retrieving user nickname:", err)
	}

	var messages []structure.PrivateMessages
	query := `
	SELECT message, timestamp, sender_id, receiver_id FROM chat_messages 
	WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
	ORDER BY timestamp DESC`
	rows, err := db.Query(query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		log.Println("Error retriveing messages", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message structure.PrivateMessages
		err := rows.Scan(&message.Content, &message.Time, &message.SenderId, &message.RecipientId)
		if err != nil {
			log.Println("Invalid UUID")
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		log.Println("Invalid UUID")
	}
	log.Println("Show message", messages)
	return messages
}

func SendMessageToUser(privateMessage  structure.PrivateMesssageSend){
	var userID string
	var otherUserID string
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", privateMessage.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}
	row = db.QueryRow("SELECT id FROM users WHERE nickname = ?", privateMessage.Nickname)
	err = row.Scan(&otherUserID)
	if err != nil {
		log.Println("Error retrieving user nickname:", err)
	}
	_, err = db.Exec("INSERT INTO chat_messages (message, sender_id, receiver_id) VALUES (?, ?, ?)",
		privateMessage.Content, userID, otherUserID)
	if err != nil {
		log.Println("Error inserting user into database:", err)
	}
	
	
	log.Println("Message send:", privateMessage.Content)
}

func ProvideGroupMessages(userData structure.UserMessageData) ([]structure.PrivateMessages){
	var userID string
	var otherUserID string
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", userData.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}

	row = db.QueryRow("SELECT id FROM groups WHERE title = ?", userData.Nickname)
	err = row.Scan(&otherUserID)
	if err != nil {
		log.Println("Error retrieving user nickname:", err)
	}

	var messages []structure.PrivateMessages
	query := `
	SELECT message, timestamp, sender_id, group_id FROM chat_messages 
	WHERE (sender_id = ? AND group_id = ?) OR (sender_id = ? AND group_id = ?) 
	ORDER BY timestamp DESC`
	rows, err := db.Query(query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		log.Println("Error retriveing messages", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message structure.PrivateMessages
		err := rows.Scan(&message.Content, &message.Time, &message.SenderId, &message.RecipientId)
		if err != nil {
			log.Println("Invalid UUID")
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		log.Println("Invalid UUID")
	}
	log.Println("Show message", messages)
	return messages
}

func SendMessageToGroup(privateMessage  structure.PrivateMesssageSend){
	log.Println(privateMessage)
	var userID string
	var otherUserID string
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", privateMessage.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}
	row = db.QueryRow("SELECT id FROM groups WHERE title = ?", privateMessage.Nickname)
	err = row.Scan(&otherUserID)
	if err != nil {
		log.Println("Error retrieving user nickname:", err)
	}
	_, err = db.Exec("INSERT INTO chat_messages (message, sender_id, group_id) VALUES (?, ?, ?)",
		privateMessage.Content, userID, otherUserID)
	if err != nil {
		log.Println("Error inserting user into database:", err)
	}
	
	
	log.Println("Message send:", privateMessage.Content)
}