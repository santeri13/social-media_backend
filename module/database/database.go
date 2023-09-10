package database

import (
	"database/sql"
	"log"
	"errors"
	
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"

	"social-network/module/structure"
)
const (
	dbPath = "forum.db"
)

// CreateTables creates all the required tables in the database.
func CreateTables(db *sql.DB) error {
	// Create the 'users' table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			user_id TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender TEXT NOT NULL,
			avatar TEXT,
			nickname TEXT,
			about_me TEXT,
			privacy TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'posts' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			category TEXT NOT NULL,
			image TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			post_id INTEGER,
			author_id INTEGER,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (author_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'followers' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS followers (
			id INTEGER PRIMARY KEY,
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			status TEXT CHECK(status IN ('requested', 'accepted')) NOT NULL,
			FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'groups' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY,
			creator_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'group_members' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS group_members (
			id INTEGER PRIMARY KEY,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT CHECK(status IN ('requested', 'accepted')) NOT NULL,
			FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'chat_messages' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id INTEGER PRIMARY KEY,
			sender_id INTEGER NOT NULL,
			receiver_id INTEGER,
			group_id INTEGER,
			message TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create the 'notifications' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			information TEXT NOT NULL,
			message TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_read INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(loginData structure.LoginData) (string) {
	var userID string
	var hashedPassword string
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Retrieve the hashed password from the database based on the provided email
	row := db.QueryRow("SELECT user_id, password FROM users WHERE email = ? OR nickname = ?", loginData.Email, loginData.Email)
	err = row.Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			return "Invalid email or password"
		}
		log.Println("Error retrieving user:", err)
	}
	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginData.Password))
	if err != nil {
		// Password does not match
		return "Invalid email or password"
	}
	return userID
}

func RegisterUser(registrationData structure.RegistrationData) (string){
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if IsUsernameTaken(db,registrationData.Nickname){
		return "Nickname is used"
	}
	if IsEmailTaken(db,registrationData.Email){
		return "Email is used"
	}

	// Generate a UUID for the user
	userID, err := uuid.NewV4()
	if err != nil {
		log.Println("Error generating UUID:", err)
		return ""
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return ""
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (user_id, nickname, first_name, last_name, age, gender, email, password, avatar, about_me, privacy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID.String(), registrationData.Nickname, registrationData.FirstName, registrationData.LastName,
		registrationData.Age, registrationData.Gender, registrationData.Email, hashedPassword, registrationData.Avatar, registrationData.About, "public")
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return ""
	}
	log.Println("User registered successfully:", registrationData.Email)
	return userID.String()
}

// Check if a username is already taken
func IsUsernameTaken(db *sql.DB,username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", username).Scan(&count)
	if err != nil {
		log.Println("Error checking nickname availability:", err)
		return true
	}
	// Username is taken if count > 0
	return count > 0
}

// Check if an email is already registered
func IsEmailTaken(db *sql.DB,email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Println("Error checking email availability:", err)
		return true
	}
	// Email is registered if count > 0
	return count > 0
}

