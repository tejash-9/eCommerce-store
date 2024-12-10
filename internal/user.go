package internal

import (
	"fmt"
)

type user struct {
	Id 		  	string
	Name   		string
	Email       string
	Cart        map[string]int
}

func newUser(id string, name string, email string) *user {
	return &user{
		Id:       id,
		Name:     name,
		Email:    email,
		Cart:     make(map[string]int),
	}
}

func (s *shoppingEngine) RegisterUser(name string, email string) (*user, error) {
	if s.UserMap[email] != "" {
		return nil, fmt.Errorf("email already exists!")
	}
	id := generateUUID()
	user := newUser(id, name, email)

	s.Users[id] = user
	s.UserMap[email] = id

	return user, nil
}

func (s *shoppingEngine) GetUser(userId string) (*user, error) {
	if s.Users[userId] == nil {
		return nil, fmt.Errorf("User not found")
	}
	return s.Users[userId], nil
}

func (s *shoppingEngine) GetUserByUsername(username string) (*user, error) {
	if s.UserMap[username] == "" {
		return nil, fmt.Errorf("Username %s doesn't exist", username)
	}
	user, err := s.GetUser(s.UserMap[username])
	if err != nil {
		delete(s.UserMap, username)
		return nil, err
	}
	return user, nil
}