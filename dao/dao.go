package dao

import (
	"encoding/json"
	"os"
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	users     = make(map[string]User)
	usersFile = "users.json"
	mu        sync.RWMutex
)

func init() {
	loadUsers()
}
func loadUsers() {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.Open(usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			saveUsers()
			return
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		users = make(map[string]User)
	}
}
func saveUsers() {
	file, err := os.Create(usersFile)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(users)
}
func AddUser(username, password string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[username]; exists {
		return false
	}
	users[username] = User{
		Username: username,
		Password: password,
	}
	saveUsers()
	return true
}
func SelectUser(username string) bool {
	mu.RLock()
	defer mu.RUnlock()

	_, exists := users[username]
	return exists
}
func SelectPasswordFromUsername(username string) string {
	mu.RLock()
	defer mu.RUnlock()

	if user, exists := users[username]; exists {
		return user.Password
	}
	return ""
}
func UpdatePassword(username, newPassword string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[username]; !exists {
		return false
	}

	users[username] = User{
		Username: username,
		Password: newPassword,
	}

	saveUsers()
	return true
}
func GetAllUsers() map[string]User {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]User)
	for k, v := range users {
		result[k] = v
	}
	return result
}
