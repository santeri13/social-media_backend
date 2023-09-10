package mainpage

import (
	"database/sql"
	"log"
	"errors"
	"sort"

	_ "github.com/mattn/go-sqlite3"

	"social-network/module/structure"
)


func PostCreation(postData structure.Post) {
	var userID string

	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", postData.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO posts (title, content, category, user_id, image) VALUES (?, ?, ?, ?, ?)", postData.Title, postData.Content, postData.Category, userID, postData.ImagePath)

	if err != nil {
		log.Println("Error inserting post into database:", err)
	}

	log.Println("Post is created:", postData.Title)
}

func CommentCreation(commentData structure.Comment) {
	var userID string
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", commentData.UserID)
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error retrieving user UUID:", err)
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO comments (content, post_id, author_id) VALUES (?, ?, ?)", commentData.Content, commentData.PostID, userID)

	if err != nil {
		log.Println("Error inserting comment into database:", err)
	}

	log.Println("Comment is created:", commentData.Content)
}

func GetPostsFromDatabase() []structure.Post {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// Retrieve all posts from the database
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.image , c.id, c.author_id, c.content
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id ORDER BY p.id
	`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// Map to store posts and their comments
	postsMap := make(map[int]*structure.Post)
	// Iterate over the rows and populate the postsMap
	for rows.Next() {
		var postID, commentID int
		var postUserID, postTitle, postContent, postCategory, postImagePath, commentUserID, commentContent, postUserNickname, commentUserNickname string
		err := rows.Scan(&postID, &postUserID, &postTitle, &postContent, &postCategory, &postImagePath, &commentID, &commentUserID, &commentContent)
		if err != nil {
			log.Println(err)
		}
		row := db.QueryRow("SELECT nickname FROM users WHERE id = ?", postUserID)
		err = row.Scan(&postUserNickname)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println( err)
			}
			log.Println("Error retrieving user nickname:", err)
		}
		// Check if the post already exists in the map
		post, ok := postsMap[postID]
		if !ok {
			// Create a new post if it doesn't exist
			post = &structure.Post{
				ID:       postID,
				UserID:   postUserID,
				Title:    postTitle,
				Content:  postContent,
				Category: postCategory,
				Nickname: postUserNickname,
				ImagePath: postImagePath,
			}
			postsMap[postID] = post
		}
		row = db.QueryRow("SELECT nickname FROM users WHERE id = ?", commentUserID)
		err = row.Scan(&commentUserNickname)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println( err)
			}
			log.Println("Error retrieving user nickname for comments:", err)
		}
		// Append the comment to the post's comments slice
		if commentID != 0 {
			comment := structure.Comment{
				ID:      commentID,
				UserID:  commentUserID,
				Content: commentContent,
				Nickname: commentUserNickname,
			}
			post.Comments = append(post.Comments, comment)
		}
	}
	// Collect the posts from the map
	posts := make([]structure.Post, 0, len(postsMap))
	for _, post := range postsMap {
		posts = append(posts, *post)
	}
	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	sort.Sort(employeeList(posts))
	return posts
}

type employeeList []structure.Post
func (e employeeList) Len() int {
	return len(e)
}
func (e employeeList) Less(i, j int) bool {
	return e[i].ID < e[j].ID
}
func (e employeeList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func GetNotifications(UUID string) []structure.Notification {
	var notifications []structure.Notification
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Retrieve all posts from the database
	rows, err := db.Query(`SELECT id, type, information, message, timestamp FROM notifications WHERE user_id = ? AND is_read = 0`, UUID)
	if err != nil {
		log.Println("Error retriveing notifications", err)
	}

	for rows.Next() {
		var notification structure.Notification
		err := rows.Scan(&notification.ID ,&notification.Type, &notification.Information, &notification.Message, &notification.Time)
		if err != nil {
			log.Println("Wrong in retriving notifications")
		}
		notifications = append(notifications, notification)
	}
	_, err = db.Exec("UPDATE notifications SET is_read = 1")
	if err != nil {
		log.Println("Error changing is_read in notifications:", err)
	}
	log.Println("Show notifications", notifications)
	return notifications
}

func DeleteNotification(id int){
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM notifications WHERE id = ", id)
	if err != nil {
		log.Println("Error deleteing notification:", err)
	}

	log.Println("Notification deleted")
}

func CreateNotification(userData structure.Notification){
	var message string
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	if (userData.Type == "follow"){
		row := db.QueryRow("SELECT username FROM users WHERE user_id = ? OR id = ?", userData.UUID, userData.UUID)
		err = row.Scan(&message)
		if err != nil {
			log.Fatal("Retrive UUID to check",err)
		}
		message = "User "+message+" like to follow you"
	}else if (userData.Type == "group"){
		message = "Invitation to group "+message
	}

	_, err = db.Exec("INSERT INTO notifications (user_id, type, message, is_read, information) VALUES (?, ?, ?, ?, ?)",  userData.UUID, userData.Type, message, 0, userData.Information)
	if err != nil {
		log.Println("Error inserting group member into database:", err)
	}

	log.Println("Notification created")
}