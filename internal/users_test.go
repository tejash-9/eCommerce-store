package internal

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test RegisterUser with a unique email
func TestRegisterUser_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Act
	user, err := shoppingApp.RegisterUser("hailey", "hailey@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "hailey", user.Name)
	assert.Equal(t, "hailey@example.com", user.Email)
	assert.Len(t, shoppingApp.Users, 1)
	assert.Equal(t, user.Id, shoppingApp.UserMap["hailey@example.com"])
}

// Test RegisterUser with an already registered email
func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user
	user, err := shoppingApp.RegisterUser("aaron", "aaron@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Act
	_, err = shoppingApp.RegisterUser("alex", "aaron@example.com")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "Email already exists", err.Error())
}

// Test GetUser with an existing user ID
func TestGetUser_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Shahrukh", "shah@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Act
	retrievedUser, err := shoppingApp.GetUser(user.Id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user, retrievedUser)
}

func TestRemoveUser_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user
	user, err := shoppingApp.RegisterUser("Prem", "prem@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	
	err = shoppingApp.RemoveUser(user.Id)

	// Assert
	assert.NoError(t, err)
}

func TestRemoveUser_UserNotFound(t *testing.T) {
	shoppingApp := createMockEngine()
	
	err := shoppingApp.RemoveUser("nonExistentUserId")

	// Assert
	assert.Error(t, err)
}

// Test GetUser with a non-existent user ID
func TestGetUser_UserNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Act
	retrievedUser, err := shoppingApp.GetUser("nonExistentUserId")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
	assert.Equal(t, "User not found!", err.Error())
}

// Test GetUser for deleted user
func TestGetUser_UserRemoved(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user
	user, err := shoppingApp.RegisterUser("Abdul", "abdul@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	
	err = shoppingApp.RemoveUser(user.Id)

	// Assert
	assert.NoError(t, err)

	// Act
	retrievedUser, err := shoppingApp.GetUser(user.Id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
	assert.Equal(t, "User not found!", err.Error())
}

// Test GetUserByUsername with an existing username
func TestGetUserByUsername_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user
	user, err := shoppingApp.RegisterUser("smith", "smith@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Act
	retrievedUser, err := shoppingApp.GetUserByUsername("smith@example.com")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user, retrievedUser)
}

// Test GetUserByUsername with a non-existent username
func TestGetUserByUsername_UserNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Act
	retrievedUser, err := shoppingApp.GetUserByUsername("nonexistent.email@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}

// Test GetUserByUsername with a valid username that leads to a failed GetUser
func TestGetUserByUsername_UserRemoved(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user
	user, err := shoppingApp.RegisterUser("John cena", "john@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	
	err = shoppingApp.RemoveUser(user.Id)

	// Assert
	assert.NoError(t, err)

	// Act
	retrievedUser, err := shoppingApp.GetUserByUsername("john@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}
