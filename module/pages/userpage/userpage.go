package userpage

import (
	"database/sql"
	"log"
	"errors"

	_ "github.com/mattn/go-sqlite3"

	"social-network/module/structure"
)

func UserPosts(userUUID string) []structure.Post {
	var posts []structure.Post
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// Retrieve all posts from the database
	rows, err:= db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.image , c.id, c.user_id
		FROM posts p
		LEFT JOIN users c ON c.id = p.user_id WHERE c.user_id = ? ORDER BY p.id
	`, userUUID)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var post structure.Post
		err := rows.Scan(&post.UserID, &post.Title, &post.Content, &post.Category)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	return posts
}

func FollowUser(followUser structure.UserData){
	var following_id string
	var follower_id string
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	row := db.QueryRow("SELECT id FROM users WHERE user_id = ? OR id = ?", followUser.UserID, followUser.UserID)
	err = row.Scan(&following_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			log.Println("Invalid following_id")
		}
		log.Println("Error retrieving following_id:", err)
	}

	row = db.QueryRow("SELECT id FROM users WHERE user_id = ? OR id = ?", followUser.UUID, followUser.UUID)
	err = row.Scan(&follower_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			log.Println("Invalid follower_id")
		}
		log.Println("Error retrieving follower_id:", err)
	}

	defer db.Close()

	sqlStmt := `SELECT following_id, follower_id FROM followers WHERE following_id = ? and follower_id = ?`
	err = db.QueryRow(sqlStmt, following_id, follower_id).Scan(&following_id, &follower_id)
	if err != nil {
        if err != sql.ErrNoRows {
            // a real error happened! you should change your function return
            // to "(bool, error)" and return "false, err" here
            log.Print(err)
        }

        _, err = db.Exec("INSERT INTO followers (follower_id , following_id, status) VALUES (?, ?, ?)", follower_id, following_id, "accepted")
		if err != nil {
			log.Println("Error inserting follower into database:", err)
		}
		log.Println("User", follower_id,"follow",following_id)
    } else{
		stmt, err := db.Prepare("DELETE FROM followers WHERE follower_id = ? AND following_id = ?")
		if err != nil {
			log.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(follower_id, following_id)
		if err != nil {
			log.Println("Error inserting follower into database:", err)
		}
		log.Println("User", follower_id,"delete follow",following_id)
	}
}