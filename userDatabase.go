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
	_secretKeys       map[SecretKey]UserID
)

type UserID string
type GroupID string
type SecretKey string

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
	_secretKeys = make(map[SecretKey]UserID)
}

func seedDatabase() {
	_userInfo["MIKEUID"] = UserInfo{"Mike", "Schmidt"}
	_userLogin["MIKEUID"] = UserLogin{"mschmidt", "1234"}
	_userCardID["MIKEUID"] = "1234"
	_groupPermissions["ADMIN"] = []string{"CREATE_USER", "REMOVE_USER"}
	_userGroups["MIKEUID"] = []GroupID{"ADMIN"}
	//_userPermissions["MIKEUID"] = []Permission{CREATE_USER, REMOVE_USER}
}
func secretKeyHasPermission(secretKey SecretKey, requiredPermission string) bool {
	var Permissions []string
	log.Tracef("Looking for UserId with SecretKey %v", secretKey)
	userID := _secretKeys[secretKey]
	if userID != "" {
		log.Tracef("Found UserID: %v for SecretKey %v", userID, secretKey)
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

func getUserInfoForSecretKey(secretKeyString string) (outputData map[string]interface{}) {
	outputData = make(map[string]interface{})
	secretKey := SecretKey(secretKeyString)
	guid, ok := _secretKeys[secretKey]
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
func getSecretKey(username, password, cardID string) (secretKey string) {
	log.Tracef("getSecretKey (Username: %s, Password: %s, CardID: %s)", username, password, cardID)

	if username != "" && password != "" {
		passwordHash := hashPassword(password)

		for userID, userLogin := range _userLogin {
			if userLogin.Username == username && userLogin.PasswordHash == passwordHash {
				secretKey = string(generateSecretKey(userID))
			}
		}
	} else if cardID != "" {
		for userID, CardID := range _userCardID {
			if cardID == CardID {
				secretKey = string(generateSecretKey(userID))
			}
		}
	}
	return
}
func closeSecretKey(secretKey string) {
	delete(_secretKeys, SecretKey(secretKey))
}

func hashPassword(password string) (passwordHash string) {
	passwordHash = password
	return
}
func generateSecretKey(userID UserID) (secretKey SecretKey) {
	secretKey = SecretKey(fmt.Sprintf("%vSecretKey", userID))
	_secretKeys[secretKey] = userID
	return
}
