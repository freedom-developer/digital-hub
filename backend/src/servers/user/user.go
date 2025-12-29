package user

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(255);not null;unique" json:"username"`
	Email     string    `gorm:"type:varchar(255)" json:"email"`
	Password  string    `gorm:"type:varchar(255);not pull" json:"password"`
	Role      int32     `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     int32  `json:"role" binding:"required,oneof=1 2 3"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      int32  `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (us *UserService) getUserByID(id string) (*User, error) {
	// Implementation for retrieving user by ID
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) getUserByName(username string) (*User, error) {
	// Implementation for retrieving user by username
	var user User
	err := us.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) getUserByEmail(email string) (*User, error) {
	// Implementation for retrieving user by email
	var user User
	err := us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) createUser(user *User) error {
	// Implementation for creating a new user
	return us.db.Create(user).Error
}

func (us *UserService) deleteUserByID(id string) error {
	// Implementation for deleting a user
	return us.db.Where("id = ?", id).Delete(&User{}).Error
}

func (us *UserService) deleteUserByName(username string) error {
	// Implementation for deleting a user by username
	return us.db.Where("username = ?", username).Delete(&User{}).Error
}

func (us *UserService) updateUser(user *User) error {
	// Implementation for updating user information
	return us.db.Save(user).Error
}

func (us *UserService) registerUser(req *RegisterRequest) (*UserResponse, error) {
	// Implementation for user registration
	existUser, err := us.getUserByName(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}
	newUser := &User{
		ID:        generateID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashPassword(req.Password),
		Role:      req.Role,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = us.createUser(newUser)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Role:      newUser.Role,
		Active:    newUser.Active,
		CreatedAt: newUser.CreatedAt.Unix(),
		UpdatedAt: newUser.UpdatedAt.Unix(),
	}, nil
}

func hashPassword(s string) string {
	// Example implementation for hashing a password
	return fmt.Sprintf("hashed-%s", s)
}

func generateID() string {
	// Example implementation for generating a unique ID
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}
