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
		rows, err := db.Query("SELECT id, nickname, user_id FROM users WHERE user_id != ? ORDER BY (SELECT MAX(created_at) FROM chat_messages WHERE sender_id = users.id OR receiver_id = users.id) DESC", user_id)
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
