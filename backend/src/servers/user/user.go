package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:varchar(255); primaryKey" json:"id"`
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
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=3"`
	Role     int32  `json:"role"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	newUser := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
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

func (us *UserService) login(req *LoginRequest) (*UserResponse, error) {
	user, err := us.getUserByName(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	if !user.Active {
		return nil, errors.New("用户未激活")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}, nil
}
