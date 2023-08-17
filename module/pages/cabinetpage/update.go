package cabinetpage
import (
	"database/sql"
	"fmt"

	"social-network/module/structure"
)

// UpdateUserData handles the HTTP request to update user data
func UpdateUserData(userData structure.UserData) {

	fmt.Println(userData)
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE users SET nickname=?, first_name=?, last_name=?, age=?, gender=?, email=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	// Execute the SQL statement with the updated user data
	_, err = stmt.Exec(userData.Nickname, userData.FirstName, userData.LastName, userData.Age, userData.Gender, userData.Email, userData.UserID)
	if err != nil {
		fmt.Println(err)
	}
}