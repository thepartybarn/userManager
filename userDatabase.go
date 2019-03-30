package main

import (
	"fmt"
)

var (
	_userInfo         map[UserID]UserInfo
	_userLogin        map[UserID]UserLogin
	_userCardID       map[UserID]string
	_userFriends      map[UserID][]UserID
	_userPermissions  map[UserID][]string
	_userGroups       map[UserID][]GroupID
	_groupPermissions map[GroupID][]string
	_tokens           map[Token]UserID
)

type UserID string
type GroupID string
type Token string

type UserInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
type UserLogin struct {
	Username     string `json:"username"`
	PasswordHash string `json:"PasswordHash"`
}

func setupDataBase() {
	log.Trace("Setting Up Database")
	_userInfo = make(map[UserID]UserInfo)
	_userLogin = make(map[UserID]UserLogin)
	_userCardID = make(map[UserID]string)
	_userFriends = make(map[UserID][]UserID)
	_userPermissions = make(map[UserID][]string)
	_userGroups = make(map[UserID][]GroupID)
	_groupPermissions = make(map[GroupID][]string)
	_tokens = make(map[Token]UserID)
}

func seedDatabase() {
	_userInfo["MIKEUID"] = UserInfo{"Mike", "Schmidt"}
	_userLogin["MIKEUID"] = UserLogin{"mschmidt", "1234"}
	_userCardID["MIKEUID"] = "1234"
	_groupPermissions["ADMIN"] = []string{"CREATE_USER", "REMOVE_USER"}
	_userGroups["MIKEUID"] = []GroupID{"ADMIN"}
	//_userPermissions["MIKEUID"] = []Permission{CREATE_USER, REMOVE_USER}
}
func tokenHasPermission(token Token, requiredPermission string) bool {
	var Permissions []string
	log.Tracef("Looking for UserId with Token %v", token)
	userID := _tokens[token]
	if userID != "" {
		log.Tracef("Found UserID: %v for Token %v", userID, token)
		Permissions = getAllPermissionsForUserID(userID)
	}
	return PermissionInPermissions(requiredPermission, Permissions)
}
func getAllPermissionsForUserID(userID UserID) []string {
	permissions := _userPermissions[userID]
	for _, groupID := range _userGroups[userID] {
		permissions = append(permissions, _groupPermissions[groupID]...)
	}
	return permissions
}
func PermissionInPermissions(requiredPermission string, Permissions []string) bool {
	for _, permission := range Permissions {
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
	guid, ok := _tokens[token]
	if ok {
		userInfo, ok := _userInfo[guid]
		if ok {
			outputData["Firstname"] = userInfo.Firstname
			outputData["Lastname"] = userInfo.Lastname
		}
		userLogin, ok := _userLogin[guid]
		if ok {
			outputData["Username"] = userLogin.Username
		}
	}

	return
}
func getToken(username, password, cardID string) (token string) {
	log.Tracef("getToken (Username: %s, Password: %s, CardID: %s)", username, password, cardID)

	if username != "" && password != "" {
		passwordHash := hashPassword(password)

		for userID, userLogin := range _userLogin {
			if userLogin.Username == username && userLogin.PasswordHash == passwordHash {
				token = string(generateToken(userID))
			}
		}
	} else if cardID != "" {
		for userID, CardID := range _userCardID {
			if cardID == CardID {
				token = string(generateToken(userID))
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
func generateToken(userID UserID) (token Token) {
	token = Token(fmt.Sprintf("%vToken", userID))
	_tokens[token] = userID
	return
}
