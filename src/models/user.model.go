package models

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password,omitempty" json:"password"`
}

// Validate validates the user input for registration.
func (u *User) Validate() error {
	// Check name
	if u.Name == "" {
		return errors.New("name is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Check password
	if utf8.RuneCountInString(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (u *User) GenerateAccessToken() (string,error){
	expiration,err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))

	if err!=nil{
		return "",fmt.Errorf("error while parsing expiration time : %v",err)
		
	}

	claims := jwt.MapClaims{
		"id":u.ID,
		"exp":time.Now().Add(expiration).Unix(),//expiration time
		"iat":time.Now().Unix(), //Issued at time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	signedToken,err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET"))) 

	if err!=nil{
		return "",fmt.Errorf("error while signing token : %v",err)
	}

	return signedToken,nil
}

func (u *User) HashPassword() error{
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(input *LoginInfo) (bool,error){
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return false, errors.New("invalid credentials")
	}
	return true,nil
}



// Login Structure
type LoginInfo struct{
	Email    string 
	Password string 
}

func (u *LoginInfo) Validate() error {
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Check password
	if utf8.RuneCountInString(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long ")
	}

	return nil
}


// A new struct for the response without the password field
type UserResponse struct {
    ID   primitive.ObjectID `json:"id,omitempty"`
    Name string             `json:"name"`
	Email string            `json:"email"`
}

// Function to convert from User to UserResponse
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:   u.ID,
        Name: u.Name,
		Email: u.Email,
    }
}