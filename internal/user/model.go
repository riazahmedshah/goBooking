package user

type User struct {
	ID        string  `json:"id" db:"id"`
	FirstName string  `json:"firstName" db:"first_name"`
	LastName  *string `json:"lastName" db:"last_name"`
	Email     string  `json:"email" db:"email"`
	Password  string  `json:"password" db:"password"`
	Role      string  `json:"role" db:"role"`
}

type Host struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"userId" db:"user_id"`
	StateName string `json:"stateName" db:"state_name"`
	City      string `json:"city" db:"city"`
	Area      string `json:"area" db:"area"`
}
