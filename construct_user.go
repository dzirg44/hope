package main

import (
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
	//"io"
	"fmt"
)

type User struct {
	Id             string
	Email          string
	HashedPassword string
	Username       string
}
const (
	passwordLength = 8
	hashCost = 10
	userIDLenght
)

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.Id
}
func NewUser(username, email, password string) (User, []error) {
	user := User {
		Email: email,
		Username: username,
	}
	var validationError []error


	if username == "" {
        validationError = append(validationError, errNoUsername)
	}
	if email == "" {
		validationError = append(validationError, errNoEmail)
	}
	if password == "" {
		validationError = append(validationError, errNoPassword)
	}
	if len(password) < passwordLength && password != "" {
		validationError = append(validationError, errPasswordTooShort)
	}
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		validationError = append(validationError, err)
	}
	if existingUser != nil {
		validationError = append(validationError, errUsernameExists)
	}
	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		validationError = append(validationError, err)
	}
	if existingUser != nil {
		validationError = append(validationError, errEmailExists)
	}
	if validationError != nil {
		return user, validationError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	user.HashedPassword = string(hashedPassword)
	user.Id = GenerateID("usr", userIDLenght)
	return user,  validationError
}
func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error){
	out := *user
	out.Email = email
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.Id != user.Id {
		return out, errEmailExists
	}
	user.Email = email
	if currentPassword == "" {
		return out, nil
	}
	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
			[]byte(currentPassword),
	) != nil {
		return out, errPasswordIncorrect
	}
	if newPassword == "" {
		return out, errNoPassword
	}
	if len(newPassword) < passwordLength {
		return out, errPasswordTooShort
	}
	hashedPassword, err :=  bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return  out, errCredentialsIncorrect
	}
	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
			[]byte(password),
	) != nil {
		return out, errCredentialsIncorrect
	}
	return existingUser, nil
}