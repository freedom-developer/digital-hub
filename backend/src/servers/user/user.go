package user

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      int32  `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
