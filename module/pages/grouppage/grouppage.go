package grouppage
import (
	"database/sql"
	"log"
	"errors"

	"social-network/module/structure"
)

func CreateGroup(groupData structure.Group) {
	var userID string
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE user_id = ? OR id = ?", groupData.UUID, groupData.UUID)
	err = row.Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			log.Println("Invalid id")
		}
		log.Println("Error retrieving id:", err)
	}

	_, err = db.Exec("INSERT INTO groups (creator_id, title, description) VALUES (?, ?, ?)", 
		&userID, groupData.Name, groupData.Description)
	if err != nil {
		log.Println("Error inserting group into database:", err)
	}
	log.Println("Group registered successfully:", groupData.Name)

	InsertUser(groupData.Name, userID)
}

func GetGroups(groupInfromaion structure.Message) []structure.Group{
	var groups []structure.Group
	var userID int 
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT title, description FROM groups `)
	if groupInfromaion.Page == "group"{
		rows, err = db.Query(`SELECT title, description FROM groups `)
	}else if groupInfromaion.Page  == "messages"{
		row := db.QueryRow("SELECT id FROM users WHERE user_id = ?", groupInfromaion.UUID)
		err = row.Scan(&userID)
		if err != nil {
			log.Fatal("Retrive UUID to check",err)
		}
		rows, err = db.Query(`SELECT g.title, g.description
		FROM groups AS g
		JOIN group_members AS gm ON g.id = gm.group_id
		WHERE gm.user_id = ?`, userID)
	}
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var group structure.Group
		err := rows.Scan(&group.Name, &group.Description)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	return groups
}

func InsertUser(groupName string, userID string){
	var GroupID string

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM groups WHERE title = ? ", groupName)
	err = row.Scan(&GroupID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			log.Println("Invalid id")
		}
		log.Println("Error retrieving id:", err)
	}

	_, err = db.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)",  &GroupID, userID, "accepted")
	if err != nil {
		log.Println("Error inserting group member into database:", err)
	}

	log.Println("User added to group")
}