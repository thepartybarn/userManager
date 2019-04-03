package main

import (
	"encoding/json"
	"fmt"
)

var (
	_users            []*UserInfo
	_groupPermissions map[GroupID][]string
	_tokens           map[Token]*UserInfo
)

type GroupID string
type Token string

type UserInfo struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	PasswordHash string    `json:"passwordHash"`
	CardID       string    `json:"cardID"`
	FriendOf     *UserInfo `json:"friendOf"`
	Permissions  []string  `json:"permissions"`
	Groups       []GroupID `json:"groups"`
}

func setupDataBase() {
	log.Trace("Setting Up Database")
	_groupPermissions = make(map[GroupID][]string)
	_tokens = make(map[Token]*UserInfo)
}

func seedDatabase() {
	user := createUser("Mike", "Schmidt", "1234", "1234", nil)
	_groupPermissions["ADMIN"] = []string{"CREATE_USER", "REMOVE_USER"}
	user.Groups = []GroupID{"ADMIN", "OWNER"}
	user.Permissions = []string{"CREATE_USER", "REMOVE_USER"}
	log.Tracef("Default User: %+v", user)
}
func tokenHasPermission(token Token, requiredPermission string) bool {
	HasPermission := false
	log.Tracef("Looking for UserId with Token %v", token)
	if _tokens[token] != nil {
		log.Tracef("Found: %+v for Token %v", _tokens[token], token)
		HasPermission = _tokens[token].HasPermission(requiredPermission)
	}
	return HasPermission
}

func (userInfo *UserInfo) GenerateToken() (token Token) {
	log.Tracef("User: %+v", userInfo)
	token = Token(fmt.Sprintf("%v.%v-Token", userInfo.Firstname, userInfo.Lastname))
	log.Tracef("Token Generated %v", token)
	_tokens[token] = userInfo
	return
}
func (userInfo *UserInfo) HasPermission(requiredPermission string) bool {
	log.Tracef("User: %+v", userInfo)

	for _, group := range userInfo.Groups {
		if group == "OWNER" {
			log.Trace("Is Owner")
			return true
		}
	}

	for _, permission := range userInfo.Permissions {
		if permission == requiredPermission {
			log.Trace("Has Permission")
			return true
		}
	}
	return false
}

func getUserInfoForToken(tokenString string) (outputData map[string]interface{}) {
	outputData = make(map[string]interface{})
	token := Token(tokenString)
	if _tokens[token] != nil {
		log.Tracef("User: %+v", _tokens[token])
		inBytes, _ := json.Marshal(_tokens[token])
		json.Unmarshal(inBytes, &outputData)
	}

	return
}
func getToken(username, password, cardID string) (token Token) {
	log.Tracef("getToken (Username: %s, Password: %s, CardID: %s)", username, password, cardID)

	if username != "" && password != "" {
		passwordHash := hashPassword(password)

		for _, userInfo := range _users {
			if fmt.Sprintf("%v.%v", userInfo.Firstname, userInfo.Lastname) == username && userInfo.PasswordHash == passwordHash {
				token = userInfo.GenerateToken()
			}
		}
	} else if cardID != "" {
		for _, userInfo := range _users {
			if cardID == userInfo.CardID {
				token = userInfo.GenerateToken()
			}
		}
	}
	return
}
func closeToken(token string) {
	delete(_tokens, Token(token))
}

func hashPassword(password string) (passwordHash string) {
	passwordHash = password
	return
}

func createUser(firstname, lastname, password, cardID string, friendOf *UserInfo) *UserInfo {
	userToAdd := UserInfo{firstname, lastname, hashPassword(password), cardID, friendOf, make([]string, 0), make([]GroupID, 0)}
	_users = append(_users, &userToAdd)
	return &userToAdd
}
