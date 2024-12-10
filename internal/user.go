package internal

import (
	"fmt"
)

// user represents a customer or user in the system
type user struct {
	Id 		  	string            	// Unique user ID
	Name   		string            	// Name of the user
	Email       string            	// Email address of the user
	Cart        map[string]int   	// Map of product IDs and quantities in the user's cart
}

// newUser creates and returns a new user instance
func newUser(id string, name string, email string) *user {
	return &user{
		Id:       id,        
		Name:     name,      
		Email:    email, 
		Cart:     make(map[string]int),
	}
}

// RegisterUser registers a new user in the system, using email as a unique identifier
func (s *shoppingEngine) RegisterUser(name string, email string) (*user, error) {
	// Check if the email is already registered
	if s.UserMap[email] != "" {
		return nil, fmt.Errorf("email already exists!") // Error if email is already registered
	}
	
	// Generate a unique ID and create a new user
	id := generateUUID()
	user := newUser(id, name, email)

	// Store the user in the system's user map and map email to user ID
	s.Users[id] = user
	s.UserMap[email] = id

	return user, nil
}

// GetUser retrieves a user by their unique user ID
func (s *shoppingEngine) GetUser(userId string) (*user, error) {
	// Check if the user exists in the system
	if s.Users[userId] == nil {
		return nil, fmt.Errorf("User not found") // Error if the user is not found
	}
	return s.Users[userId], nil
}

// GetUserByUsername retrieves a user by their email/username
func (s *shoppingEngine) GetUserByUsername(username string) (*user, error) {
	// Check if the username exists in the user map
	if s.UserMap[username] == "" {
		return nil, fmt.Errorf("Username %s doesn't exist", username) // Error if the username does not exist
	}
	
	// Retrieve the user using the user ID from the email
	user, err := s.GetUser(s.UserMap[username])
	if err != nil {
		// If error occurs, remove the invalid username from the map and return the error
		delete(s.UserMap, username)
		return nil, err
	}
	return user, nil
}