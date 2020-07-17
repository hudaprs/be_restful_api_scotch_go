package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User Model
type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);unique_index" json:"email"`
	FirstName    string `gorm:"size:100;not null" json:"firstName"`
	LastName     string `gorm:"size:100;not null" json:"lastName"`
	Password     string `gorm:"size:100;not null" json:"password"`
	ProfileImage string `gorm:"size:255" json:"profileImage"`
}

// HashPassword hashes password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("Password Incorrect")
	}
	return nil
}

// BeforeSave hashes user password
func (user *User) BeforeSave() error {
	password := strings.TrimSpace(user.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Prepare strips user input of any white spaces
func (user *User) Prepare() {
	user.Email = strings.TrimSpace(user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.ProfileImage = strings.TrimSpace(user.ProfileImage)
}

// Validate user input
func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Email == "" {
			return errors.New("Email is required")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	default:
		if user.FirstName == "" {
			return errors.New("FirstName is required")
		}
		if user.LastName == "" {
			return errors.New("LastName is required")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	}
}

// SaveUser add a user to the database
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	// Debug a single operation, show detailed log for this operation
	err = db.Debug().Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GetUser returns a user based on Email
func (user *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("email = ?", user.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// GetAllUsers returns a list of all the user
func GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}
