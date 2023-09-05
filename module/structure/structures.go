package structure


type RegistrationData struct {
	Nickname  string `json: "nickname"`
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
	Age       int 	 `json: "age"`
	Gender    string `json: "gender"`
	Email     string `json: "email"`
	Password  string `json: "password"`
	Avatar	  string `json:"avatar"`
	About     string `json:"about_me"`
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
	ID       int       
	UserID   string    
	Title    string    
	Content  string   
	Category string    
	Nickname string
	ImagePath string
	Comments []Comment
}

type Comment struct {
	ID      int    
	UserID  string 
	Content string 
	PostID  int 
	Nickname string
}

type UserData struct {
	UserID    string `json: "userid"`
	UUID      string `json:"user_uuid"`
	Nickname  string `json: "nickname"`
    FirstName string `json: "firstname"`
    LastName  string `json: "lastname"`
    Age       int 	 `json: "age"`
    Gender    string `json: "gender"`
    Email     string `json: "email"`
	Activity  string `json:"activity"`
	Avatar	  string `json:"avatar"`
	About     string `json:"about_me"`
	Privacy	  string `json:"privacy"`
}

type PrivateMessages struct {
	Content  string 
	Time     string 
	SenderId string 
	RecipientId string 
}

type PrivateMesssageSend struct {
	UserID   string 
	Nickname string 
	Content  string 
}

type UserMessageData struct {
	UserID   string
	Nickname string
	Offset   int
}

type Group struct{
	UUID		string
	Name		string
	Description	string
}
