package structure


type RegistrationData struct {
	Nickname  string `json: "nickname"`
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
	Age       int 	 `json: "age"`
	Gender    string `json: "gender"`
	Email     string `json: "email"`
	Password  string `json: "password"`
}

type LoginData struct {
	Email    string	
	Password string
}

type Message struct {
	Type 	string `json:"type"`
	UUID 	string `json:"uuid"`
}

type Post struct {
	ID       int       `json:"id"`
	UserID   string    `json:"uuid"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Category string    `json:"category"`
	Nickname string		`json:"nickname"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID      int    `json:"id"`
	UserID  string `json:"uuid"`
	Content string `json:"content"`
	PostID  string `json:"post_id"`
	Nickname string `json:"nickname"`
}

type UserData struct {
	UserID    string `json: "userid"`
	Nickname  string `json: "nickname"`
    FirstName string `json: "firstname"`
    LastName  string `json: "lastname"`
    Age       int 	 `json: "age"`
    Gender    string `json: "gender"`
    Email     string `json: "email"`
}
