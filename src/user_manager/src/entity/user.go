package entity

// id, email, nickname are unique.
// Storage system implementation should throw error if duplicates are inserted
type User struct {
	ID         uint            `json:"id"`          // primary key
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	Nickname   string          `json:"nickname"`
	Password   string          `json:"password"`
	Email      string          `json:"email"`
	Country    string          `json:"country"`
}
