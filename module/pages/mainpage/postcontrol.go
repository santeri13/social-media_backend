package mainpage

import (
	"database/sql"
	"log"

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
	_, err = db.Exec("INSERT INTO posts (title, content, category, user_id) VALUES (?, ?, ?, ?)", postData.Title, postData.Content, postData.Category, userID)

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

func getPostsFromDatabase(postCounter int) structure.Post {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// Retrieve all posts from the database
	rows, err := db.Query(`
		SELECT p.id, p.author_id, p.title, p.content, p.category, c.id, c.author_id, c.content
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id ORDER BY p.id
	`)
	if postCounter > 0 {
		rows, err = db.Query(`
		SELECT p.id, p.author_id, p.title, p.content, p.category, c.id, c.author_id, c.content
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id ORDER BY p.id DESC LIMIT 1`)
	}
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// Map to store posts and their comments
	postsMap := make(map[string]*structure.Post)
	// Iterate over the rows and populate the postsMap
	for rows.Next() {
		var postID, commentID int
		var postUserID, postTitle, postContent, postCategory, commentUserID, commentContent, postUserNickname, commentUserNickname string
		err := rows.Scan(&postID, &postUserID, &postTitle, &postContent, &postCategory, &commentID, &commentUserID, &commentContent)
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
			post = &Post{
				ID:       postID,
				UserID:   postUserID,
				Title:    postTitle,
				Content:  postContent,
				Category: postCategory,
				Nickname: postUserNickname,
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
			comment := Comment{
				ID:      commentID,
				UserID:  commentUserID,
				Content: commentContent,
				Nickname: commentUserNickname,
			}
			post.Comments = append(post.Comments, comment)
		}
	}
	// Collect the posts from the map
	posts := make(*structure.Post, 0, len(postsMap))
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